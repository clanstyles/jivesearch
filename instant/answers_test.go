package instant

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// TestDetect runs the test cases for each instant answer.
func TestDetect(t *testing.T) {
	cases := []test{
		test{
			query: "testing an empty answer here",
			expected: []Solution{
				Solution{},
			},
		},
	}

	for i, ia := range answers() {
		if len(ia.tests()) == 0 {
			t.Fatalf("No tests for answer #%d", i)
		}
		cases = append(cases, ia.tests()...)
	}

	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			ctx := fmt.Sprintf(`(query: %q, user agent: %q)`, c.query, c.userAgent)

			v := url.Values{}
			v.Set("q", c.query)

			r := &http.Request{
				Form:   v,
				Header: make(http.Header),
			}

			r.Header.Set("User-Agent", c.userAgent)

			got := Detect(r)

			var solved bool

			for _, expected := range c.expected {
				if reflect.DeepEqual(got, expected) {
					solved = true
					break
				}
			}

			if !solved {
				t.Errorf("Instant answer failed %v", ctx)
				t.Errorf("got %+v;", got)
				t.Errorf("want ")
				for _, expected := range c.expected {
					t.Errorf("    %+v\n", expected)
				}
				t.FailNow()
			}
		})
	}
}
