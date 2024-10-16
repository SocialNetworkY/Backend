package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/report/model"
	"net/http"
)

type Gateway struct {
	client *http.Client
	addr   string
}

func New(addr string) *Gateway {
	return &Gateway{
		addr:   addr,
		client: &http.Client{},
	}
}

func (g *Gateway) PostInfo(ctx context.Context, postID uint) (*model.Post, error) {
	url := fmt.Sprintf("%s/api/v1/posts/%d", g.addr, postID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var post *model.Post
	if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
		return nil, err
	}

	return post, nil
}
