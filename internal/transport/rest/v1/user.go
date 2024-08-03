package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
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
	var body userLoginReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch {
	case body.Login == "":
		return echo.NewHTTPError(http.StatusBadRequest, "login is empty")
	case body.Password == "":
		return echo.NewHTTPError(http.StatusBadRequest, "password is empty")
	}

	accessToken, refreshToken, err := h.userService.Login(body.Login, body.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokensResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

type (
	userLoginReq struct {
		Login    string `json:"login" binding:"required"`
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
	var body userRegisterReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch {
	case body.Email == "":
		return echo.NewHTTPError(http.StatusBadRequest, "email is empty")
	case body.Username == "":
		return echo.NewHTTPError(http.StatusBadRequest, "username is empty")
	case body.Password == "":
		return echo.NewHTTPError(http.StatusBadRequest, "password is empty")
	case !isEmailValid(body.Email):
		return echo.NewHTTPError(http.StatusBadRequest, "email is invalid")
	}

	accessToken, refreshToken, err := h.userService.Register(body.Username, body.Email, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokensResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type (
	userRegisterReq struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)
