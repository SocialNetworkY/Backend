package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// TODO: Implement functional in refreshToken
func (h *Handler) initTokenApi(api *echo.Group) {
	api.POST("/refresh", h.refreshToken)
}

// @Summary      Refresh token
// @Description  Refresh jwt token
// @Tags         Token
// @Accept       json
// @Produce      json
// @Param input body refreshTokenReq true "refresh token"
// @Success      200  {object}  tokensResp
// @Failure      default  {object}  echo.HTTPError
// @Router       /refresh [post]
func (h *Handler) refreshToken(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "not implemented")
}

type (
	refreshTokenReq struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	tokensResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
