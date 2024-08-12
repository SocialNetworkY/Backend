package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterAuthService/internal/core"
	"net/http"
)

func (h *Handler) initUserChangeApi(api *echo.Group) {
	change := api.Group("/change")
	{
		change.PATCH("/email", h.userChangeEmail, h.authenticationMiddleware)
		change.PATCH("/username", h.userChangeUsername, h.authenticationMiddleware)
		change.PATCH("/password", h.userChangePassword, h.authenticationMiddleware)
	}
}

// @Summary      Change Email
// @Description  Change user email
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param input body userChangeEmailReq true "new email"
// @Success      204
// @Failure      default  {object}  echo.HTTPError
// @Router       /change/email [patch]
func (h *Handler) userChangeEmail(c echo.Context) error {
	user, ok := c.Get(userLocals).(*core.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	var body userChangeEmailReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Email(body.Email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.userService.ChangeEmail(user.ID, body.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

type (
	userChangeEmailReq struct {
		Email string `json:"email" binding:"required"`
	}
)

// @Summary      Change Username
// @Description  Change user username
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param input body userChangeUsernameReq true "new username"
// @Success      204
// @Failure      default  {object}  echo.HTTPError
// @Router       /change/username [patch]
func (h *Handler) userChangeUsername(c echo.Context) error {
	user, ok := c.Get(userLocals).(*core.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	var body userChangeUsernameReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Username(body.Username); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.userService.ChangeUsername(user.ID, body.Username); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

type (
	userChangeUsernameReq struct {
		Username string `json:"username" binding:"required"`
	}
)

// @Summary      Change Password
// @Description  Change user password
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param input body userChangePasswordReq true "new password"
// @Success      204
// @Failure      default  {object}  echo.HTTPError
// @Router       /change/password [patch]
func (h *Handler) userChangePassword(c echo.Context) error {
	user, ok := c.Get(userLocals).(*core.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	var body userChangePasswordReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Password(body.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.userService.ChangePassword(user.ID, body.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

type (
	userChangePasswordReq struct {
		Password string `json:"password" binding:"required"`
	}
)
