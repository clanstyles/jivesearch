package frontend

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jivesearch/jivesearch/instant"
	"github.com/jivesearch/jivesearch/log"
	"github.com/jivesearch/jivesearch/search"
	"github.com/jivesearch/jivesearch/wikipedia"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

// Context holds a user's request context so we can pass it to our template's form.
// Query, Language, and Region are the RAW query string variables.
type Context struct {
	Q         string          `json:"query"`
	L         string          `json:"-"`
	R         string          `json:"-"`
	N         string          `json:"-"`
	Preferred []language.Tag  `json:"-"`
	Region    language.Region `json:"-"`
	Number    int             `json:"-"`
	Page      int             `json:"-"`
}

// Results is the results from search, instant, wikipedia, etc
type Results struct {
	Alternative string           `json:"alternative"`
	Instant     instant.Solution `json:"instant"`
	Search      *search.Results  `json:"search"`
	Wikipedia   *wikipedia.Item  `json:"wikipedia"`
}

type data struct {
	Context `json:"-"`
	Results
}

// Detect the user's preferred language(s).
// The "l" param takes precendence over the "Accept-Language" header.
func (f *Frontend) detectLanguage(r *http.Request) []language.Tag {
	preferred := []language.Tag{}
	if lang := strings.TrimSpace(r.FormValue("l")); lang != "" {
		if l, err := language.Parse(lang); err == nil {
			preferred = append(preferred, l)
		}
	}

	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		log.Info.Println(err)
		return preferred
	}

	preferred = append(preferred, tags...)
	return preferred
}

// Detect the user's region. "r" param takes precedence over the language's region (if any).
func (f *Frontend) detectRegion(lang language.Tag, r *http.Request) language.Region {
	reg, err := language.ParseRegion(strings.TrimSpace(r.FormValue("r")))
	if err != nil {
		reg, _ = lang.Region()
	}

	return reg.Canonicalize()
}

func (f *Frontend) addQuery(q string) error {
	exists, err := f.Suggest.Exists(q)
	if err != nil {
		return err
	}

	if !exists {
		if err := f.Suggest.Insert(q); err != nil {
			return err
		}
	}

	return f.Suggest.Increment(q)
}

func (f *Frontend) searchHandler(w http.ResponseWriter, r *http.Request) *response {
	d := data{
		Context{
			Q: strings.TrimSpace(r.FormValue("q")),
			L: strings.TrimSpace(r.FormValue("l")),
			N: strings.TrimSpace(r.FormValue("n")),
			R: strings.TrimSpace(r.FormValue("r")),
		},
		Results{
			Search: &search.Results{},
			Wikipedia: &wikipedia.Item{
				Wikidata: &wikipedia.Wikidata{
					Claims: &wikipedia.Claims{},
				},
			},
		},
	}

	resp := &response{
		status:   http.StatusOK,
		template: "search",
		err:      nil,
	}

	if d.Context.Q == "" { // render start page if no query
		return resp
	}

	d.Context.Preferred = f.detectLanguage(r)
	lang, _, _ := f.Document.Matcher.Match(d.Context.Preferred...) // will use first supported tag in case of error

	d.Context.Region = f.detectRegion(lang, r)

	// is it a !bang? Redirect them
	if loc, ok := f.Bangs.Detect(d.Context.Q, d.Context.Region.String(), lang.String()); ok {
		return &response{
			status:   302,
			redirect: loc,
		}
	}

	// Let's get them their results
	// what page are they on? Give them first page by default
	var err error
	d.Context.Page, err = strconv.Atoi(strings.TrimSpace(r.FormValue("p")))
	if err != nil || d.Context.Page < 1 {
		d.Context.Page = 1
	}

	// how many results wanted?
	d.Context.Number, err = strconv.Atoi(strings.TrimSpace(r.FormValue("n")))
	if err != nil || d.Context.Number > 100 {
		d.Context.Number = 25
	}

	channels := 1
	sc := make(chan *search.Results)
	var ac chan error
	var ic chan instant.Solution
	var wc chan *wikipedia.Item

	strt := time.Now() // we already have total response time in nginx...we want the breakdown

	if d.Context.Page == 1 {
		channels += 3

		ac = make(chan error)
		ic = make(chan instant.Solution)
		wc = make(chan *wikipedia.Item)

		go func(q string, ch chan error) {
			ch <- f.addQuery(q)
		}(d.Context.Q, ac)

		go func(r *http.Request) {
			ic <- instant.Detect(r)
		}(r)

		go func(d data) {
			w, err := f.wikiHandler(d.Context.Q, d.Context.Preferred)
			if err != nil && err != sql.ErrNoRows {
				log.Info.Println(err)
			}
			wc <- w
		}(d)
	}

	go func(d data, lang language.Tag, region language.Region) {
		// get the votes
		offset := d.Context.Page*d.Context.Number - d.Context.Number
		votes, err := f.Vote.Get(d.Context.Q, d.Context.Number*10) // get votes for first 10 pages
		if err != nil {
			log.Info.Println(err)
		}

		res, err := f.Search.Fetch(d.Context.Q, lang, region, d.Context.Number, offset, votes)
		if err != nil {
			log.Info.Println(err)
		}

		for _, doc := range res.Documents {
			for _, v := range votes {
				if doc.ID == v.URL {
					doc.Votes = v.Votes
				}
			}
		}

		res = res.AddPagination(d.Context.Number, d.Context.Page) // move this to javascript??? (Wouldn't be available in API....)
		sc <- res
	}(d, lang, d.Context.Region)

	stats := struct {
		autocomplete time.Duration
		instant      time.Duration
		wikipedia    time.Duration
		search       time.Duration
	}{}

	for i := 0; i < channels; i++ {
		select {
		case d.Instant = <-ic:
			if d.Instant.Err != nil {
				log.Info.Println(d.Instant.Err)
			}
			stats.instant = time.Since(strt).Round(time.Microsecond)
		case d.Wikipedia = <-wc:
			stats.wikipedia = time.Since(strt).Round(time.Millisecond)
		case d.Search = <-sc:
			stats.search = time.Since(strt).Round(time.Millisecond)
		case err := <-ac:
			if err != nil {
				log.Info.Println(err)
			}
			stats.autocomplete = time.Since(strt).Round(time.Millisecond)
		case <-r.Context().Done():
			// TODO: add info on which items took too long...
			// Perhaps change status code of response so it isn't cached by nginx
			log.Info.Println(errors.Wrapf(r.Context().Err(), "timeout on retrieving results"))
		}
	}

	log.Info.Printf("ac:%v, instant:%v, search:%v, wiki:%v\n", stats.autocomplete, stats.instant, stats.search, stats.wikipedia)

	if r.FormValue("o") == "json" {
		resp.template = r.FormValue("o")
	}

	resp.data = d
	return resp
}

func (f *Frontend) wikiHandler(query string, preferred []language.Tag) (*wikipedia.Item, error) {
	var err error
	item := &wikipedia.Item{}

	lang, _, _ := f.Wikipedia.Matcher.Match(preferred...)
	item, err = f.Wikipedia.Fetch(query, lang)
	if err != nil {
		return item, err
	}

	return item, err
}
