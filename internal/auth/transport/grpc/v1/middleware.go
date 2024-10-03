package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	authorizationMetadata = "authorization"
)

func (h *Handler) getUserFromMetadata(ctx context.Context) (*model.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(http.StatusUnauthorized, "missing metadata")
	}

	authHeader, ok := md[constant.GRPCAuthorizationMetadata]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(http.StatusUnauthorized, "missing Authorization header")
	}

	token := authHeader[0]
	return h.as.Auth(token)
}
