package frontend

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jivesearch/jivesearch/suggest"
	"github.com/spf13/pflag"
)

func TestMiddleware(t *testing.T) {
	type want struct {
		status int
		body   string
	}

	for _, c := range []struct {
		name  string
		tmpl  string
		ct    string
		cl    string
		sniff string
		want
	}{
		{"json", "json", "application/json", "28", "",
			want{http.StatusOK, "{\"response\":\"hello world!\"}\n"},
		},
		{"wrong template", "", "text/plain; charset=utf-8", "22", "nosniff",
			want{http.StatusInternalServerError, "Internal Server Error\n"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			want := make(http.Header)
			want["Content-Length"] = []string{c.cl}
			want["Content-Type"] = []string{c.ct}

			if c.sniff != "" {
				want["X-Content-Type-Options"] = []string{c.sniff}
			}

			f := &Frontend{}
			ParseTemplates()

			fn := func(w http.ResponseWriter, r *http.Request) *response {
				return &response{
					status:   200,
					template: c.tmpl,
					data:     map[string]string{"response": "hello world!"},
					err:      nil,
				}
			}

			ts := httptest.NewServer(f.middleware(appHandler(fn)))
			defer ts.Close()

			resp, err := http.Get(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != c.want.status {
				t.Fatalf("got %d; want %d", resp.StatusCode, c.want.status)
			}

			h := resp.Header
			delete(h, "Date") // is there a way to mock this instead???

			if !reflect.DeepEqual(h, want) {
				t.Fatalf("got %v; want %v", h, want)
			}

			bdy, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			got := string(bdy)
			if got != c.want.body {
				t.Fatalf("got %q; want %q", got, c.want.body)
			}
		})
	}
}

func TestAutocompleteHandler(t *testing.T) {
	for _, c := range []struct {
		name string
		q    string
		want *response
	}{
		{"basic", "r",
			&response{
				status:   http.StatusOK,
				template: "json",
				data: suggest.Results{
					Suggestions: []string{
						"radiohead",
						"rage against the machine",
						"red hot chili peppers",
						"r.e.m.",
						"rolling stones",
						"rollins band",
						"rusted root",
					},
				},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &Frontend{
				Suggest: &mockSuggester{},
			}

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			q := req.URL.Query()
			q.Add("q", c.q)
			req.URL.RawQuery = q.Encode()

			got := f.autocompleteHandler(httptest.NewRecorder(), req)

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

type mockProvider struct {
	m map[string]interface{}
}

func (p *mockProvider) SetDefault(key string, value interface{}) {
	p.m[key] = value
}
func (p *mockProvider) SetTypeByDefaultValue(bool) {}
func (p *mockProvider) BindPFlag(key string, flg *pflag.Flag) error {
	return nil
}
func (p *mockProvider) Get(key string) interface{} {
	return p.m[key]
}
func (p *mockProvider) GetString(key string) string {
	return p.m[key].(string)
}
func (p *mockProvider) GetInt(key string) int {
	return p.m[key].(int)
}
func (p *mockProvider) GetStringSlice(key string) []string {
	return p.m[key].([]string)
}

type mockSuggester struct {
	ex bool
}

func (ms *mockSuggester) Exists(q string) (bool, error) {
	return ms.ex, nil
}

func (ms *mockSuggester) Insert(q string) error {
	return nil
}

func (ms *mockSuggester) Increment(q string) error {
	return nil
}

func (ms *mockSuggester) Completion(q string, size int) (suggest.Results, error) {
	s := suggest.Results{}

	if q == "r" {
		s.Suggestions = []string{
			"radiohead",
			"rage against the machine",
			"red hot chili peppers",
			"r.e.m.",
			"rolling stones",
			"rollins band",
			"rusted root",
		}
	}
	return s, nil
}

func (ms *mockSuggester) IndexExists() (bool, error) {
	return ms.ex, nil
}

func (ms *mockSuggester) Setup() error { return nil }

func TestParseTemplates(t *testing.T) {
	ParseTemplates()

	if _, ok := templates["search"]; !ok {
		t.Fatal("Our search template is not in our templates map.")
	}
}
