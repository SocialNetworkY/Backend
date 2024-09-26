package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"net/http"
	"strconv"
)

const (
	paramUserID   = "userID"
	paramUsername = "username"
)

func (h *Handler) initUserApi(api *echo.Group) {
	api.GET(fmt.Sprintf("/@:%s", paramUsername), h.getUser)
	api.GET(fmt.Sprintf("/:%s", paramUserID), h.putUser, h.authenticationMiddleware)
}

// @Summary Get user by username
// @Description Get user by username
// @Tags user
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} getUserResponse
// @Router /{username} [get]
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

// @Summary Update user
// @Description Update user
// @Tags user
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param Authorization header string true "Authorization"
// @Param request body putUserRequest true "Request body"
// @Success 200
// @Router /{userID} [put]
func (h *Handler) putUser(c echo.Context) error {
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

	user, err := h.us.Find(uint(userID))
	if err != nil {
		return err
	}

	if !(requester.ID == user.ID || requester.Role >= model.RoleAdminLvl1) {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to update this user")
	}

	switch {
	case requestBody.Nickname != "":
		if err := h.us.ChangeNickname(user.ID, requestBody.Nickname); err != nil {
			return err
		}
	case requestBody.Username != "":
		if err := h.us.ChangeUsername(user.ID, authHeader, requestBody.Username); err != nil {
			return err
		}
	case requestBody.Email != "":
		if err := h.us.ChangeEmail(user.ID, authHeader, requestBody.Email); err != nil {
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

// @Summary Delete user
// @Description Delete user
// @Tags user
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param Authorization header string true "Authorization"
// @Success 200
// @Router /{userID} [delete]
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

	user, err := h.us.Find(uint(userID))
	if err != nil {
		return err
	}

	if !(requester.ID == user.ID || requester.Role >= model.RoleAdminLvl3) {
		return echo.NewHTTPError(http.StatusForbidden, "you don't have permission to delete this user")
	}

	if err := h.us.Delete(user.ID, authHeader); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
