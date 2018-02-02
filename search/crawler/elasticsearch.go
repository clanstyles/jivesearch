package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"jivesearch/search/document"
	"sync"
	"time"

	"github.com/olivere/elastic"
)

// ElasticSearch satisfies the crawler's Backend interface
type ElasticSearch struct {
	*document.ElasticSearch
	Bulk *elastic.BulkProcessor
	sync.Mutex
}

// Upsert updates a document or inserts it if it doesn't exist
// NOTE: Elasticsearch has a 512-byte limit on an insert operation.
// Upsert does not have that limit.
func (e *ElasticSearch) Upsert(doc *document.Document) error {
	a, err := e.Analyzer(doc.Language)
	if err != nil {
		return err
	}

	idx := e.IndexName(a)

	item := elastic.NewBulkUpdateRequest().
		Index(idx).
		Type(e.Type).
		Id(doc.ID).
		DocAsUpsert(true).
		Doc(doc)

	e.Bulk.Add(item)
	return nil
}

// CrawledAndCount returns the crawled date of the url (if any) and
// the total number of links a domain has
func (e *ElasticSearch) CrawledAndCount(u, domain string) (time.Time, int, error) {
	body := fmt.Sprintf(`{
		"bool": {
			"filter": [
				{
					"term": {
						"domain": "%v"
					}
				},
				{
					"term": {
						"index": "true"
					}
				}
			]
		}
	}`, domain)

	var crawled, cnt = time.Time{}, 0

	// even though this technically could be a count request
	// it s/b faster using multisearch.
	countReq := elastic.NewSearchRequest().
		Index(e.Index + "-*").
		Type(e.Type).Source(elastic.NewSearchSource().
		Query(elastic.RawStringQuery(body)),
	)

	crawledRequest := elastic.NewSearchRequest().
		Index(e.Index + "-*").
		Type(e.Type).
		Source(elastic.NewSearchSource().
			Query(elastic.NewTermQuery("_id", u)).
			FetchSourceContext(elastic.NewFetchSourceContext(true).Include("crawled")),
		)

	// Concurrently calling this results in Error 429 [reduce_search_phase_exception] error.
	// Seems to only happen when searching multiple indices (e.g. search-*) as it
	// doesn't happen when searching one at a time (e.g. search-english, etc...)
	e.Lock()

	res, err := e.Client.MultiSearch().
		Add(countReq, crawledRequest).
		Do(context.TODO())

	e.Unlock()

	if err != nil {
		return crawled, cnt, err
	}

	r1, r2 := res.Responses[0], res.Responses[1]

	cnt = int(r1.TotalHits())

	if err != nil && !elastic.IsNotFound(r2.Error) {
		return crawled, cnt, fmt.Errorf(r2.Error.Reason)
	}

	for _, h := range r2.Hits.Hits {
		c := make(map[string]string)
		if err := json.Unmarshal(*h.Source, &c); err != nil {
			return crawled, cnt, err
		}
		crawled, err = time.Parse("20060102", c["crawled"])
	}

	return crawled, cnt, err
}
