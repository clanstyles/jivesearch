package frontend

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jivesearch/jivesearch/bangs"
	"github.com/jivesearch/jivesearch/instant"
	"github.com/jivesearch/jivesearch/search"
	"github.com/jivesearch/jivesearch/search/document"
	"github.com/jivesearch/jivesearch/search/vote"
	"github.com/jivesearch/jivesearch/wikipedia"
	"golang.org/x/text/language"
)

func TestDetectLanguage(t *testing.T) {
	for _, c := range []struct {
		name           string
		acceptLanguage string
		l              string
		want           []language.Tag
	}{
		{
			"blank", "", "", []language.Tag{},
		},
		{
			"basic", "", "en", []language.Tag{language.English},
		},
		{
			"french", "", "fr", []language.Tag{language.French},
		},
		{
			"Accept-Language header",
			"fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7",
			"",
			[]language.Tag{
				language.MustParse("fr-CH"),
				language.French,
				language.English,
				language.German,
			},
		},
		{
			"param overrides Accept-Language header",
			"fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7",
			"hr",
			[]language.Tag{
				language.Croatian,
				language.MustParse("fr-CH"),
				language.French,
				language.English,
				language.German,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &Frontend{}

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Accept-Language", c.acceptLanguage)

			q := req.URL.Query()
			q.Add("l", c.l)

			req.URL.RawQuery = q.Encode()

			got := f.detectLanguage(req)

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

func TestDetectRegion(t *testing.T) {
	for _, c := range []struct {
		name string
		lang language.Tag
		r    string
		want language.Region
	}{
		{
			"empty", language.Tag{}, "", language.MustParseRegion("US").Canonicalize(),
		},
		{
			"basic", language.Tag{}, "us", language.MustParseRegion("US").Canonicalize(),
		},
		{
			"region from language", language.BrazilianPortuguese, "", language.MustParseRegion("BR").Canonicalize(),
		},
		{
			"param overrides language's region", language.CanadianFrench, "gb", language.MustParseRegion("GB").Canonicalize(),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &Frontend{}

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			q := req.URL.Query()
			q.Add("r", c.r)

			req.URL.RawQuery = q.Encode()

			got := f.detectRegion(c.lang, req)

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

func TestSearchHandler(t *testing.T) {
	for _, c := range []struct {
		name     string
		language string
		query    string
		output   string
		want     *response
	}{
		{
			"empty", "en", "", "",
			&response{
				status:   http.StatusOK,
				template: "search",
				data:     nil,
			},
		},
		{
			"basic", "en", " some query ", "",
			&response{
				status:   http.StatusOK,
				template: "search",
				data: data{
					Context: Context{
						Q:         "some query",
						L:         "en",
						Preferred: []language.Tag{language.MustParse("en")},
						Region:    language.MustParseRegion("US"),
						Number:    25,
						Page:      1,
					},
					Results: Results{
						Search: &search.Results{
							Count:      int64(25),
							Page:       "1",
							Previous:   "",
							Next:       "2",
							Last:       "72",
							Pagination: []string{"1"},
							Documents:  []*document.Document{},
						},
						Wikipedia: &wikipedia.Item{},
					},
				},
			},
		},
		{
			"json", "en", " some query", "json",
			&response{
				status:   http.StatusOK,
				template: "json",
				data: data{
					Context: Context{
						Q:         "some query",
						L:         "en",
						Preferred: []language.Tag{language.MustParse("en")},
						Region:    language.MustParseRegion("US"),
						Number:    25,
						Page:      1,
					},
					Results: Results{
						Search: &search.Results{
							Count:      int64(25),
							Page:       "1",
							Previous:   "",
							Next:       "2",
							Last:       "72",
							Pagination: []string{"1"},
							Documents:  []*document.Document{},
						},
						Wikipedia: &wikipedia.Item{},
					},
				},
			},
		},
		{
			"!bang", "", "!g something", "",
			&response{
				status:   http.StatusFound,
				redirect: "https://encrypted.google.com/search?hl=en&q=something",
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var matcher = language.NewMatcher(
				[]language.Tag{
					language.English,
					language.French,
				},
			)

			f := &Frontend{
				Document: Document{
					Matcher: matcher,
				},
				Bangs:   bangs.New(),
				Suggest: &mockSuggester{},
				Search:  &mockSearch{},
				Wikipedia: Wikipedia{
					Matcher: matcher,
					Fetcher: &mockWikipedia{},
				},
				Vote: &mockVoter{},
			}

			// override instant answer detection for mocking
			instant.Detect = func(r *http.Request) instant.Solution {
				return instant.Solution{}
			}

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			q := req.URL.Query()
			q.Add("q", c.query)
			q.Add("l", c.language)
			q.Add("o", c.output)
			req.URL.RawQuery = q.Encode()

			got := f.searchHandler(httptest.NewRecorder(), req)

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

type mockSearch struct{}

func (s *mockSearch) Fetch(q string, lang language.Tag, region language.Region, page int, number int, votes []vote.Result) (*search.Results, error) {
	r := &search.Results{
		Count:      int64(25),
		Page:       "1",
		Next:       "2",
		Last:       "72",
		Pagination: []string{"2", "3", "4", "5"},
		Documents:  []*document.Document{},
	}

	return r, nil
}

type mockWikipedia struct{}

func (w *mockWikipedia) Fetch(query string, lang language.Tag) (*wikipedia.Item, error) {
	return &wikipedia.Item{}, nil
}

func (w *mockWikipedia) Setup() error {
	return nil
}
