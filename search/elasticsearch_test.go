package search

import (
	"jivesearch/search/document"
	"jivesearch/search/vote"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/olivere/elastic"
	"golang.org/x/text/language"
)

func TestFetch(t *testing.T) {
	type want struct {
		*Results
		err error
	}

	for _, c := range []struct {
		name   string
		query  string
		lang   language.Tag
		region language.Region
		number int
		page   int
		votes  []vote.Result
		status int
		resp   string
		want
	}{
		{
			name:   "basic",
			query:  "Bob Dylan",
			lang:   language.English,
			region: language.MustParseRegion("US"),
			number: 25,
			page:   1,
			status: http.StatusOK,
			resp: `{
			  "took": 5,
			  "timed_out": false,
			  "_shards": {
			    "total": 5,
			    "successful": 5,
			    "failed": 0
			  },
			  "hits": {
			    "total": 2,
			    "max_score": 6.6914043,
			    "hits": [
					{
						"_index": "search-english",
						"_type": "search",
						"_id": "http://example.com/articles/is-bob-dylan-literature-1476401068",
						"_score": 5.6856093,
						"_source": {
						"title": "Is Bob Dylan Literature? - WSJ",
						"description": "The Nobel committee says ‘Yes’ to Bob Dylan."
					}
			      },
					{
						"_index": "search-english",
						"_type": "search",
						"_id": "http://www.example.com/book-search/author/DYLAN",
						"_score": 6.6914043,
						"_source": {
							"title": "Dylan, Bob - example",
							"description": "It's Easy to Play Bob Dylan by Dylan"
						}
					}
			    ]
			  }
			}`,
			want: want{
				&Results{
					Count: 2,
					Documents: []*document.Document{
						&document.Document{
							ID: "http://example.com/articles/is-bob-dylan-literature-1476401068",
							Content: document.Content{
								Title:       "Is Bob Dylan Literature? - WSJ",
								Description: "The Nobel committee says ‘Yes’ to Bob Dylan.",
							},
						},
						&document.Document{
							ID: "http://www.example.com/book-search/author/DYLAN",
							Content: document.Content{
								Title:       "Dylan, Bob - example",
								Description: "It's Easy to Play Bob Dylan by Dylan",
							},
						},
					},
				},
				nil,
			},
		},
		{
			name:   "language",
			query:  "jimi hendrix",
			lang:   language.BrazilianPortuguese,
			region: language.MustParseRegion("BR"),
			number: 2500,
			page:   5,
			status: http.StatusOK,
			resp: `{
			  "took": 5,
			  "timed_out": false,
			  "_shards": {
			    "total": 5,
			    "successful": 5,
			    "failed": 0
			  },
			  "hits": {
			    "total": 2500,
			    "max_score": 6.6914043,
			    "hits": [
					{
						"_index": "search-brazilian",
						"_type": "search",
						"_id": "http://example.com.br/articles/is-bob-dylan-literature-1476401068",
						"_score": 5.6856093,
						"_source": {
						"title": "Is Bob Dylan Literature? - WSJ",
						"description": "The Nobel committee says ‘Yes’ to Bob Dylan."
					}
			      },
					{
						"_index": "search-brazilian",
						"_type": "search",
						"_id": "http://www.example.com.br/book-search/author/DYLAN",
						"_score": 6.6914043,
						"_source": {
							"title": "Dylan, Bob - example",
							"description": "It's Easy to Play Bob Dylan by Dylan"
						}
					}
			    ]
			  }
			}`,
			want: want{
				&Results{
					Count: 2500,
					Documents: []*document.Document{
						&document.Document{
							ID: "http://example.com.br/articles/is-bob-dylan-literature-1476401068",
							Content: document.Content{
								Title:       "Is Bob Dylan Literature? - WSJ",
								Description: "The Nobel committee says ‘Yes’ to Bob Dylan.",
							},
						},
						&document.Document{
							ID: "http://www.example.com.br/book-search/author/DYLAN",
							Content: document.Content{
								Title:       "Dylan, Bob - example",
								Description: "It's Easy to Play Bob Dylan by Dylan",
							},
						},
					},
				},
				nil,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			handler := http.NotFound
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				handler(w, r)
			}))
			defer ts.Close()

			handler = func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(c.resp))
			}

			e, err := MockService(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			got, err := e.Fetch(c.query, c.lang, c.region, c.number, c.page, c.votes)
			if err != c.want.err {
				t.Fatalf("got err %q; want %q", err, c.want.err)
			}

			got.Documents = []*document.Document{}
			c.want.Results.Documents = []*document.Document{}

			if !reflect.DeepEqual(got, c.want.Results) {
				t.Fatalf("got %+v; want %+v", got, c.want.Results)
			}
		})
	}
}

func MockService(url string) (*ElasticSearch, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	e := &ElasticSearch{
		ElasticSearch: &document.ElasticSearch{
			Client: client,
			Index:  "search",
		},
	}

	return e, nil
}
