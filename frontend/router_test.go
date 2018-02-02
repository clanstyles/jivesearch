package frontend

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestRouter(t *testing.T) {
	type route struct {
		name   string
		method string
		url    string
	}

	for _, c := range []*route{
		&route{
			name:   "search",
			method: "GET",
			url:    "https://www.example.com/?q=search+term",
		},
		&route{
			name:   "autocomplete",
			method: "GET",
			url:    "http://127.0.0.1/autocomplete",
		},
		&route{
			name:   "vote",
			method: "POST",
			url:    "http://localhost/vote",
		},
		&route{
			name:   "favicon",
			method: "GET",
			url:    "http://example.com/favicon.ico",
		},
		&route{
			name:   "static",
			method: "GET",
			url:    "https://example.com/static/main.js",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cfg := &mockProvider{
				m: make(map[string]interface{}),
			}
			cfg.SetDefault("hmac.secret", "very secret")

			f := &Frontend{}
			router := f.Router(cfg)

			expected, err := http.NewRequest(
				c.method,
				c.url,
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			route := router.Get(c.name)

			if !route.Match(expected, &mux.RouteMatch{}) {
				t.Fatalf("expected route for %q to exist. It doesn't", c.url)
			}
		})
	}
}
