package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handler) initTokenApi(api *echo.Group) {
	api.POST("/refresh", h.refreshToken)
}

func (h *Handler) refreshToken(c echo.Context) error {
	refreshTokenCookie, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, err := h.tokenService.VerifyRefreshToken(refreshTokenCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	accessToken, refreshToken, err := h.tokenService.Generate(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return h.setAndReturnTokens(c, accessToken, refreshToken)
}

type (
	accessTokenResp struct {
		Token string `json:"token"`
	}
)

// func set refresh token cookie and return access token
func (h *Handler) setAndReturnTokens(c echo.Context, accessToken, refreshToken string) error {
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(h.RefreshTokenDuration),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, accessTokenResp{
		Token: accessToken,
	})
}
