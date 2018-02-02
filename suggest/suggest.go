// Package suggest handles AutoComplete and Phrase Suggester (Did you mean?) queries
package suggest

// Suggester outlines methods to fetch & store Autocomplete & PhraseSuggester results
type Suggester interface {
	IndexExists() (bool, error)
	Setup() error
	Exists(q string) (bool, error)
	Insert(q string) error
	Increment(q string) error
	Completion(q string, size int) (Results, error)
	//phrase(q string) Results //  TODO: "Did you mean?"
}

// Results are the results of an autocomplete query
type Results struct { // remember top-level arrays = no-no in javascript/json
	RawQuery    string   `json:"-"`
	Suggestions []string `json:"suggestions"`
}
