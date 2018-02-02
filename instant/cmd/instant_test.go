package main

import (
	"encoding/json"
	"jivesearch/instant"
	"jivesearch/instant/contributors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	for _, c := range []struct {
		query     string
		userAgent string
		want      *instant.Solution
	}{
		{
			query: "january birthstone",
			want: &instant.Solution{
				Type:         "birthstone",
				Triggered:    true,
				Contributors: contributors.Load([]string{"brentadamson"}),
				Text:         "Garnet",
				Cache:        true,
			},
		},
		{
			query:     "user agent",
			userAgent: "firefox",
			want: &instant.Solution{
				Type:         "user agent",
				Triggered:    true,
				Contributors: contributors.Load([]string{"brentadamson"}),
				Text:         "firefox",
				Cache:        false,
			},
		},
		{
			query: "not an instant answer",
			want:  &instant.Solution{},
		},
	} {
		t.Run(c.query, func(t *testing.T) {
			v := url.Values{}
			v.Set("q", c.query)

			r := &http.Request{
				Form:   v,
				Header: make(http.Header),
			}

			r.Header.Set("User-Agent", c.userAgent)

			rr := httptest.NewRecorder()
			http.HandlerFunc(handler).ServeHTTP(rr, r)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: expected %v got %v",
					http.StatusOK, status)
			}

			got := &instant.Solution{}

			if err := json.NewDecoder(rr.Body).Decode(got); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}
