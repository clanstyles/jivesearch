package instant

import (
	"net/http"

	"github.com/jivesearch/jivesearch/instant/contributors"
)

// BirthStone is an instant answer
type BirthStone struct {
	Answer
}

func (b *BirthStone) setQuery(r *http.Request) answerer {
	b.Answer.setQuery(r)
	return b
}

func (b *BirthStone) setUserAgent(r *http.Request) answerer {
	return b
}

func (b *BirthStone) setType() answerer {
	b.Type = "birthstone"
	return b
}

func (b *BirthStone) setContributors() answerer {
	b.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return b
}

func (b *BirthStone) setTriggers() answerer {
	b.triggers = []string{
		"birthstones",
		"birth stones",
		"birthstone",
		"birth stone",
	}
	return b
}

func (b *BirthStone) setTriggerFuncs() answerer {
	b.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return b
}

func (b *BirthStone) setSolution() answerer {
	switch b.remainder {
	case "january":
		b.Text = "Garnet"
	case "february":
		b.Text = "Amethyst"
	case "march":
		b.Text = "Aquamarine, Bloodstone"
	case "april":
		b.Text = "Diamond"
	case "may":
		b.Text = "Emerald"
	case "june":
		b.Text = "Pearl, Moonstone, Alexandrite"
	case "july":
		b.Text = "Ruby"
	case "august":
		b.Text = "Peridot, Spinel"
	case "september":
		b.Text = "Sapphire"
	case "october":
		b.Text = "Opal, Tourmaline"
	case "november":
		b.Text = "Topaz, Citrine"
	case "december":
		b.Text = "Turquoise, Zircon, Tanzanite"
	}

	return b
}

func (b *BirthStone) setCache() answerer {
	b.Cache = true
	return b
}

func (b *BirthStone) tests() []test {
	typ := "birthstone"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query: "January birthstone",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Garnet",
					Cache:        true,
				},
			},
		},
		test{
			query: "birthstone february",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Amethyst",
					Cache:        true,
				},
			},
		},
		test{
			query: "march birth stone",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Aquamarine, Bloodstone",
					Cache:        true,
				},
			},
		},
		test{
			query: "birth stone April",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Diamond",
					Cache:        true,
				},
			},
		},
		test{
			query: "birth stones may",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Emerald",
					Cache:        true,
				},
			},
		},
		test{
			query: "birthstones June",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Pearl, Moonstone, Alexandrite",
					Cache:        true,
				},
			},
		},
		test{
			query: "July Birth Stones",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Ruby",
					Cache:        true,
				},
			},
		},
		test{
			query: "birthstones August",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Peridot, Spinel",
					Cache:        true,
				},
			},
		},
		test{
			query: "september birthstones",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Sapphire",
					Cache:        true,
				},
			},
		},
		test{
			query: "October birthstone",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Opal, Tourmaline",
					Cache:        true,
				},
			},
		},
		test{
			query: "birthstone November",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Topaz, Citrine",
					Cache:        true,
				},
			},
		},
		test{
			query: "December birthstone",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Turquoise, Zircon, Tanzanite",
					Cache:        true,
				},
			},
		},
	}

	return tests
}
