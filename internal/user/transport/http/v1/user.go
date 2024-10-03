package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"net/http"
	"strconv"
)

const (
	paramUserID   = "userID"
	paramUsername = "username"
)

func (h *Handler) initUserApi(api *echo.Group) {
	api.GET(fmt.Sprintf("/@:%s", paramUsername), h.getUser)
	api.PATCH(fmt.Sprintf("/:%s", paramUserID), h.patchUser, h.authenticationMiddleware)
	api.DELETE(fmt.Sprintf("/:%s", paramUserID), h.deleteUser, h.authenticationMiddleware)
}

func (h *Handler) getUser(c echo.Context) error {
	username := c.Param(paramUsername)
	user, err := h.us.FindByUsername(username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &getUserResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     user.Role,
	})
}

type (
	getUserResponse struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Role     uint   `json:"role"`
	}
)

func (h *Handler) patchUser(c echo.Context) error {
	requester, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	authHeader := c.Request().Header.Get(authorizationHeader)
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
	}

	userID, err := strconv.Atoi(c.Param(paramUserID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	var requestBody putUserRequest
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if requester.ID != uint(userID) && requester.Role < constant.RoleAdminLvl1 {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to update this user")
	}

	if requestBody.Nickname != "" {
		if err := h.us.ChangeNickname(uint(userID), requestBody.Nickname); err != nil {
			return err
		}
	}
	if requestBody.Username != "" {
		if err := h.us.ChangeUsername(uint(userID), authHeader, requestBody.Username); err != nil {
			return err
		}
	}
	if requestBody.Email != "" {
		if err := h.us.ChangeEmail(uint(userID), authHeader, requestBody.Email); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

type (
	putUserRequest struct {
		Nickname string `json:"nickname"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)

func (h *Handler) deleteUser(c echo.Context) error {
	requester, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	authHeader := c.Request().Header.Get(authorizationHeader)
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
	}

	userID, err := strconv.Atoi(c.Param(paramUserID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	if requester.ID != uint(userID) && requester.Role < constant.RoleAdminLvl3 {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to delete this user")
	}

	if err := h.us.Delete(uint(userID), authHeader); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
