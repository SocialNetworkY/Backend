package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

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
	var body refreshTokenReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if body.RefreshToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "token is empty")
	}

	userID, err := h.tokenService.VerifyRefreshToken(body.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	accessToken, refreshToken, err := h.tokenService.Generate(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokensResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
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
