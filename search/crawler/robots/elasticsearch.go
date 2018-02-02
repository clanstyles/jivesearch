package robots

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
)

// ElasticSearch holds our Elasticsearch connection and index information
type ElasticSearch struct {
	Client *elastic.Client
	Index  string
	Type   string
	Bulk   *elastic.BulkProcessor
}

// Get retrieves a single cached robots.txt file
func (e *ElasticSearch) Get(sh string) (*Robots, error) {
	rob := &Robots{SchemeHost: sh}

	res, err := e.Client.Get().Index(e.Index).Type(e.Type).Id(sh).Do(context.TODO())
	if err != nil {
		if elastic.IsNotFound(err) {
			err = nil
		}
		return rob, err
	}

	if err := json.Unmarshal(*res.Source, rob); err != nil {
		return rob, err
	}

	if res.Found {
		rob.Cached = true
	}

	return rob, nil
}

// Put caches a robots.txt file
func (e *ElasticSearch) Put(r *Robots) {
	item := elastic.NewBulkIndexRequest().
		Index(e.Index).
		Type(e.Type).
		Id(r.SchemeHost).
		Doc(r)

	e.Bulk.Add(item)
}

// Mapping is the mapping of our robots Index
func (e *ElasticSearch) Mapping() string {
	return fmt.Sprintf(`{
    "mappings": {
			"%v": {
				"dynamic": "strict",
        "properties": {
					"body": {
            "type": "text",
    				"index": false
          },
					"status": {
            "type": "short"
          },
					"expires": {
						"type": "date",
						"format": "yyyyMMddHHmm"
          }
        }
      }
    }
  }`, e.Type)
}

// Setup creates an index for caching robots.txt files
func (e *ElasticSearch) Setup() error {
	_, err := e.Client.CreateIndex(e.Index).Body(e.Mapping()).Do(context.TODO())
	return err
}

// IndexExists returns true if the index exists
func (e *ElasticSearch) IndexExists() (bool, error) {
	return e.Client.IndexExists(e.Index).Do(context.TODO())
}
