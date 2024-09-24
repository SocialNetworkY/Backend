package v1

import (
	"context"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/model"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
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
	log.Printf("md: %+v", md)

	authHeader, ok := md[authorizationMetadata]
	log.Printf("authHeader: %v", authHeader)
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(http.StatusUnauthorized, "missing Authorization header")
	}

	token := authHeader[0]
	return h.as.Auth(token)
}
