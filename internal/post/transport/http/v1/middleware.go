package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/labstack/echo/v4"
)

const (
	postIDParam    = "post_id"
	userIDParam    = "user_id"
	commentIDParam = "comment_id"

	postLocals      = "post"
	userLocals      = "user"
	requesterLocals = "requester"
	commentLocals   = "comment"

	skipQuery  = "skip"
	limitQuery = "limit"
	queryQuery = "query"

	defaultSkip  = 0
	defaultLimit = 10
)

func (h *Handler) authenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(constant.HTTPAuthorizationHeader)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		requesterID, err := h.ag.Authenticate(c.Request().Context(), authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		requester, err := h.ug.UserInfo(c.Request().Context(), requesterID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(requesterLocals, requester)

		return next(c)
	}
}

func (h *Handler) setUserByIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUintParam(c, userIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		user, err := h.ug.UserInfo(c.Request().Context(), userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(userLocals, user)

		return next(c)
	}
}

func (h *Handler) banMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requester, ok := c.Get(requesterLocals).(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
		}

		if requester.Banned {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("You are banned: %s.\nUntil: %s", requester.BanReason, requester.BanExpiredAt))
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

func (h *Handler) setCommentByIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		commentID, err := getUintParam(c, commentIDParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		comment, err := h.cs.Find(commentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		c.Set(commentLocals, comment)

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

func skipLimitQuery(c echo.Context) (int, int) {
	skip := defaultSkip
	if s, err := strconv.Atoi(c.QueryParam(skipQuery)); err == nil {
		skip = s
	}
	limit := defaultLimit
	if l, err := strconv.Atoi(c.QueryParam(limitQuery)); err == nil {
		limit = l
	}

	return skip, limit
}
