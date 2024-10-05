package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/user/model"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	authorizationMetadata = "authorization"
)

func (h *Handler) getRequesterFromMetadata(ctx context.Context) (*model.User, error) {
	requesterID, err := h.getRequesterIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return h.us.Find(requesterID)
}

func (h *Handler) getRequesterIDFromMetadata(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(http.StatusUnauthorized, "missing metadata")
	}

	authHeader, ok := md[authorizationMetadata]
	if !ok || len(authHeader) == 0 {
		return 0, status.Errorf(http.StatusUnauthorized, "missing Authorization header")
	}

	token := authHeader[0]
	requesterID, err := h.ag.Authenticate(ctx, token)
	if err != nil {
		return 0, status.Errorf(http.StatusUnauthorized, err.Error())
	}

	return requesterID, nil
}
