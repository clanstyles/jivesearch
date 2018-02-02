// Package search provides the core search results.
package search

import (
	"math"
	"strconv"

	"github.com/jivesearch/jivesearch/search/document"
	"github.com/jivesearch/jivesearch/search/vote"
	"golang.org/x/text/language"
)

// Fetcher outlines the methods used to retrieve the core search results
type Fetcher interface {
	Fetch(q string, lang language.Tag, region language.Region, number int, page int, votes []vote.Result) (*Results, error)
}

// Results are the core search results from a query
type Results struct {
	Count      int64                `json:"count"`
	Page       string               `json:"page"`
	Previous   string               `json:"previous"`
	Next       string               `json:"next"`
	Last       string               `json:"last"`
	Pagination []string             `json:"-"`
	Documents  []*document.Document `json:"links"`
}

// AddPagination adds pagination to the search results
func (r *Results) AddPagination(number, page int) *Results {
	r.Pagination = []string{}
	r.Page = strconv.Itoa(page)
	if page > 1 {
		r.Previous = strconv.Itoa(page - 1)
	}

	min, max := 1, int(math.Ceil(float64(r.Count)/float64(number))) // round up
	if page > max {
		page = max
	}

	if max > page {
		r.Next = strconv.Itoa(page + 1)
	}

	if page > 6 {
		min = page - 5
		tmp := max
		for i := 0; i < 5; i++ {
			if tmp < max {
				max = page + i
			}
		}
	}

	for i := min; i <= max; i++ {
		if len(r.Pagination) < 10 {
			r.Pagination = append(r.Pagination, strconv.Itoa(i))
		}
	}

	return r
}
