package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"golang.org/x/net/context"
	"log"
	"strings"
)

type (
	Ban struct {
		client *elasticsearch.TypedClient
	}
)

const banIndex = "bans"

func NewBan(addr string) (*Ban, error) {
	log.Printf("Connecting to Elasticsearch for User: %s\n", addr)
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}
	// Check if connection
	ok, err := client.Ping().Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("connection failed")
	}

	log.Printf("Connected to Elasticsearch for User: %s\n", addr)
	return &Ban{client: client}, nil
}

// Index creates or updates a ban in the Elasticsearch index
func (b *Ban) Index(ban *model.Ban) error {
	res, err := b.client.Index(banIndex).Id(fmt.Sprintf("%d", ban.ID)).Request(ban).Do(context.TODO())
	if err != nil {
		return err
	}
	if res.Result.Name != "created" && res.Result.Name != "updated" {
		return errors.New("failed to index ban")
	}

	return nil
}

func (b *Ban) Delete(id uint) error {
	_, err := b.client.Delete(banIndex, fmt.Sprintf("%d", id)).Do(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (b *Ban) Search(query string, skip, limit int) ([]uint, error) {
	query = strings.ToLower(query)
	fields := []string{"banReason", "unbanReason"}

	should := make([]types.Query, 0, len(fields))

	should = append(should, types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query:  query,
			Fields: fields,
			Type: &textquerytype.TextQueryType{
				Name: "phrase_prefix",
			},
		},
	})

	// Example of prefix query
	/*for _, field := range fields {
		prefix := make(map[string]types.PrefixQuery, 1)
		prefix[field] = types.PrefixQuery{
			Value: query,
		}

		should = append(should, types.Query{
			Prefix: prefix,
		})
	}*/

	res, err := b.client.Search().Index(banIndex).
		From(skip).Size(limit).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Should: should,
				},
			},
		}).
		Do(context.TODO())
	if err != nil {
		return nil, err
	}

	var bansIDs []uint
	for _, hit := range res.Hits.Hits {
		ban := &model.Ban{}
		if err := json.Unmarshal(hit.Source_, ban); err != nil {
			return nil, err
		}
		bansIDs = append(bansIDs, ban.ID)
	}

	return bansIDs, nil
}
