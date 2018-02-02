package instant

import (
	"jivesearch/instant/contributors"
	"net/http"
)

// UserAgent is an instant answer
type UserAgent struct {
	Answer
}

func (u *UserAgent) setQuery(r *http.Request) answerer {
	u.Answer.setQuery(r)
	return u
}

func (u *UserAgent) setUserAgent(r *http.Request) answerer {
	u.Answer.userAgent = r.UserAgent()
	return u
}

func (u *UserAgent) setType() answerer {
	u.Type = "user agent"
	return u
}

func (u *UserAgent) setContributors() answerer {
	u.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)

	return u
}

func (u *UserAgent) setTriggers() answerer {
	u.triggers = []string{
		"user agent", "user agent?",
		"useragent", "useragent?",
		"my user agent", "my user agent?",
		"my useragent", "my useragent?",
		"what's my user agent", "what's my user agent?",
		"what's my useragent", "what's my useragent?",
		"what is my user agent", "what is my user agent?",
		"what is my useragent", "what is my useragent?",
	}

	return u
}

func (u *UserAgent) setTriggerFuncs() answerer {
	u.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}

	return u
}

func (u *UserAgent) setSolution() answerer {
	u.Text = u.userAgent

	return u
}

func (u *UserAgent) setCache() answerer {
	// caching would cache the query but the browser could change
	u.Cache = false
	return u
}

func (u *UserAgent) tests() []test {
	typ := "user agent"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query:     "user agent",
			userAgent: "firefox",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "firefox",
					Cache:        false,
				},
			},
		},
		test{
			query:     "useragent?",
			userAgent: "opera",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "opera",
					Cache:        false,
				},
			},
		},
		test{
			query:     "my user agent",
			userAgent: "some random ua",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "some random ua",
					Cache:        false,
				},
			},
		},
		test{
			query:     "what's my user agent?",
			userAgent: "chrome",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "chrome",
					Cache:        false,
				},
			},
		},
		test{
			query:     "what is my useragent?",
			userAgent: "internet explorer",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "internet explorer",
					Cache:        false,
				},
			},
		},
	}

	return tests
}
