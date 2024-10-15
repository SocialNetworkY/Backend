package v1

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/report/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	postIDParam = "post_id"
	userIDParam = "user_id"

	postLocals      = "post"
	userLocals      = "user"
	requesterLocals = "requester"

	skipQuery  = "skip"
	limitQuery = "limit"

	defaultSkip  = 0
	defaultLimit = 10
)

func (h *Handler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
