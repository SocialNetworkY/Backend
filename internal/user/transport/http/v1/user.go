package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"net/http"
)

func (h *Handler) initUserApi(group *echo.Group) {
	initUserEndpoints := func(user *echo.Group) {
		user.GET("", h.getUser)
		user.GET("/full", h.getFullUser, h.authenticationMiddleware)
		user.PATCH("", h.patchUser, h.authenticationMiddleware)
		user.DELETE("", h.deleteUser, h.authenticationMiddleware)
		user.POST("/avatar", h.postAvatar, h.authenticationMiddleware)
	}

	users := group.Group("/users")
	{
		initUserEndpoints(users.Group(fmt.Sprintf("/@:%s", paramUsername), h.setUserByUsernameFromParam))
		initUserEndpoints(users.Group(fmt.Sprintf("/:%s", paramUserID), h.setUserByIDFromParam))
	}
}

func (h *Handler) getUser(c echo.Context) error {
	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	})
}

func (h *Handler) getFullUser(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if requester.ID != user.ID && !requester.Admin {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to view this user")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) patchUser(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	authHeader := c.Request().Header.Get(authorizationHeader)
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
	}

	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	var requestBody struct {
		Nickname string `json:"nickname"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if requester.ID != user.ID && requester.Role < constant.RoleAdminLvl1 {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to update this user")
	}

	if requestBody.Nickname != "" {
		if err := h.us.ChangeNickname(user.ID, requestBody.Nickname); err != nil {
			return err
		}
	}
	if requestBody.Username != "" {
		if err := h.us.ChangeUsername(user.ID, authHeader, requestBody.Username); err != nil {
			return err
		}
	}
	if requestBody.Email != "" {
		if err := h.us.ChangeEmail(user.ID, authHeader, requestBody.Email); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) deleteUser(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	authHeader := c.Request().Header.Get(authorizationHeader)
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
	}

	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if requester.ID != user.ID && requester.Role < constant.RoleAdminLvl3 {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to delete this user")
	}

	if err := h.us.Delete(user.ID, authHeader); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) postAvatar(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if requester.ID != user.ID && requester.Role <= user.Role {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to update this user")
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing avatar file")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	if err := h.us.ChangeAvatar(user.ID, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
