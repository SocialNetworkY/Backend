package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"net/http"
)

func (h *Handler) authenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(constant.HTTPAuthorizationHeader)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		user, err := h.authenticationService.Auth(authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(userLocals, user)

		return next(c)
	}
}
