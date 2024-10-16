package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/pkg/gen"
)

type (
	ReportService interface {
		DeleteByUser(userID uint) error
		DeleteByPost(postID uint) error
	}

	Handler struct {
		gen.ReportServiceServer
		rs ReportService
	}
)

func New(rs ReportService) *Handler {
	return &Handler{
		rs: rs,
	}
}

func (h *Handler) DeleteUserReports(ctx context.Context, r *gen.DeleteUserReportsRequest) (*gen.DeleteUserReportsResponse, error) {
	userID := uint(r.GetUserId())
	if err := h.rs.DeleteByUser(userID); err != nil {
		return nil, err
	}

	return &gen.DeleteUserReportsResponse{
		Success: true,
	}, nil
}

func (h *Handler) DeletePostReports(ctx context.Context, r *gen.DeletePostReportsRequest) (*gen.DeletePostReportsResponse, error) {
	postID := uint(r.GetPostId())
	if err := h.rs.DeleteByPost(postID); err != nil {
		return nil, err
	}

	return &gen.DeletePostReportsResponse{
		Success: true,
	}, nil
}
