package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
)

type (
	Post struct {
		client *elasticsearch.TypedClient
	}

	indexedPost struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		ImageURLs []string  `json:"image_urls"`
		VideoURLs []string  `json:"video_urls"`
		Tags      []string  `json:"tags"`
		PostedAt  time.Time `json:"posted_at"`
		Edited    bool      `json:"edited" gorm:"-"`
		EditedBy  uint      `json:"edited_by"`
		CreatedAt time.Time `json:"-"`
		UpdatedAt time.Time `json:"edited_at"`
	}
)

const postIndex = "posts"

func NewPost(addr string) (*Post, error) {
	log.Printf("Connecting to Elasticsearch for Post: %s\n", addr)
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Printf("Error connecting to Elasticsearch: %v\n", err)
		return nil, err
	}
	// Check if connection
	ok, err := client.Ping().Do(context.Background())
	if err != nil {
		log.Printf("Error pinging Elasticsearch: %v\n", err)
		return nil, err
	}
	if !ok {
		log.Printf("Connection to Elasticsearch failed\n")
		return nil, errors.New("connection failed")
	}

	log.Printf("Connected to Elasticsearch for Post: %s\n", addr)
	return &Post{client: client}, nil
}

// Index creates or updates a post in the Elasticsearch index
func (p *Post) Index(post *model.Post) error {
	log.Printf("Indexing post: %v\n", post)
	indexed := indexedPost{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		ImageURLs: post.ImageURLs,
		VideoURLs: post.VideoURLs,
		Tags:      make([]string, len(post.Tags)),
		PostedAt:  post.PostedAt,
		Edited:    post.Edited,
		EditedBy:  post.EditedBy,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
	for i, tag := range post.Tags {
		indexed.Tags[i] = tag.Name
	}

	res, err := p.client.Index(postIndex).Id(fmt.Sprintf("%d", indexed.ID)).Request(indexed).Do(context.TODO())
	if err != nil {
		log.Printf("Error indexing post: %v\n", err)
		return err
	}
	if res.Result.Name != "created" && res.Result.Name != "updated" {
		log.Printf("Failed to index post: %v\n", res.Result.Name)
		return errors.New("failed to index post")
	}

	log.Printf("Post indexed successfully: %v\n", post)
	return nil
}

func (p *Post) Delete(id uint) error {
	log.Printf("Deleting post with ID: %d\n", id)
	_, err := p.client.Delete(postIndex, fmt.Sprintf("%d", id)).Do(context.TODO())
	if err != nil {
		log.Printf("Error deleting post: %v\n", err)
		return err
	}

	log.Printf("Post deleted successfully: %d\n", id)
	return nil
}

func (p *Post) Search(query string, skip, limit int) ([]uint, error) {
	log.Printf("Searching posts with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	query = strings.ToLower(query)
	fields := []string{"title", "content", "tags"}

	var should []types.Query

	should = append(should, types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query:  query,
			Fields: fields,
			Type: &textquerytype.TextQueryType{
				Name: "phrase_prefix",
			},
		},
	})

	res, err := p.client.Search().Index(postIndex).
		Request(&search.Request{
			From: &skip,
			Size: &limit,
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Should: should,
				},
			},
		}).
		Do(context.TODO())
	if err != nil {
		log.Printf("Error searching posts: %v\n", err)
		return nil, err
	}

	var ids []uint
	for _, hit := range res.Hits.Hits {
		post := &indexedPost{}
		if err := json.Unmarshal(hit.Source_, post); err != nil {
			log.Printf("Error unmarshalling post: %v\n", err)
			return nil, err
		}
		ids = append(ids, post.ID)
	}

	log.Printf("Posts found: %v\n", ids)
	return ids, nil
}
