package instant

import (
	"net/http"
	"strings"

	"github.com/jivesearch/jivesearch/instant/contributors"
)

// CamelCase is an instant answer
type CamelCase struct {
	Answer
}

func (c *CamelCase) setQuery(r *http.Request) answerer {
	c.Answer.setQuery(r)
	return c
}

func (c *CamelCase) setUserAgent(r *http.Request) answerer {
	return c
}

func (c *CamelCase) setType() answerer {
	c.Type = "camelcase"
	return c
}

func (c *CamelCase) setContributors() answerer {
	c.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return c
}

func (c *CamelCase) setTriggers() answerer {
	c.triggers = []string{
		"camelcase",
		"camel case",
	}
	return c
}

func (c *CamelCase) setTriggerFuncs() answerer {
	c.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return c
}

func (c *CamelCase) setSolution() answerer {
	titled := []string{}
	for _, w := range strings.Fields(c.remainder) {
		titled = append(titled, strings.Title(w))
	}

	c.Text = strings.Join(titled, "")

	return c
}

func (c *CamelCase) setCache() answerer {
	c.Cache = true
	return c
}

func (c *CamelCase) tests() []test {
	typ := "camelcase"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query: "camelcase metallica rocks",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "MetallicaRocks",
					Cache:        true,
				},
			},
		},
		test{
			query: "aliCE in chAins Is better camel case",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "AliceInChainsIsBetter",
					Cache:        true,
				},
			},
		},
		test{
			query: "camel case O'doyle ruLES",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "O'DoyleRules",
					Cache:        true,
				},
			},
		},
	}

	return tests
}
