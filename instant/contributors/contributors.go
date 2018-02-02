// Package contributors contains the contributors of the instant answers
package contributors

import "sync"

// Contributor is an individual contributor for an instant answer, etc...
type Contributor struct {
	Name    string `json:"contributor"`
	Github  string `json:"github"`
	Twitter string `json:"twitter"`
}

type contributors struct {
	sync.Mutex
	M map[string]Contributor
}

// Contributors is a map of all the contributors to IA's.
var Contributors = contributors{}

// Load the contributor info
func Load(names []string) []Contributor {
	c := []Contributor{}

	// probably not necessary but haven't seen a definitive answer
	Contributors.Lock()
	defer Contributors.Unlock()

	for _, n := range names {
		if val, ok := Contributors.M[n]; ok {
			c = append(c, val)
		}
	}

	return c
}

func init() {
	// add your info here
	Contributors.M = make(map[string]Contributor)
	Contributors.M["brentadamson"] = Contributor{Name: "Brent Adamson", Github: "brentadamson", Twitter: "thebrentadamson"}
}
