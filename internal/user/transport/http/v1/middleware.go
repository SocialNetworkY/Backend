package v1

import (
	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	userLocals      = "user"
	requesterLocals = "requester"

	paramUserID   = "userID"
	paramUsername = "username"

	queryQuery = "query"
	skipQuery  = "skip"
	limitQuery = "limit"

	defaultSkip  = 0
	defaultLimit = 10
)

func (h *Handler) authenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(constant.HTTPAuthorizationHeader)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		userID, err := h.ag.Authenticate(c.Request().Context(), authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		requester, err := h.us.Find(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(requesterLocals, requester)

		return next(c)
	}
}

func (h *Handler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requester, ok := c.Get(requesterLocals).(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
		}

		if !requester.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "User is not admin")
		}

		return next(c)
	}
}

func (h *Handler) setUserByIDFromParam(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUintParam(c, paramUserID)
		if err != nil {
			return err
		}

		user, err := h.us.Find(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		c.Set(userLocals, user)

		return next(c)
	}
}

func (h *Handler) setUserByUsernameFromParam(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param(paramUsername)
		user, err := h.us.FindByUsername(username)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		c.Set(userLocals, user)

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
