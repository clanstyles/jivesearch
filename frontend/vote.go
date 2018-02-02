package frontend

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jivesearch/jivesearch/search/vote"
)

func (f *Frontend) voteHandler(w http.ResponseWriter, r *http.Request) *response {
	fail := &response{
		status:   http.StatusInternalServerError,
		template: "json",
	}

	q := strings.TrimSpace(r.FormValue("q"))
	u := strings.TrimSpace(r.FormValue("u"))
	value := strings.TrimSpace(r.FormValue("v"))

	val, err := strconv.Atoi(value)
	if err != nil {
		fail.err = vote.ErrInvalidVote
		fail.status = http.StatusBadRequest
		return fail
	}

	v := &vote.Vote{}
	v, fail.err = vote.New(
		vote.Query(q),
		vote.URL(u),
		vote.Value(val),
	)

	if fail.err != nil {
		if fail.err == vote.ErrInvalidQuery || fail.err == vote.ErrInvalidURL || fail.err == vote.ErrInvalidVote {
			fail.status = http.StatusBadRequest
		}

		return fail
	}

	if fail.err = f.Vote.Insert(v); fail.err != nil {
		return fail
	}

	return &response{
		status:   http.StatusOK,
		template: "json",
		data:     "success",
	}
}
