package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initUserApi(api *echo.Group) {
	api.POST("/login", h.userLogin)
	api.POST("/register", h.userRegister)
}

// @Summary      User login
// @Description  login bruh
// @Tags         User
// @Accept       json
// @Produce      json
// @Param input body userLoginReq true "user credentials"
// @Success      200  {object}  tokensResp
// @Failure      default  {object}  echo.HTTPError
// @Router       /login [post]
func (h *Handler) userLogin(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Not Implemented")

}

type (
	userLoginReq struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}
)

// @Summary      User register
// @Description  reg bruh
// @Tags         User
// @Accept       json
// @Produce      json
// @Param input body userRegisterReq true "user credentials"
// @Success      200  {object}  tokensResp
// @Failure      default  {object}  echo.HTTPError
// @Router       /register [post]
func (h *Handler) userRegister(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Not Implemented")
}

type (
	userRegisterReq struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)
