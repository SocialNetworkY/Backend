package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) authenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authorizationHeader)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		userID, err := h.ag.CheckAuth(c.Request().Context(), authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		user, err := h.us.Find(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(userLocals, user)

		return next(c)
	}
}
