// Package vote handles storing and retrieving user votes on urls
package vote

import (
	"errors"
	"jivesearch/search/document"
	"time"
)

// Vote represents a user's vote for a url
type Vote struct {
	Query  string `db:"query" json:"query"`
	URL    string `db:"url" json:"url"`
	Domain string `db:"domain" json:"-"`
	Vote   int    `db:"vote" json:"vote"`
	Date   string `db:"date" json:"date"`
}

// Result represents the vote total for a url
type Result struct {
	URL   string `db:"url" json:"url"`
	Votes int    `db:"votes" json:"votes"`
}

// Voter outlines the methods to store & retrieve votes
type Voter interface {
	Setup() error
	Get(query string, limit int) ([]Result, error)
	Insert(v *Vote) error
}

// Option is a function option
type Option func(v *Vote) (*Vote, error)

var (
	now = func() string { return time.Now().Format("20060102") }
	// ErrInvalidQuery indicates a malformed query was passed
	ErrInvalidQuery = errors.New("invalid query")
	// ErrInvalidURL indicates a malformed URL was passed
	ErrInvalidURL = errors.New("invalid url")
	// ErrInvalidVote indicates an invalid vote was passed
	ErrInvalidVote = errors.New("invalid vote; must be -1 or 1")
)

// New creates a new Vote, setting the date to today
func New(opts ...Option) (*Vote, error) {
	v := &Vote{
		Date: now(),
	}

	for _, opt := range opts {
		if _, err := opt(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

// Query sets the query
func Query(q string) Option {
	return func(v *Vote) (*Vote, error) {
		if q == "" {
			return nil, ErrInvalidQuery
		}

		v.Query = q
		return v, nil
	}
}

// URL sets the URL and domain the user is upvoting/downvoting
func URL(lnk string) Option {
	return func(v *Vote) (*Vote, error) {
		u, err := document.ValidateURL(lnk)
		if err != nil {
			return nil, ErrInvalidURL
		}

		v.URL = u.String()

		v.Domain, err = document.ExtractDomain(u)
		if err != nil {
			return nil, ErrInvalidURL
		}

		return v, nil
	}
}

// Value sets the vote value
func Value(vote int) Option {
	return func(v *Vote) (*Vote, error) {
		if vote != -1 && vote != 1 {
			return nil, ErrInvalidVote
		}

		v.Vote = vote
		return v, nil
	}
}
