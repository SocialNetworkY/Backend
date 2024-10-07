package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"net/http"
	"strconv"
	"time"
)

type user struct {
	ID        uint
	Role      uint
	Banned    bool
	Reason    string
	ExpiredAt time.Time
}

const (
	requesterLocals = "requester"
	postIDParam     = "post_id"
	postLocals      = "post"
)

func (h *Handler) authenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(constant.HTTPAuthorizationHeader)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		var err error
		requester := &user{}

		requester.ID, err = h.ag.Authenticate(c.Request().Context(), authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		requester.Role, err = h.ug.GetUserRole(c.Request().Context(), authHeader, requester.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		requester.Banned, requester.Reason, requester.ExpiredAt, err = h.ug.IsUserBan(c.Request().Context(), authHeader, requester.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(requesterLocals, requester)

		return next(c)
	}
}

func (h *Handler) banMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requester, ok := c.Get(requesterLocals).(*user)
		if !ok {
			return echo.ErrUnauthorized
		}

		if requester.Banned {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("You are banned: %s.\nUntil: %s", requester.Reason, requester.ExpiredAt))
		}

		return next(c)
	}
}

func (h *Handler) setPostByIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := getUintParam(c, postIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		post, err := h.ps.Find(postID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		c.Set(postLocals, post)

		return next(c)
	}
}

func getUintParam(c echo.Context, key string) (uint, error) {
	param := c.Param(key)
	if param == "" {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "missing parameter "+key)
	}

	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid parameter "+key)
	}

	return uint(id), nil
}
