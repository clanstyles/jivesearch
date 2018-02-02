package instant

import (
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jivesearch/jivesearch/instant/contributors"
)

// Stats is an instant answer that
// returns the average, median, etc.
type Stats struct {
	Answer
}

var reStats *regexp.Regexp

func (s *Stats) setQuery(r *http.Request) answerer {
	s.Answer.setQuery(r)
	return s
}

func (s *Stats) setUserAgent(r *http.Request) answerer {
	return s
}

func (s *Stats) setType() answerer {
	s.Type = "stats"
	return s
}

func (s *Stats) setContributors() answerer {
	s.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return s
}

func (s *Stats) setTriggers() answerer {
	s.triggers = []string{
		"avg", "average", "mean", "median", "sum", "total",
	}
	return s
}

func (s *Stats) setTriggerFuncs() answerer {
	s.triggerFuncs = []triggerFunc{
		startsWith, endsWith,
	}
	return s
}

func (s *Stats) setSolution() answerer {
	// get all the numbers..this regexp will correctly grab e notation
	numbersStrings := reStats.FindAllString(s.remainder, -1)
	numbers := []float64{}

	for _, n := range numbersStrings {
		if i, err := strconv.ParseFloat(n, 64); err == nil {
			numbers = append(numbers, i)
		}
	}

	for _, a := range s.triggers {
		if !strings.Contains(s.query, a) {
			continue
		}

		var txt string
		var ans float64

		if a == "avg" || a == "average" || a == "mean" {
			txt = "Average: "
			ans = average(numbers)
		} else if a == "median" {
			txt = "Median: "
			ans = median(numbers)
		} else if a == "sum" || a == "total" {
			txt = "Sum: "
			ans = sum(numbers)
		}
		s.Text = txt + strconv.FormatFloat(ans, 'f', -1, 64)
		break
	}

	return s
}

func (s *Stats) setCache() answerer {
	s.Cache = true
	return s
}

func (s *Stats) tests() []test {
	typ := "stats"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		test{
			query: "avg 3 4e6",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Average: 2000001.5",
					Cache:        true,
				},
			},
		},
		test{
			query: "11 18 -142 Average",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Average: -37.666666666666664",
					Cache:        true,
				},
			},
		},
		test{
			query: "6 3 -5 23 Median",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Median: 4.5",
					Cache:        true,
				},
			},
		},
		test{
			query: "median 17 12 -18",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Median: 12",
					Cache:        true,
				},
			},
		},
		test{
			query: "58 96 -41 sum",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Sum: 113",
					Cache:        true,
				},
			},
		},
		test{
			query: "Total -17 3 87 -476",
			expected: []Solution{
				Solution{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Text:         "Sum: -403",
					Cache:        true,
				},
			},
		},
	}

	return tests
}

func average(numbers []float64) float64 {
	total := sum(numbers)
	return total / float64(len(numbers))
}

func median(numbers []float64) float64 {
	sort.Float64s(numbers)
	middle := len(numbers) / 2
	result := numbers[middle]
	if len(numbers)%2 == 0 {
		result = (result + numbers[middle-1]) / 2
	}
	return result
}

func sum(numbers []float64) float64 {
	var total float64
	for _, value := range numbers {
		total += value
	}
	return total
}

func init() {
	reStats = regexp.MustCompile(`[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`)
}
