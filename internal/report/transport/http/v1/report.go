package v1

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/report/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initReportApi(api *echo.Group) {
	reports := api.Group("/reports", h.authMiddleware, h.banMiddleware)
	{
		reports.POST("", h.postReport)
		report := reports.Group(fmt.Sprintf("/:%s", reportIDParam), h.setReportByIDMiddleware, h.checkAccessMiddleware)
		{
			report.GET("", h.get)
			report.DELETE("", h.delete, h.adminMiddleware)
			report.POST("/answer", h.answer, h.adminMiddleware)
			report.POST("/reject", h.reject, h.adminMiddleware)
		}
	}
}

func (h *Handler) postReport(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
	}

	req := &struct {
		PostID uint   `form:"post_id" validate:"required"`
		Reason string `form:"reason" validate:"required,min=1,max=255"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	report, err := h.rs.Create(requester.ID, req.PostID, req.Reason)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, report)
}

func (h *Handler) get(c echo.Context) error {
	report, ok := c.Get(reportLocals).(*model.Report)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get report")
	}

	return c.JSON(http.StatusOK, report)
}

func (h *Handler) delete(c echo.Context) error {
	report, ok := c.Get(reportLocals).(*model.Report)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get report")
	}

	if err := h.rs.Delete(report.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) answer(c echo.Context) error {
	report, ok := c.Get(reportLocals).(*model.Report)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get report")
	}

	if report.Closed {
		return echo.NewHTTPError(http.StatusForbidden, "report is closed")
	}

	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
	}

	req := &struct {
		Answer string `form:"answer" validate:"required,min=1,max=255"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	report, err := h.rs.Answer(report.ID, requester.ID, req.Answer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, report)
}

func (h *Handler) reject(c echo.Context) error {
	report, ok := c.Get(reportLocals).(*model.Report)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get report")
	}

	if report.Closed {
		return echo.NewHTTPError(http.StatusForbidden, "report is closed")
	}

	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get requester")
	}

	req := &struct {
		Answer string `form:"answer" validate:"required,min=1,max=255"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	report, err := h.rs.Reject(report.ID, requester.ID, req.Answer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, report)
}
