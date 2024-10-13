package v1

import (
	"github.com/SocialNetworkY/Backend/internal/auth/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handler) initUserApi(api *echo.Group) {
	api.POST("/login", h.userLogin)
	api.POST("/register", h.userRegister)
	api.GET("/activate/:token", h.userActivate)
	api.GET("/authenticate", h.userAuthenticate, h.authenticationMiddleware)
	api.GET("/info", h.userInfo, h.authenticationMiddleware)
	api.PATCH("/change-password", h.userChangePassword, h.authenticationMiddleware)
}

func (h *Handler) userLogin(c echo.Context) error {
	var body struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required,password"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, refreshToken, err := h.userService.Login(body.Login, body.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return h.setAndReturnTokens(c, accessToken, refreshToken)
}

func (h *Handler) userRegister(c echo.Context) error {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,username"`
		Password string `json:"password" validate:"required,password"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	activationToken, err := h.userService.Register(body.Username, body.Email, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ActivationToken string `json:"activation_token"`
	}{
		ActivationToken: activationToken,
	})
}

func (h *Handler) userActivate(c echo.Context) error {
	activationToken := c.Param("token")

	accessToken, refreshToken, err := h.userService.Activate(activationToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return h.setAndReturnTokens(c, accessToken, refreshToken)
}

func (h *Handler) userAuthenticate(c echo.Context) error {
	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	return c.JSON(http.StatusOK, struct {
		UserID uint `json:"user_id"`
	}{
		UserID: user.ID,
	})
}

func (h *Handler) userInfo(c echo.Context) error {
	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	return c.JSON(http.StatusOK, struct {
		ID        uint      `json:"id"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *Handler) userChangePassword(c echo.Context) error {
	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	var body struct {
		Password string `json:"password" validate:"required,password"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.userService.ChangePassword(user.ID, body.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
