package v1

import (
	"context"
	"github.com/SocialNetworkY/Backend/pkg/gen"
)

type (
	PostService interface {
		DeleteByUser(userID uint) error
	}

	CommentService interface {
		DeleteByUser(userID uint) error
	}

	LikeService interface {
		DeleteByUser(userID uint) error
	}

	Handler struct {
		gen.PostServiceServer
		ps PostService
		cs CommentService
		ls LikeService
	}
)

func New(ps PostService, cs CommentService, ls LikeService) *Handler {
	return &Handler{
		ps: ps,
		cs: cs,
		ls: ls,
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

func (h *Handler) DeleteUserComments(ctx context.Context, r *gen.DeleteUserCommentsRequest) (*gen.DeleteUserCommentsResponse, error) {
	userID := uint(r.GetUserId())
	if err := h.cs.DeleteByUser(userID); err != nil {
		return nil, err
	}

	return &gen.DeleteUserCommentsResponse{
		Success: true,
	}, nil
}

func (h *Handler) DeleteUserLikes(ctx context.Context, r *gen.DeleteUserLikesRequest) (*gen.DeleteUserLikesResponse, error) {
	userID := uint(r.GetUserId())
	if err := h.ls.DeleteByUser(userID); err != nil {
		return nil, err
	}

	return &gen.DeleteUserLikesResponse{
		Success: true,
	}, nil
}
