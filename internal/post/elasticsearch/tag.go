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

type (
	Tag struct {
		client *elasticsearch.TypedClient
	}
)

const tagIndex = "tags"

func NewTag(addr string) (*Tag, error) {
	log.Printf("Connecting to Elasticsearch for Tag: %s\n", addr)
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

	log.Printf("Connected to Elasticsearch for Tag: %s\n", addr)
	return &Tag{client: client}, nil
}

// Index creates or updates a tag in the Elasticsearch index
func (t *Tag) Index(tag *model.Tag) error {
	log.Printf("Indexing tag: %v\n", tag)
	res, err := t.client.Index(tagIndex).Id(fmt.Sprintf("%d", tag.ID)).Request(tag).Do(context.TODO())
	if err != nil {
		log.Printf("Error indexing tag: %v\n", err)
		return err
	}
	if res.Result.Name != "created" && res.Result.Name != "updated" {
		log.Printf("Failed to index tag: %v\n", res.Result.Name)
		return errors.New("failed to index tag")
	}

	log.Printf("Tag indexed successfully: %v\n", tag)
	return nil
}

func (t *Tag) Delete(id uint) error {
	log.Printf("Deleting tag with ID: %d\n", id)
	_, err := t.client.Delete(tagIndex, fmt.Sprintf("%d", id)).Do(context.TODO())
	if err != nil {
		log.Printf("Error deleting tag: %v\n", err)
		return err
	}

	log.Printf("Tag deleted successfully: %d\n", id)
	return nil
}

func (t *Tag) Search(query string, skip, limit int) ([]uint, error) {
	log.Printf("Searching tags with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	query = strings.ToLower(query)
	fields := []string{"name"}

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

	res, err := t.client.Search().Index(tagIndex).
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
		log.Printf("Error searching tags: %v\n", err)
		return nil, err
	}

	var ids []uint
	for _, hit := range res.Hits.Hits {
		tag := &model.Tag{}
		if err := json.Unmarshal(hit.Source_, tag); err != nil {
			log.Printf("Error unmarshalling tag: %v\n", err)
			return nil, err
		}
		ids = append(ids, tag.ID)
	}

	log.Printf("Tags found: %v\n", ids)
	return ids, nil
}
