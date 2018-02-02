package instant

import (
	"jivesearch/instant/contributors"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
)

// Random is an instant answer
type Random struct {
	Answer
}

var reRandom *regexp.Regexp

func (r *Random) setQuery(req *http.Request) answerer {
	r.Answer.setQuery(req)
	return r
}

func (r *Random) setUserAgent(req *http.Request) answerer {
	return r
}

func (r *Random) setType() answerer {
	r.Type = "random"
	return r
}

func (r *Random) setContributors() answerer {
	r.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return r
}

func (r *Random) setTriggers() answerer {
	r.triggers = []string{
		"random number", "random number between",
	}
	return r
}

func (r *Random) setTriggerFuncs() answerer {
	r.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return r
}

func (r *Random) setSolution() answerer {
	matches := make(map[string]int)
	matches["min"], matches["max"] = 1, 100 // if no range specified

	match := reRandom.FindStringSubmatch(r.remainder)

	if len(match) > 0 {
		for i, name := range reRandom.SubexpNames() {
			if i == 0 {
				continue
			}
			if integer, err := strconv.Atoi(match[i]); err == nil {
				matches[name] = integer
			}
		}
		if matches["max"] < matches["min"] {
			matches["min"], matches["max"] = matches["max"], matches["min"]
		}
	}

	r.Text = strconv.Itoa(rand.Intn(matches["max"]+1-matches["min"]) + matches["min"])

	return r
}

func (r *Random) setCache() answerer {
	r.Cache = false
	return r
}

func (r *Random) tests() []test {
	typ := "random"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{}

	solutions := func(choices []string) []Solution {
		sol := []Solution{}

		for _, c := range choices {
			sol = append(sol,
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         c,
					Cache:        false,
				},
			)
		}

		return sol
	}

	for _, c := range []struct {
		q   string
		sol []string
	}{
		{"Random number between 1 and 3", []string{"1", "2", "3"}},
		{"Random number between 5431 and 5434", []string{"5431", "5432", "5433", "5434"}},
		{"Random number between -18 and -21", []string{"-18", "-19", "-20", "-21"}},
	} {
		t := test{
			query:    c.q,
			expected: solutions(c.sol),
		}
		tests = append(tests, t)
	}

	return tests
}

func init() {
	reRandom = regexp.MustCompile(`(?P<min>-?\d+).*?(?P<max>-?\d+)`)
}
