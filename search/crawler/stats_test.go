package crawler

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	s := &Stats{
		StatusCodes: map[int]int64{
			-1:  25,
			200: 1000,
			400: 50,
			500: 10,
		},
		elapsed: 5 * time.Second,
	}

	want := "[stats] Crawled: 1085 Elapsed: 5s\n[stats] Rate: 217 per second, 13,020 per minute, 781,200 per hour, 18,748,800 per day\n[stats]2xx (92%)  4xx (4%)  5xx (0%)  Not Crawled (2%)  \n"
	got := s.String()

	if got != want {
		t.Fatalf("got %q; want %q", got, want)
	}
}
