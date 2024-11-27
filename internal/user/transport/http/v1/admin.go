package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/labstack/echo/v4"
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
			users.GET("", h.getUsers)
			users.GET("/stats", h.getUsersStats)
			initUserEndpoints(users.Group(fmt.Sprintf("/:%s", paramUserID), h.setUserByIDFromParam))
			initUserEndpoints(users.Group(fmt.Sprintf("/@:%s", paramUsername), h.setUserByUsernameFromParam))
		}

		bans := admin.Group("/bans")
		{
			bans.GET("", h.getBans)
			bans.GET("/stats", h.getBansStats)
			bans.GET("/search", h.searchBans)
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

	if err := h.bs.Ban(user.ID, requester.ID, req.Reason, duration); err != nil {
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

	ban, err := h.bs.Find(req.BanID)
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

	if err := h.bs.Unban(ban.ID, requester.ID, req.Reason); err != nil {
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

func (h *Handler) getBans(c echo.Context) error {
	skip, limit := skipLimitQuery(c)
	bans, err := h.bs.FindSome(skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bans)
}

func (h *Handler) searchBans(c echo.Context) error {
	query := c.QueryParam(queryQuery)
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query parameter")
	}

	skip, limit := skipLimitQuery(c)
	bans, err := h.bs.Search(query, skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bans)
}

func (h *Handler) getBansStats(c echo.Context) error {
	stats, err := h.bs.Statistic()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) getUsersStats(c echo.Context) error {
	stats, err := h.us.Statistic()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) getUsers(c echo.Context) error {
	skip, limit := skipLimitQuery(c)
	users, err := h.us.FindSome(skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
