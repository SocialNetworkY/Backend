package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
)

type Comment struct {
	client *elasticsearch.TypedClient
}

const commentIndex = "comments"

func NewComment(addr string) (*Comment, error) {
	log.Printf("Connecting to Elasticsearch for Comment: %s\n", addr)
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

	log.Printf("Connected to Elasticsearch for Comment: %s\n", addr)
	return &Comment{client: client}, nil
}

// Index creates or updates a comment in the Elasticsearch index
func (c *Comment) Index(comment *model.Comment) error {
	log.Printf("Indexing comment: %v\n", comment)
	res, err := c.client.Index(commentIndex).Id(fmt.Sprintf("%d", comment.ID)).Request(comment).Do(context.TODO())
	if err != nil {
		log.Printf("Error indexing comment: %v\n", err)
		return err
	}
	if res.Result.Name != "created" && res.Result.Name != "updated" {
		log.Printf("Failed to index comment: %v\n", res.Result.Name)
		return errors.New("failed to index comment")
	}

	log.Printf("Comment indexed successfully: %v\n", comment)
	return nil
}

func (c *Comment) Delete(id uint) error {
	log.Printf("Deleting comment with ID: %d\n", id)
	_, err := c.client.Delete(commentIndex, fmt.Sprintf("%d", id)).Do(context.TODO())
	if err != nil {
		log.Printf("Error deleting comment: %v\n", err)
		return err
	}

	log.Printf("Comment deleted successfully: %d\n", id)
	return nil
}

func (c *Comment) Search(query string, skip, limit int) ([]uint, error) {
	log.Printf("Searching comments with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	query = strings.ToLower(query)
	fields := []string{"content"}

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

	res, err := c.client.Search().Index(commentIndex).
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
		log.Printf("Error searching comments: %v\n", err)
		return nil, err
	}

	var ids []uint
	for _, hit := range res.Hits.Hits {
		comment := &model.Comment{}
		if err := json.Unmarshal(hit.Source_, comment); err != nil {
			log.Printf("Error unmarshalling comment: %v\n", err)
			return nil, err
		}
		ids = append(ids, comment.ID)
	}

	log.Printf("Comments found: %v\n", ids)
	return ids, nil
}
