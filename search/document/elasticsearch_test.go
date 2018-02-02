package document

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olivere/elastic"
	"golang.org/x/text/language"
)

func TestAnalyzers(t *testing.T) {
	// make sure we have an analyzer mapped for each language listed in document.go
	for _, lang := range available {
		if _, ok := langAnalyzer[lang]; !ok {
			t.Fatalf("no elasticsearch analyzer found for %q", lang.String())
		}
	}
}

func TestAnalyzer(t *testing.T) {
	for _, c := range []struct {
		name string
		lang language.Tag
		want string
	}{
		{"English", language.English, "english"},
		{"British English", language.BritishEnglish, "english"},
		{"Spanish", language.Spanish, "spanish"},
		{"European Spanish", language.EuropeanSpanish, "spanish"},
		{"Latin American Spanish", language.LatinAmericanSpanish, "spanish"},
		{"German", language.German, "german"},
		{"Portuguese", language.Portuguese, "portuguese"},
		{"European Portuguese", language.EuropeanPortuguese, "portuguese"},
		{"Brazilian Portuguese", language.BrazilianPortuguese, "brazilian"},
	} {
		t.Run(c.name, func(t *testing.T) {
			handler := http.NotFound
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler(w, r)
			}))
			defer ts.Close()

			handler = func(w http.ResponseWriter, r *http.Request) {}

			e, err := MockService(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			got, err := e.Analyzer(c.lang)
			if err != nil {
				t.Fatal(err)
			}

			if got != c.want {
				t.Fatalf("got %q; want %q", got, c.want)
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

	return &ElasticSearch{
		Client: client, Index: "search", Type: "document",
	}, nil
}
