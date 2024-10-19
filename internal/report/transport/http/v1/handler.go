package v1

import (
	"context"
	"log"

	"github.com/SocialNetworkY/Backend/internal/report/model"
	"github.com/labstack/echo/v4"
)

type (
	ReportService interface {
		Create(userID, postID uint, reason string) (*model.Report, error)
		Answer(reportID, adminID uint, answer string) (*model.Report, error)
		Reject(reportID, adminID uint, answer string) (*model.Report, error)
		Get(reportID uint) (*model.Report, error)
		GetSome(skip, limit int, status string) ([]*model.Report, error)
		GetByPost(postID uint, skip, limit int, status string) ([]*model.Report, error)
		GetByUser(userID uint, skip, limit int, status string) ([]*model.Report, error)
		GetByAdmin(adminID uint, skip, limit int, status string) ([]*model.Report, error)
		Delete(reportID uint) error
		Search(query string, skip, limit int) ([]*model.Report, error)
	}

	AuthGateway interface {
		Authenticate(ctx context.Context, auth string) (uint, error)
	}

	UserGateway interface {
		UserInfo(ctx context.Context, userID uint) (*model.User, error)
	}

	Handler struct {
		rs ReportService
		ag AuthGateway
		ug UserGateway
	}
)

func New(rs ReportService, ag AuthGateway, ug UserGateway) *Handler {
	return &Handler{
		rs: rs,
		ag: ag,
		ug: ug,
	}
}

func (h *Handler) Init(api *echo.Group) {
	log.Println("Initializing V1 api")
	v1 := api.Group("/v1")
	{
		h.initReportApi(v1)
	}
}
