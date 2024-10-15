package v1

import (
	"github.com/SocialNetworkY/Backend/internal/report/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initReportApi(api *echo.Group) {
	api.POST("", h.postReport, h.authMiddleware, h.banMiddleware)
	api.GET("/:id", h.get)
	api.POST("/:id/answer", h.answer)
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

}
