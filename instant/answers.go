// Package instant provides instant answers
package instant

import (
	"jivesearch/instant/contributors"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// QueryVar is the http request variable to parse
var QueryVar = "q"

// answerer outlines methods for an instant answer
type answerer interface {
	setQuery(r *http.Request) answerer
	setUserAgent(r *http.Request) answerer
	setType() answerer
	setContributors() answerer
	setTriggers() answerer
	setTriggerFuncs() answerer
	trigger() bool
	setSolution() answerer
	setCache() answerer
	solution() Solution
	tests() []test
}

// Answer holds an instant answer when triggered
type Answer struct {
	query        string
	userAgent    string
	triggers     []string
	triggerFuncs []triggerFunc
	remainder    string
	Solution
}

// Solution holds the Text, Data and HTML of an answer
type Solution struct {
	Type         string                     `json:"type,omitempty"`
	Triggered    bool                       `json:"triggered"`
	Contributors []contributors.Contributor `json:"contributors,omitempty"`
	Text         string                     `json:"text,omitempty"`
	HTML         string                     `json:"html,omitempty"` // TODO: custom html
	Err          error                      `json:"error,omitempty"`
	Cache        bool                       `json:"cache,omitempty"`
}

// Detect loops through all instant answers to find a solution
var Detect = func(r *http.Request) Solution {
	for _, ia := range answers() {
		ia.setUserAgent(r)
		ia.setQuery(r).setTriggers().setTriggerFuncs()
		if triggered := ia.trigger(); triggered {
			ia.setType().
				setContributors().
				setCache().
				setSolution()
			return ia.solution()
		}
	}

	return Solution{}
}

// setQuery sets the query field
// If future answers need custom setQuery methods we
// could implement same model as we do for setTriggerFuncs()
func (a *Answer) setQuery(r *http.Request) {
	q := strings.ToLower(strings.TrimSpace(r.FormValue(QueryVar)))
	q = strings.Trim(q, "?")
	a.query = strings.Join(strings.Fields(q), " ") // Replace multiple whitespace w/ single whitespace
}

// trigger executes the triggerer for an instant answer
func (a *Answer) trigger() bool {
	for _, t := range a.triggerFuncs {
		if a := t(a); a.Triggered {
			return a.Triggered
		}
	}
	return a.Triggered
}

type triggerFunc func(a *Answer) *Answer

var startsWith triggerFunc = func(a *Answer) *Answer {
	for _, w := range a.triggers {
		if pre := strings.TrimPrefix(a.query, w); pre != a.query {
			a.remainder = strings.TrimSpace(pre)
			a.Triggered = true
			return a
		}
	}
	return a
}

var endsWith triggerFunc = func(a *Answer) *Answer {
	for _, w := range a.triggers {
		if suff := strings.TrimSuffix(a.query, w); suff != a.query {
			a.remainder = strings.TrimSpace(suff)
			a.Triggered = true
			return a
		}
	}
	return a
}

func (a *Answer) solution() Solution {
	return a.Solution
}

type test struct {
	query     string
	userAgent string
	expected  []Solution
}

// answers returns a slice of all instant answers
// Note: Since we modify fields of the answers we probably shouldn't reuse them....
func answers() []answerer {
	return []answerer{
		&BirthStone{},
		&CamelCase{},
		&Characters{},
		&Coin{},
		&Frequency{},
		&Potus{},
		&Prime{},
		&Random{},
		&Reverse{},
		&Stats{},
		&Temperature{},
		&UserAgent{},
	}
}

func init() {
	// for coin, random & probably others down the road
	rand.Seed(time.Now().UTC().UnixNano())
}
