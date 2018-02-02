package suggest

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/olivere/elastic"
)

func TestCompletion(t *testing.T) {
	for _, c := range []struct {
		term   string
		size   int
		status int
		resp   string
		want   Results
	}{
		{
			term:   "b",
			size:   10,
			status: http.StatusOK,
			resp: `{
				"took": 0,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"skipped": 0,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": 0,
					"hits": []
				},
				"suggest": {
					"completion_suggest": [
						{
							"text": "b",
							"offset": 0,
							"length": 1,
							"options": [
								{
									"text": "brad",
									"_index": "test-queries",
									"_type": "query",
									"_id": "brad",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "brad",
											"weight": 1
										}
									}
								},
								{
									"text": "bros",
									"_index": "test-queries",
									"_type": "query",
									"_id": "bros",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "bros",
											"weight": 1
										}
									}
								},
								{
									"text": "bob",
									"_index": "test-queries",
									"_type": "query",
									"_id": "bob",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "bob",
											"weight": 1
										}
									}
								},
								{
									"text": "blondie",
									"_index": "test-queries",
									"_type": "query",
									"_id": "blondie",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "blondie",
											"weight": 1
										}
									}
								},
								{
									"text": "brad pitt",
									"_index": "test-queries",
									"_type": "query",
									"_id": "brad pitt",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "brad pitt",
											"weight": 1
										}
									}
								},
								{
									"text": "buster",
									"_index": "test-queries",
									"_type": "query",
									"_id": "buster",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "buster",
											"weight": 1
										}
									}
								}
							]
						}
					]
				}
			}`,
			want: Results{
				RawQuery:    `{"completion_suggest":{"text":"b","completion":{"field":"completion_suggest","size":10}}}`,
				Suggestions: []string{"brad", "bros", "bob", "blondie", "brad pitt", "buster"},
			},
		},
		{
			term:   "ji",
			size:   7,
			status: http.StatusOK,
			resp: `{
				"took": 0,
				"timed_out": false,
				"_shards": {
					"total": 5,
					"successful": 5,
					"skipped": 0,
					"failed": 0
				},
				"hits": {
					"total": 0,
					"max_score": 0,
					"hits": []
				},
				"suggest": {
					"completion_suggest": [
						{
							"text": "ji",
							"offset": 0,
							"length": 1,
							"options": [
								{
									"text": "jiffy",
									"_index": "test-queries",
									"_type": "query",
									"_id": "jiffy",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "jiffy",
											"weight": 1
										}
									}
								},
								{
									"text": "jill",
									"_index": "test-queries",
									"_type": "query",
									"_id": "jill",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "jill",
											"weight": 1
										}
									}
								},
								{
									"text": "jim",
									"_index": "test-queries",
									"_type": "query",
									"_id": "jim",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "jim",
											"weight": 1
										}
									}
								},
								{
									"text": "jimi",
									"_index": "test-queries",
									"_type": "query",
									"_id": "jimi",
									"_score": 1,
									"_source": {
										"completion_suggest": {
											"input": "jimi",
											"weight": 1
										}
									}
								}
							]
						}
					]
				}
			}`,
			want: Results{
				RawQuery:    `{"completion_suggest":{"text":"ji","completion":{"field":"completion_suggest","size":7}}}`,
				Suggestions: []string{"jiffy", "jill", "jim", "jimi"},
			},
		},
	} {
		t.Run(c.term, func(t *testing.T) {
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

			got, err := e.Completion(c.term, c.size)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	for _, c := range []struct {
		term   string
		status int
		resp   string
		want   bool
	}{
		{
			term:   "a search term",
			status: http.StatusNotFound,
			resp:   `{"exists": false}`,
			want:   false,
		},
		{
			term:   "another search term",
			status: http.StatusOK,
			resp:   `{"exists": true}`,
			want:   true,
		},
	} {
		t.Run(c.term, func(t *testing.T) {
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

			got, err := e.Exists(c.term)
			if err != nil {
				t.Fatal(err)
			}

			if got != c.want {
				t.Fatalf("got %v; want %v", got, c.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	for _, c := range []struct {
		term   string
		status int
		resp   string
	}{
		{
			term:   "a search term",
			status: http.StatusOK,
			resp:   `{"acknowledged": true}`,
		},
		{
			term:   "another search term",
			status: http.StatusOK,
			resp:   `{"acknowledged": true}`,
		},
	} {
		t.Run(c.term, func(t *testing.T) {
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

			if err := e.Insert(c.term); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIncrement(t *testing.T) {
	for _, c := range []struct {
		term   string
		status int
		resp   string
	}{
		{
			term:   "a search term",
			status: http.StatusOK,
			resp:   `{"acknowledged": true}`,
		},
		{
			term:   "another search term",
			status: http.StatusOK,
			resp:   `{"acknowledged": true}`,
		},
	} {
		t.Run(c.term, func(t *testing.T) {
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

			if err := e.Increment(c.term); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIndexExists(t *testing.T) {
	for _, c := range []struct {
		name   string
		status int
		want   bool
	}{
		{
			"exists", http.StatusOK, true,
		},
		{
			"doesn't exist", http.StatusNotFound, false,
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
			}

			e, err := MockService(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			got, err := e.IndexExists()
			if err != nil {
				t.Fatal(err)
			}

			if got != c.want {
				t.Fatalf("got %v; want %v", got, c.want)
			}
		})
	}
}

func TestSetup(t *testing.T) {
	for _, c := range []struct {
		name   string
		status int
		resp   string
	}{
		{
			"ok", http.StatusOK, `{"acknowledged": true}`,
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

			if err := e.Setup(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func MockService(url string) (*ElasticSearch, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &ElasticSearch{Client: client, Index: "queries", Type: "query"}, nil
}
