package v1

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/report/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	reportIDParam = "report_id"

	reportLocals    = "report"
	requesterLocals = "requester"

	skipQuery  = "skip"
	limitQuery = "limit"

	defaultSkip  = 0
	defaultLimit = 10
)

func (h *Handler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(constant.HTTPAuthorizationHeader)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		requesterID, err := h.ag.Authenticate(c.Request().Context(), authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		requester, err := h.ug.UserInfo(c.Request().Context(), requesterID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(requesterLocals, requester)

		return next(c)
	}
}

func (h *Handler) banMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requester, ok := c.Get(requesterLocals).(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
		}

		if requester.Banned {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("You are banned: %s.\nUntil: %s", requester.BanReason, requester.BanExpiredAt))
		}

		return next(c)
	}
}

func (h *Handler) setReportByIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reportID, err := getUintParam(c, reportIDParam)
		if err != nil {
			return err
		}

		report, err := h.rs.Get(reportID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.Set(reportLocals, report)

		return next(c)
	}
}

func (h *Handler) checkAccessMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		report, ok := c.Get(reportLocals).(*model.Report)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to get report")
		}

		requester, ok := c.Get(requesterLocals).(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
		}

		if report.UserID != requester.ID && requester.Role < constant.RoleAdminLvl1 {
			return echo.NewHTTPError(http.StatusForbidden, "You don't have access to this report")

		}

		return next(c)
	}
}

func (h *Handler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requester, ok := c.Get(requesterLocals).(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
		}

		if requester.Role < constant.RoleAdminLvl1 {
			return echo.NewHTTPError(http.StatusForbidden, "You don't have access to this report")
		}

		return next(c)
	}
}

func skipLimitQuery(c echo.Context) (int, int) {
	skip := defaultSkip
	if s, err := strconv.Atoi(c.QueryParam(skipQuery)); err == nil {
		skip = s
	}
	limit := defaultLimit
	if l, err := strconv.Atoi(c.QueryParam(limitQuery)); err == nil {
		limit = l
	}

	return skip, limit
}

func getUintParam(c echo.Context, key string) (uint, error) {
	param := c.Param(key)
	if param == "" {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "missing parameter "+key)
	}

	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid parameter "+key)
	}

	return uint(id), nil
}
