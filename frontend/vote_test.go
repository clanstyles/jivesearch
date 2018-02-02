package frontend

import (
	"jivesearch/search/vote"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestVoteHandler(t *testing.T) {
	for _, c := range []struct {
		name  string
		query string
		url   string
		vote  string
		want  *response
	}{
		{
			"basic", "a simple search", "http://www.example.com", "1",
			&response{
				status:   http.StatusOK,
				template: "json",
				data:     "success",
			},
		},
		{
			"invalid vote", "a simple search", "http://www.example.com", "wrong",
			&response{
				status:   http.StatusBadRequest,
				template: "json",
				err:      vote.ErrInvalidVote,
			},
		},
		{
			"downvote", "some query", "https://www.example.com", "-1",
			&response{
				status:   http.StatusOK,
				template: "json",
				data:     "success",
			},
		},
		{
			"invalid vote 0", "another query", "https://www.example.com/", "0",
			&response{
				status:   http.StatusBadRequest,
				template: "json",
				err:      vote.ErrInvalidVote,
			},
		},
		{
			"invalid vote 100", "yet another query", "https://www.example.com", "100",
			&response{
				status:   http.StatusBadRequest,
				template: "json",
				err:      vote.ErrInvalidVote,
			},
		},
		{
			"blank query", "", "https://www.example.com", "1",
			&response{
				status:   http.StatusBadRequest,
				template: "json",
				err:      vote.ErrInvalidQuery,
			},
		},
		{
			"blank url", "some query", "", "1",
			&response{
				status:   http.StatusBadRequest,
				template: "json",
				err:      vote.ErrInvalidURL,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &Frontend{
				Vote: &mockVoter{},
			}

			req, err := http.NewRequest("GET", "/vote", nil)
			if err != nil {
				t.Fatal(err)
			}

			q := req.URL.Query()
			q.Add("q", c.query)
			q.Add("u", c.url)
			q.Add("v", c.vote)

			req.URL.RawQuery = q.Encode()

			got := f.voteHandler(httptest.NewRecorder(), req)

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}

type mockVoter struct{}

func (m *mockVoter) Get(q string, l int) ([]vote.Result, error) {
	res := []vote.Result{}
	return res, nil
}
func (m *mockVoter) Setup() error {
	return nil
}
func (m *mockVoter) Insert(v *vote.Vote) error {
	return nil
}
