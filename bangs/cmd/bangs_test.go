package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"jivesearch/bangs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBangsHandler(t *testing.T) {
	for _, c := range []struct {
		name  string
		query string
		want  int
	}{
		{"empty", "", 200},
		{"no !bang", "a pants", 200},
		{"standard !bang", "!a pants", 302},
		{"backwards bang!", "Brad Pitt g!", 302},
	} {
		t.Run(c.name, func(t *testing.T) {
			cfg := config{
				bangs.New(),
			}

			v := url.Values{}
			v.Set("q", c.query)

			r := &http.Request{
				Form:   v,
				Header: make(http.Header),
			}

			rec := httptest.NewRecorder()
			cfg.bangsHandler(rec, r)

			if rec.Code != c.want {
				t.Fatalf("got %v; want %v", rec.Code, c.want)
			}
		})
	}
}

func TestRouting(t *testing.T) {
	cfg := config{
		bangs.New(),
	}

	srv := httptest.NewServer(cfg.handler())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%v/?q=Brad+Pitt", srv.URL))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK; got %v", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	want := `No !bang found for "brad pitt"`
	got := string(bytes.TrimSpace(b))
	if got != want {
		t.Fatalf("got %q; want %q", got, want)
	}
}
