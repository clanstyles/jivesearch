// Package bangs detects when queries should be redirected to 3rd party sites
package bangs

import (
	"fmt"
	"strings"
	"sync"
)

// Bangs holds a map of !bangs
type Bangs struct {
	sync.Mutex
	M map[string]map[string]string
}

// Bang holds a single !bang
type Bang struct {
	Triggers []string
	Name     string
	Category string
	Regions  []Region
}

// Region holds the regional information and url of a !bang
type Region struct {
	Region   string
	Location string
}

var def = "default"

// New creates a pointer with the default !bangs.
// Use default url unless a region is provided.
// Region: US, Language: French !a ---> Amazon.com
// Region: France, Language: English !a ---> Amazon.fr
// !afr ---> Amazon.fr
// Note: Some !bangs don't respect the language passed in or
// may not support it (eg they may support pt but not pt-BR)
//
// TODO: Allow overrides...perhaps add a method or use a config.
// Note: If we end up using viper for this don't use "SetDefault"
// as overriding one !bang will replace ALL !bangs. Instead, use "Set".
func New() *Bangs {
	// Not sure about the structure here...slice of Bangs makes it easy to add bangs
	// Would like to add autocomplete feature so that people can find !bangs easier.
	b := &Bangs{
		M: make(map[string]map[string]string),
	}

	bngs := []Bang{
		Bang{
			[]string{"a", "amazon"},
			"Amazon", "shopping",
			[]Region{
				Region{def, "https://www.amazon.com/s/ref=nb_sb_noss?url=search-alias%3Daps&field-keywords={{{term}}}"},
				Region{"ca", "https://www.amazon.ca/s/ref=nb_sb_noss?url=search-alias%3Daps&field-keywords={{{term}}}"},
				Region{"fr", "https://www.amazon.fr/s/ref=nb_sb_noss?url=search-alias%3Daps&field-keywords={{{term}}}"},
				Region{"uk", "https://www.amazon.co.uk/s/ref=nb_sb_noss?url=search-alias%3Daps&field-keywords={{{term}}}"},
			},
		},
		Bang{
			[]string{"g", "google"},
			"Google", "search",
			[]Region{
				Region{def, "https://encrypted.google.com/search?hl={{{lang}}}&q={{{term}}}"},
				Region{"ca", "https://www.google.ca/search?q={{{term}}}"},
				Region{"fr", "https://www.google.fr/search?hl={{{lang}}}&q={{{term}}}"},
				Region{"ru", "https://www.google.ru/search?hl={{{lang}}}&q={{{term}}}"},
			},
		},
		Bang{
			[]string{"gfr", "googlefr"},
			"Google France", "search",
			[]Region{
				Region{def, "https://www.google.fr/search?hl={{{lang}}}&q={{{term}}}"},
			},
		},
		Bang{
			[]string{"gru", "googleru"},
			"Google Russia", "search",
			[]Region{
				Region{def, "https://www.google.ru/search?hl={{{lang}}}&q={{{term}}}"},
			},
		},
		Bang{
			[]string{"reddit"},
			"Reddit", "social media",
			[]Region{
				Region{def, "https://www.reddit.com/search?q={{{term}}}&restrict_sr=&sort=relevance&t=all"},
			},
		},
	}

	// create a map for faster lookups
	for _, bng := range bngs {
		for _, t := range bng.Triggers {
			if _, ok := b.M[t]; ok {
				panic(fmt.Sprintf("duplicate trigger found %v", t))
			}
			b.M[t] = make(map[string]string)
			for _, r := range bng.Regions {
				b.M[t][r.Region] = r.Location
			}
		}
	}

	return b
}

// Detect lets us know if we have a !bang match.
func (b *Bangs) Detect(q, region, language string) (string, bool) {
	b.Lock()
	defer b.Unlock()

	fields := strings.Fields(q)

	for i, field := range fields {
		if field == "!" || (!strings.HasPrefix(field, "!") && !strings.HasSuffix(field, "!")) {
			continue
		}

		if bng, ok := b.M[strings.ToLower(strings.Trim(field, "!"))]; ok { // find the bang
			for _, reg := range []string{strings.ToLower(region), def} { // use default region if no region specified
				if u, ok := bng[reg]; ok {
					remainder := strings.Join(append(fields[:i], fields[i+1:]...), " ")
					u = strings.Replace(u, "{{{term}}}", remainder, -1)
					return strings.Replace(u, "{{{lang}}}", language, -1), true
				}
			}
		}
	}
	return "", false
}
