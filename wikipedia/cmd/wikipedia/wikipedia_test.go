package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jivesearch/jivesearch/wikipedia"
	"golang.org/x/text/language"
)

func TestHandler(t *testing.T) {
	for _, c := range []struct {
		name  string
		query string
		want  int
	}{
		{"james", "james", 200},
	} {
		t.Run(c.name, func(t *testing.T) {
			mf := &mockFetcher{}
			f := &fetcher{
				Fetcher: mf,
			}

			v := url.Values{}
			v.Set("q", c.query)

			r := &http.Request{
				URL:    &url.URL{},
				Form:   v,
				Header: make(http.Header),
			}

			rec := httptest.NewRecorder()
			f.wikiHandler(rec, r)

			if rec.Code != c.want {
				t.Fatalf("got %v; want %v", rec.Code, c.want)
			}
		})
	}
}

func TestSetup(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"basic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := setup()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type mockFetcher struct{}

func (mf *mockFetcher) Fetch(query string, lang language.Tag) (*wikipedia.Item, error) {
	return &wikipedia.Item{}, nil
}

func (mf *mockFetcher) Setup() error {
	return nil
}
