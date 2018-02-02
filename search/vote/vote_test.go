package vote

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type want struct {
		vote *Vote
		err  error
	}

	for _, c := range []struct {
		name string
		now  string
		q    string
		u    string
		v    int
		want
	}{
		{
			"upvote", "20120102", "a simple search", "https://www.example.com", 1,
			want{
				&Vote{
					Date:   "20120102",
					Query:  "a simple search",
					URL:    "https://www.example.com",
					Domain: "example.com",
					Vote:   1,
				},
				nil,
			},
		},
		{
			"downvote", "20141130", "another simple search", "http://example2.com", -1,
			want{
				&Vote{
					Date:   "20141130",
					Query:  "another simple search",
					URL:    "http://example2.com",
					Domain: "example2.com",
					Vote:   -1,
				},
				nil,
			},
		},
		{
			"invalid vote 0", "20170102", "very stupid search", "https://www.example.com", 0,
			want{
				nil,
				ErrInvalidVote,
			},
		},
		{
			"invalid vote 2", "20160929", "another query", "https://www.example.com/", 2,
			want{
				nil,
				ErrInvalidVote,
			},
		},
		{
			"invalid vote -2", "20160929", "another query", "https://www.example.com/", -2,
			want{
				nil,
				ErrInvalidVote,
			},
		},
		{
			"blank query", "20170811", "", "https://www.example.com", 1,
			want{
				nil,
				ErrInvalidQuery,
			},
		},
		{
			"blank url", "20160929", "a random query", "", 1,
			want{
				nil,
				ErrInvalidURL,
			},
		},
		{
			"invalid url (no colon)", "20160929", "another random query", "http//www.example.com", 1,
			want{
				nil,
				ErrInvalidURL,
			},
		},
		{
			"invalid url (missing scheme)", "20160929", "another random query", "www.example.com", 1,
			want{
				nil,
				ErrInvalidURL,
			},
		},
		{
			"invalid url (wrong scheme)", "20160929", "another random query", "ftp://www.example.com", 1,
			want{
				nil,
				ErrInvalidURL,
			},
		},
		{
			"invalid url (relative url)", "20160929", "yet another random query", "/this-is-a-relative-url", 1,
			want{
				nil,
				ErrInvalidURL,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			now = func() string { return c.now }
			got, err := New(
				Query(c.q),
				URL(c.u),
				Value(c.v),
			)

			if !reflect.DeepEqual(err, c.want.err) {
				t.Fatalf("got err %v; want %v", err, c.want.err)
			}

			if !reflect.DeepEqual(got, c.want.vote) {
				t.Fatalf("got %+v; want %+v", got, c.want.vote)
			}
		})
	}
}
