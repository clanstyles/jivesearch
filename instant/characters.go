package instant

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jivesearch/jivesearch/instant/contributors"
)

// Characters is an instant answer
type Characters struct {
	Answer
}

func (c *Characters) setQuery(r *http.Request) answerer {
	c.Answer.setQuery(r)
	return c
}

func (c *Characters) setUserAgent(r *http.Request) answerer {
	return c
}

func (c *Characters) setType() answerer {
	c.Type = "characters"
	return c
}

func (c *Characters) setContributors() answerer {
	c.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return c
}

func (c *Characters) setTriggers() answerer {
	c.triggers = []string{
		"number of characters in", "number of characters",
		"number of chars in", "number of chars",
		"char count of", "char count",
		"chars count of", "chars count",
		"character count of", "character count",
		"characters count of", "characters count",
		"length in chars", "length in characters",
	}
	return c
}

func (c *Characters) setTriggerFuncs() answerer {
	c.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return c
}

func (c *Characters) setSolution() answerer {
	for _, ch := range []string{`"`, `'`} {
		c.remainder = strings.TrimPrefix(c.remainder, ch)
		c.remainder = strings.TrimSuffix(c.remainder, ch)
	}

	c.Text = strconv.Itoa(len(c.remainder))

	return c
}

func (c *Characters) setCache() answerer {
	c.Cache = true
	return c
}

func (c *Characters) tests() []test {
	typ := "characters"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query: `number of chars in "Jimi Hendrix"`,
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "12",
					Cache:        true,
				},
			},
		},
		test{
			query: "number of chars   in Pink   Floyd",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "10",
					Cache:        true,
				},
			},
		},
		test{
			query: "Bob Dylan   number of characters in",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "9",
					Cache:        true,
				},
			},
		},
		test{
			query: "number of characters Janis   Joplin",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "12",
					Cache:        true,
				},
			},
		},
		test{
			query: "char count Led Zeppelin",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "12",
					Cache:        true,
				},
			},
		},
		test{
			query: "char count of ' 87 '",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "4",
					Cache:        true,
				},
			},
		},
		test{
			query: "they're chars count",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "7",
					Cache:        true,
				},
			},
		},
		test{
			query: "chars count of something",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "9",
					Cache:        true,
				},
			},
		},
		test{
			query: "Another something chars count of",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "17",
					Cache:        true,
				},
			},
		},
		test{
			query: "1234567 character count",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "7",
					Cache:        true,
				},
			},
		},
		test{
			query: "character count of house of cards",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "14",
					Cache:        true,
				},
			},
		},
		test{
			query: "characters count 50 cent",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "7",
					Cache:        true,
				},
			},
		},
		test{
			query: "characters count of 1 dollar",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "8",
					Cache:        true,
				},
			},
		},
		test{
			query: "chars in saved by the bell",
			expected: []Solution{
				Solution{},
			},
		},
		test{
			query: "chars 21 jump street",
			expected: []Solution{
				Solution{},
			},
		},
		test{
			query: "characters in house of cards",
			expected: []Solution{
				Solution{},
			},
		},
		test{
			query: "characters beavis and butthead",
			expected: []Solution{
				Solution{},
			},
		},
		test{
			query: "char count equity",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "6",
					Cache:        true,
				},
			},
		},
		test{
			query: "characters count seal",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "4",
					Cache:        true,
				},
			},
		},
		test{
			query: "length in chars lion",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "4",
					Cache:        true,
				},
			},
		},
		test{
			query: "length in characters mountain",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "8",
					Cache:        true,
				},
			},
		},
		test{
			query: "length of 1 meter",
			expected: []Solution{
				Solution{},
			},
		},
	}

	return tests
}
