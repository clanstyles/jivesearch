// Sample bangs demonstrates how to run a simple !bangs server
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jivesearch/jivesearch/bangs"
	"golang.org/x/text/language"
)

type config struct {
	*bangs.Bangs
}

func (c *config) handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", c.bangsHandler)
	r.HandleFunc("/favicon.ico", favHandler)
	return r
}

func (c *config) bangsHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(strings.TrimSpace(r.FormValue("q")))

	// Language and region parsing are outside the
	// scope of this package. See github.com/jivesearch/frontend/search.go for an example.
	l := language.MustParse("en-US")
	reg, _ := l.Region()

	if loc, ok := c.Bangs.Detect(q, reg.String(), l.String()); ok {
		http.Redirect(w, r, loc, http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("No !bang found for %q", q)))
}

func favHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	c := &config{
		bangs.New(),
	}

	port := 8000

	log.Printf("Listening at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), c.handler()))
}
