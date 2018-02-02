// Package frontend provides the routing and middleware for the web app
package frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"jivesearch/bangs"
	"jivesearch/log"
	"jivesearch/search"
	"jivesearch/search/vote"
	"jivesearch/suggest"
	"jivesearch/wikipedia"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/oxtoacart/bpool"
	"github.com/pkg/errors"
)

// Frontend holds settings for our languages supported, backend, etc.
// better name???
type Frontend struct {
	Document
	*bangs.Bangs
	Suggest suggest.Suggester
	Search  search.Fetcher
	Wikipedia
	Vote vote.Voter
}

// Document has the languages we support
type Document struct {
	Languages []language.Tag
	language.Matcher
}

// Wikipedia holds our settings for wikipedia/wikidata
// Note: language matcher here may be different than that for
// document due to available languages Wikipedia supports
type Wikipedia struct {
	language.Matcher
	wikipedia.Fetcher
}

var (
	bufpool   *bpool.BufferPool // makes sure no errors when writing to our templates
	templates map[string]*template.Template
)

func init() {
	bufpool = bpool.NewBufferPool(48) // what is the appropriate size??? 48??? 64???
}

type response struct {
	status   int
	redirect string
	template string
	data     interface{}
	err      error
}

type appHandler func(http.ResponseWriter, *http.Request) *response

// middleware sets a timeout and then serves.
func (f *Frontend) middleware(next appHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rsp := fn(w, r); rsp != nil {
		switch rsp.status {
		case http.StatusOK:
			buf := bufpool.Get()
			defer bufpool.Put(buf)

			switch rsp.template {
			case "json":
				w.Header().Set("Content-Type", "application/json") // the default for json is utf-8
				if err := json.NewEncoder(buf).Encode(rsp.data); err != nil {
					rsp.status, rsp.err = http.StatusInternalServerError, err
					errHandler(w, rsp)
					return
				}
			default: // html by default
				w.Header().Set("Content-Type", "text/html; charset=utf-8")

				tmpl, ok := templates[rsp.template]
				if !ok {
					rsp.status = http.StatusInternalServerError
					rsp.err = fmt.Errorf("template doesn't exist: %q", rsp.template)
					errHandler(w, rsp)
					return
				}

				if err := tmpl.Execute(buf, rsp.data); err != nil {
					rsp.status, rsp.err = http.StatusInternalServerError, err
					errHandler(w, rsp)
					return
				}
			}

			if _, err := buf.WriteTo(w); err != nil {
				rsp.status, rsp.err = http.StatusInternalServerError, err
				errHandler(w, rsp)
			}
		case http.StatusFound: // !bang
			http.Redirect(w, r, rsp.redirect, http.StatusFound)
		case http.StatusBadRequest, http.StatusInternalServerError:
			errHandler(w, rsp)
		default:
			log.Info.Printf("Unknown status %d\n", rsp.status)
		}
	}
}

func errHandler(w http.ResponseWriter, rsp *response) {
	switch rsp.status {
	case http.StatusBadRequest:
		log.Debug.Println(rsp.err)
	case http.StatusInternalServerError:
		log.Info.Println(rsp.err)
	}

	http.Error(w, http.StatusText(rsp.status), rsp.status)
}

func (f *Frontend) autocompleteHandler(w http.ResponseWriter, r *http.Request) *response {
	q := strings.TrimSpace(r.FormValue("q"))
	res, err := f.Suggest.Completion(q, 10)
	if err != nil {
		return &response{
			status: http.StatusInternalServerError,
			err:    errors.Wrapf(err, "autocomplete error %q (%v)", q, res.RawQuery),
		}
	}
	return &response{
		status:   http.StatusOK,
		template: "json",
		data:     res,
	}
}

// ParseTemplates parses our html templates.
var ParseTemplates = func() {
	templates = make(map[string]*template.Template)
	templates["search"] = template.Must(
		template.New("base.html").
			Funcs(funcMap).
			ParseFiles(
				"templates/base.html",
				"templates/main.css",
				"templates/search.html",
			),
	)
}
