package instant

import (
	"math/rand"
	"net/http"

	"github.com/jivesearch/jivesearch/instant/contributors"
)

// Coin is an instant answer
type Coin struct {
	Answer
}

func (c *Coin) setQuery(r *http.Request) answerer {
	c.Answer.setQuery(r)
	return c
}

func (c *Coin) setUserAgent(r *http.Request) answerer {
	return c
}

func (c *Coin) setType() answerer {
	c.Type = "coin toss"
	return c
}

func (c *Coin) setContributors() answerer {
	c.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return c
}

func (c *Coin) setTriggers() answerer {
	c.triggers = []string{
		"flip a coin", "heads or tails", "coin toss",
	}
	return c
}

func (c *Coin) setTriggerFuncs() answerer {
	c.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}

	return c
}

func (c *Coin) setSolution() answerer {
	choices := []string{"Heads", "Tails"}

	c.Text = choices[rand.Intn(2)]

	return c
}

func (c *Coin) setCache() answerer {
	c.Cache = false
	return c
}

func (c *Coin) tests() []test {
	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{}

	for _, q := range []string{"flip a coin", "heads or tails", "Coin Toss"} {
		tst := test{
			query: q,
			expected: []Solution{
				Solution{
					Type:         "coin toss",
					Triggered:    true,
					Contributors: contrib,
					Text:         "Heads",
					Cache:        false,
				},
				Solution{
					Type:         "coin toss",
					Triggered:    true,
					Contributors: contrib,
					Text:         "Tails",
					Cache:        false,
				},
			},
		}

		tests = append(tests, tst)
	}

	return tests
}
