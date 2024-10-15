package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/pkg/gen"
)

type (
	PostService interface {
		DeleteByUser(userID uint) error
	}

	Handler struct {
		gen.PostServiceServer
		ps PostService
	}
)

func New(ps PostService) *Handler {
	return &Handler{
		ps: ps,
	}
}

func (h *Handler) DeleteUserPosts(ctx context.Context, r *gen.DeleteUserPostsRequest) (*gen.DeleteUserPostsResponse, error) {
	userID := uint(r.GetUserId())
	if err := h.ps.DeleteByUser(userID); err != nil {
		return nil, err
	}

	return &gen.DeleteUserPostsResponse{
		Success: true,
	}, nil
}
