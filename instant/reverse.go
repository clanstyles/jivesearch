package instant

import (
	"jivesearch/instant/contributors"
	"net/http"
	"strings"
)

// Reverse is an instant answer
type Reverse struct {
	Answer
}

func (r *Reverse) setQuery(req *http.Request) answerer {
	r.Answer.setQuery(req)
	return r
}

func (r *Reverse) setUserAgent(req *http.Request) answerer {
	return r
}

func (r *Reverse) setType() answerer {
	r.Type = "reverse"
	return r
}

func (r *Reverse) setContributors() answerer {
	r.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return r
}

func (r *Reverse) setTriggers() answerer {
	r.triggers = []string{
		"reverse",
	}
	return r
}

func (r *Reverse) setTriggerFuncs() answerer {
	r.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return r
}

func (r *Reverse) setSolution() answerer {
	for _, c := range []string{`"`, `'`} {
		r.remainder = strings.TrimPrefix(r.remainder, c)
		r.remainder = strings.TrimSuffix(r.remainder, c)
	}

	var n int
	rune := make([]rune, len(r.remainder))
	for _, j := range r.remainder {
		rune[n] = j
		n++
	}
	rune = rune[0:n]

	// Reverse
	for i, j := 0, len(rune)-1; i < j; i, j = i+1, j-1 {
		rune[i], rune[j] = rune[j], rune[i]
	}

	r.Text = string(rune)

	return r
}

func (r *Reverse) setCache() answerer {
	r.Cache = true
	return r
}

func (r *Reverse) tests() []test {
	typ := "reverse"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query: "reverse ahh lights....ahh see 'em",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "me' ees hha....sthgil hha",
					Cache:        true,
				},
			},
		},
		test{
			query: "reverse 私日本語は話せません",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "んせませ話は語本日私",
					Cache:        true,
				},
			},
		},
		test{
			query: `reverse "ahh yeah"`,
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "haey hha",
					Cache:        true,
				},
			},
		},
	}

	return tests
}
