package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/SocialNetworkY/Backend/internal/user/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"golang.org/x/net/context"
)

type (
	User struct {
		client *elasticsearch.TypedClient
	}
)

const userIndex = "users"

func NewUser(addr string) (*User, error) {
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
	return &User{client: client}, nil
}

// Index creates or updates a user in the Elasticsearch index
func (u *User) Index(user *model.User) error {
	res, err := u.client.Index(userIndex).Id(fmt.Sprintf("%d", user.ID)).Request(user).Do(context.TODO())
	if err != nil {
		return err
	}
	if res.Result.Name != "created" && res.Result.Name != "updated" {
		return errors.New("failed to index user")
	}

	return nil
}

func (u *User) Delete(id uint) error {
	_, err := u.client.Delete(userIndex, fmt.Sprintf("%d", id)).Do(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Search(query string, skip, limit int) ([]uint, error) {
	query = strings.ToLower(query)
	fields := []string{"username", "nickname"}

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

	res, err := u.client.Search().Index(userIndex).
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
		return nil, err
	}

	var usersIDs []uint
	for _, hit := range res.Hits.Hits {
		user := &model.User{}
		if err := json.Unmarshal(hit.Source_, user); err != nil {
			return nil, err
		}
		usersIDs = append(usersIDs, user.ID)
	}

	return usersIDs, nil
}
