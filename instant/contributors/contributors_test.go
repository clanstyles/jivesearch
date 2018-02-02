package contributors

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	Contributors.M = make(map[string]Contributor)
	Contributors.M["sethklarman"] = Contributor{Name: "Seth Klarman", Github: "sethklarman"}
	Contributors.M["marilynmonroe"] = Contributor{Name: "Marilyn Monroe", Github: "Marilyn Monroe"}

	for _, c := range []struct {
		name  string
		names []string
		want  []Contributor
	}{
		{
			name:  "match",
			names: []string{"sethklarman", "davidabrams"},
			want: []Contributor{
				Contributor{Name: "Seth Klarman", Github: "sethklarman"},
			},
		},
		{
			name:  "empty",
			names: []string{"jimihendrix"},
			want:  []Contributor{},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var got = Load(c.names)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}
