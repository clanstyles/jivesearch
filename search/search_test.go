package search

import (
	"reflect"
	"testing"
)

func TestAddPagination(t *testing.T) {
	for _, c := range []struct {
		name   string
		count  int64
		number int
		page   int
		want   []string
	}{
		{
			name:   "basic",
			count:  250,
			number: 10,
			page:   1,
			want:   []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
		{
			name:   "page 7",
			count:  500,
			number: 10,
			page:   7,
			want:   []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
		},
		{
			name:   "maxxed out",
			count:  210,
			number: 10,
			page:   22,
			want:   []string{"16", "17", "18", "19", "20", "21"},
		},
		{
			name:   "middle",
			count:  250,
			number: 10,
			page:   22,
			want:   []string{"17", "18", "19", "20", "21", "22", "23", "24", "25"},
		},
		{
			name:   "few results",
			count:  17,
			number: 10,
			page:   1,
			want:   []string{"1", "2"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			res := &Results{
				Count: c.count,
			}

			got := res.AddPagination(c.number, c.page)

			if !reflect.DeepEqual(got.Pagination, c.want) {
				t.Fatalf("got %+v; want %+v", got.Pagination, c.want)
			}
		})
	}
}
