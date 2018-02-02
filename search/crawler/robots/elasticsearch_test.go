package robots

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/olivere/elastic"
)

func TestGet(t *testing.T) {
	for _, c := range []struct {
		name   string
		host   string
		status int
		resp   string
		want   *Robots
	}{
		{
			name:   "cached",
			host:   "https://www.example.com",
			status: http.StatusOK,
			resp: `{
			  "_index": "robots",
			  "_type": "robots",
			  "_id": "http://www.bloomberg.com",
			  "_version": 1,
			  "found": true,
			  "_source": {
			    "status": 200,
			    "body": "User-agent: *\nDisallow: /",
			    "expires": "201611071423"
			  }
			}`,
			want: &Robots{
				SchemeHost: "https://www.example.com",
				StatusCode: 200,
				Body:       "User-agent: *\nDisallow: /",
				Expires:    "201611071423",
				Cached:     true,
			},
		},
		{
			name:   "not cached",
			host:   "https://api.example.com",
			status: http.StatusNotFound,
			resp: `{
			  "_index": "robots",
			  "_type": "robots",
			  "_id": "https://api.example.com",
			  "found": false
			}`,
			want: &Robots{
				SchemeHost: "https://api.example.com",
			},
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

			got, err := e.Get(c.host)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

func TestPut(t *testing.T) {
	for _, c := range []struct {
		name   string
		robot  *Robots
		status int
		resp   string
	}{
		{
			name: "success",
			robot: &Robots{
				SchemeHost: "https://www.example.com",
				StatusCode: 200,
				Body:       "User-agent: *\nDisallow: /",
				Expires:    "201611071423",
			},
			status: http.StatusCreated,
			resp: `{
			  "took": 27,
			  "errors": false,
			  "items": [
					{
			      "create": {
			        "_index": "robots",
			        "_type": "robots",
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

			e.Put(c.robot)

			if err := e.Bulk.Flush(); err != nil {
				t.Fatal(err)
			}

			stats := e.Bulk.Stats()
			if stats.Succeeded != 1 {
				t.Fatalf("put failed: got %d", stats.Succeeded)
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
			name:   "ok",
			status: http.StatusOK,
			want:   true,
		},
		{
			name:   "not found",
			status: http.StatusNotFound,
			want:   false,
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
			name:   "ok",
			status: http.StatusOK,
			resp:   `{"acknowledged": true}`,
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

	bulk, err := client.BulkProcessor().Stats(true).Do(context.TODO())
	if err != nil {
		return nil, err
	}

	return &ElasticSearch{
		Client: client,
		Index:  "robots",
		Type:   "robots",
		Bulk:   bulk,
	}, nil
}
