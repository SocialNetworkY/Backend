package post

import (
	"context"
	"github.com/SocialNetworkY/Backend/internal/report/gateway/post/http"
	"github.com/SocialNetworkY/Backend/internal/report/model"
)

type Gateway struct {
	http *http.Gateway
}

func New(httpAddr, grpcAddr string) *Gateway {
	return &Gateway{
		http: http.New(httpAddr),
	}
}

func (g *Gateway) PostInfo(ctx context.Context, postID uint) (*model.Post, error) {
	return g.http.PostInfo(ctx, postID)
}
