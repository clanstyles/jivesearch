package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jivesearch/jivesearch/search/document"

	"github.com/olivere/elastic"
)

func TestUpsert(t *testing.T) {
	for _, c := range []struct {
		name   string
		status int
		resp   string
		doc    *document.Document
		err    error
	}{
		{
			name:   "basic",
			status: http.StatusCreated,
			resp: `{
			  "took": 27,
			  "errors": false,
			  "items": [
					{
			      "create": {
			        "_index": "search",
			        "_type": "document",
			        "_id": "AVhRlxyshqP4iSOLLnUz",
			        "_version": 1,
			        "_shards": {
			          "total": 2,
			          "successful": 1,
			          "failed": 0
			        },
			        "status": 201
			      }
				  }
				]
			}`,
			doc: &document.Document{
				ID:     "http://www.example.com/path/to/nowhere",
				Domain: "example.com",
				Host:   "http://www.example.com",
			},
			err: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			handler := http.NotFound
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler(w, r)
			}))
			defer ts.Close()

			handler = func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				w.Write([]byte(c.resp))
			}

			e, err := MockService(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			e.Upsert(c.doc)

			if err := e.Bulk.Flush(); err != nil {
				t.Fatal(err)
			}

			stats := e.Bulk.Stats()
			if stats.Succeeded != 1 {
				t.Fatalf("upsert failed: got %d", stats.Succeeded)
			}
		})
	}
}

func TestCrawledAndCount(t *testing.T) {
	type want struct {
		count   int
		crawled time.Time
		err     error
	}

	for _, c := range []struct {
		name   string
		url    string
		domain string
		status int
		resp   string
		want
	}{
		{
			name:   "exists",
			url:    "http://www.example.com/path/to/somewhere",
			domain: "example.com",
			status: http.StatusOK,
			resp: `{
				"responses": [
					{
						"took": 1,
						"timed_out": false,
						"_shards": {
							"total": 5,
							"successful": 5,
							"failed": 0
						},
						"hits": {
							"total": 593,
							"max_score": 0,
							"hits": [
								{
									"_index": "test-search-english",
									"_type": "document",
									"_id": "https://ast.wikipedia.org/wiki/RSS",
									"_score": 0,
									"_source": {
										"scheme": "https",
										"host": "ast.wikipedia.org",
										"domain": "wikipedia.org",
										"tld": "org",
										"path_parts": "wiki RSS",
										"crawled": "20170909",
										"mime": "text/html",
										"status": 200,
										"canonical": true,
										"title": "RSS - Wikipedia",
										"index": true
									}
								},
								{
									"_index": "test-search-english",
									"_type": "document",
									"_id": "https://he.wikipedia.org/wiki/RSS",
									"_score": 0,
									"_source": {
										"scheme": "https",
										"host": "he.wikipedia.org",
										"domain": "wikipedia.org",
										"tld": "org",
										"path_parts": "wiki RSS",
										"crawled": "20170909",
										"mime": "text/html",
										"status": 200,
										"canonical": true,
										"title": "RSS – Wikipedia",
										"index": true
									}
								}
							]
						},
						"status": 200
					},
					{
						"took": 21,
						"timed_out": false,
						"_shards": {
							"total": 5,
							"successful": 5,
							"failed": 0
						},
						"hits": {
							"total": 1,
							"max_score": 9.395768,
							"hits": [
							{
								"_index": "search-english",
								"_type": "document",
								"_id": "http://www.example.com/path/to/somewhere",
								"_score": 9.395768,
								"_source": {
								"crawled": "20170706"
								}
							}
							]
						}
					}
				]
			}`,
			want: want{593, time.Date(2017, time.July, 06, 0, 0, 0, 0, time.UTC), nil},
		},
		{
			name:   "does not exist",
			url:    "http://www.example.com/path/to/nowhere",
			domain: "example.com",
			resp: `{
				"responses": [
					{
						"took": 1,
						"timed_out": false,
						"_shards": {
							"total": 5,
							"successful": 5,
							"failed": 0
						},
						"hits": {
							"total": 412,
							"max_score": 0,
							"hits": [
								{
									"_index": "test-search-english",
									"_type": "document",
									"_id": "https://ast.wikipedia.org/wiki/RSS",
									"_score": 0,
									"_source": {
										"scheme": "https",
										"host": "ast.wikipedia.org",
										"domain": "wikipedia.org",
										"tld": "org",
										"path_parts": "wiki RSS",
										"crawled": "20170909",
										"mime": "text/html",
										"status": 200,
										"canonical": true,
										"title": "RSS - Wikipedia",
										"index": true
									}
								},
								{
									"_index": "test-search-english",
									"_type": "document",
									"_id": "https://he.wikipedia.org/wiki/RSS",
									"_score": 0,
									"_source": {
										"scheme": "https",
										"host": "he.wikipedia.org",
										"domain": "wikipedia.org",
										"tld": "org",
										"path_parts": "wiki RSS",
										"crawled": "20170909",
										"mime": "text/html",
										"status": 200,
										"canonical": true,
										"title": "RSS – Wikipedia",
										"index": true
									}
								}
							]
						},
						"status": 200
					},
					{
						"took": 0,
						"timed_out": false,
						"_shards": {
							"total": 5,
							"successful": 5,
							"failed": 0
						},
						"hits": {
							"total": 0,
							"max_score": null,
							"hits": []
						},
						"status": 200
					}
				]
			}`,
			status: http.StatusOK,
			want:   want{412, time.Time{}, nil},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				w.Write([]byte(c.resp))
			}))

			defer ts.Close()

			e, err := MockService(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			crawled, count, err := e.CrawledAndCount(c.url, c.domain)
			if err != c.want.err {
				t.Fatalf("got err %q; want %q", err, c.want.err)
			}

			if count != c.want.count {
				t.Fatalf("got %d; want %d", count, c.want.count)
			}

			if crawled != c.want.crawled {
				t.Fatalf("got %q; want %q", crawled, c.want.crawled)
			}
		})
	}
}

func MockService(url string) (*ElasticSearch, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	bulk, err := client.BulkProcessor().Stats(true).Do(context.TODO())
	if err != nil {
		return nil, err
	}

	return &ElasticSearch{
		ElasticSearch: &document.ElasticSearch{
			Client: client, Index: "search", Type: "document",
		},
		Bulk: bulk,
	}, nil
}
