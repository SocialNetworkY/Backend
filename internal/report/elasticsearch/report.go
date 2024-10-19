package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/SocialNetworkY/Backend/internal/report/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
)

type Report struct {
	client *elasticsearch.TypedClient
}

const reportIndex = "reports"

func NewReport(addr string) (*Report, error) {
	log.Printf("Connecting to Elasticsearch for Report: %s\n", addr)
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

	log.Printf("Connected to Elasticsearch for Report: %s\n", addr)
	return &Report{client: client}, nil
}

// Index creates or updates a report in the Elasticsearch index
func (r *Report) Index(report *model.Report) error {
	log.Printf("Indexing report: %v\n", report)
	res, err := r.client.Index(reportIndex).Id(fmt.Sprintf("%d", report.ID)).Request(report).Do(context.TODO())
	if err != nil {
		log.Printf("Error indexing report: %v\n", err)
		return err
	}
	if res.Result.Name != "created" && res.Result.Name != "updated" {
		log.Printf("Failed to index report: %v\n", res.Result.Name)
		return errors.New("failed to index report")
	}

	log.Printf("Report indexed successfully: %v\n", report)
	return nil
}

func (r *Report) Delete(id uint) error {
	log.Printf("Deleting report: %d\n", id)
	res, err := r.client.Delete(reportIndex, fmt.Sprintf("%d", id)).Do(context.TODO())
	if err != nil {
		log.Printf("Error deleting report: %v\n", err)
		return err
	}
	if res.Result.Name != "deleted" {
		log.Printf("Failed to delete report: %v\n", res.Result.Name)
		return errors.New("failed to delete report")
	}

	log.Printf("Report deleted successfully: %d\n", id)
	return nil
}

func (r *Report) Search(query string, skip, limit int) ([]uint, error) {
	log.Printf("Searching reports with query: %s, skip: %d, limit: %d\n", query, skip, limit)
	query = strings.ToLower(query)
	fields := []string{"reason", "answer", "status"}

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

	res, err := r.client.Search().Index(reportIndex).
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
		log.Printf("Error searching reports: %v\n", err)
		return nil, err
	}

	var ids []uint
	for _, hit := range res.Hits.Hits {
		report := &model.Report{}
		if err := json.Unmarshal(hit.Source_, report); err != nil {
			log.Printf("Error unmarshalling report: %v\n", err)
			return nil, err
		}
		ids = append(ids, report.ID)
	}

	log.Printf("Reports found: %v\n", ids)
	return ids, nil
}
