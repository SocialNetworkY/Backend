package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"net/http"
	"time"
)

func (h *Handler) initAdminApi(group *echo.Group) {
	initUserEndpoints := func(user *echo.Group) {
		user.POST("/ban", h.banUser)
		user.GET("/bans", h.getUserBans)
	}

	admin := group.Group("/admin", h.authenticationMiddleware, h.adminMiddleware)
	{
		users := admin.Group("/users")
		{
			initUserEndpoints(users.Group(fmt.Sprintf("/:%s", paramUserID), h.setUserByIDFromParam))
			initUserEndpoints(users.Group(fmt.Sprintf("/@:%s", paramUsername), h.setUserByUsernameFromParam))
		}

		admin.POST("/unban", h.unbanByBanID)
	}
}

func (h *Handler) banUser(c echo.Context) error {
	var req struct {
		Reason   string `json:"reason"`
		Duration string `json:"duration"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if req.Reason == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid reason")
	}

	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid duration")
	}

	if duration == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid duration")
	}

	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if user.Role > requester.Role {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	if user.Banned {
		return echo.NewHTTPError(http.StatusBadRequest, "user already banned")
	}

	if err := h.bs.BanUser(user.ID, requester.ID, req.Reason, duration); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) unbanByBanID(c echo.Context) error {
	var req struct {
		BanID  uint   `json:"banID"`
		Reason string `json:"reason"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	ban, err := h.bs.FindBan(req.BanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "ban not found")
	}

	if ban.UserID == requester.ID {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	if !ban.Active {
		return echo.NewHTTPError(http.StatusBadRequest, "ban already inactive")
	}

	bannedBy, err := h.us.Find(ban.BannedBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "admin who banned not found")
	}

	if bannedBy.Role > requester.Role {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	if err := h.bs.UnbanByBanID(ban.ID, requester.ID, req.Reason); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) getUserBans(c echo.Context) error {
	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, struct {
		Bans []*model.Ban `json:"bans"`
	}{
		Bans: user.Bans,
	})
}
