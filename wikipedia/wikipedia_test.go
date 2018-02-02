package wikipedia

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

var shaqWikipedia = &Wikipedia{
	"Q169452", "en", "Shaquille O'Neal", `Shaquille Rashaun O'Neal, Ed.D, nicknamed "Shaq", is an American retired professional ...`, 90,
}

var shaqRawWikipediaJSON = []byte(`{"wikibase_item": "Q169452", "language": "en", "title": "Shaquille O'Neal",	"text": "Shaquille Rashaun O'Neal, Ed.D (born March 6, 1972), nicknamed \"Shaq\" (SHAK), is an American retired professional basketball player and rapper, currently serving as a sports analyst on the television program Inside the"}`)

func TestWikipedia_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name     string
		truncate int
		args     args
		want     *Wikipedia
	}{
		{
			name:     "shaq",
			truncate: 90,
			args:     args{shaqRawWikipediaJSON},
			want:     shaqWikipedia,
		},
		{
			name:     "madonna",
			truncate: 50,
			args: args{
				[]byte(`{"wikibase_item": "Q1744", "language": "en", "title": "Madonna", "text": "Madonna Louise Ciccone (born August 16, 1958) is an American singer, songwriter, actress, and businesswoman. Referred to as the \"Queen of Pop\" since the 1980s, Madonna is known for pushing the boundaries of lyrical content in mainstream popular music, as well as visual imagery in music videos"}`),
			},
			want: &Wikipedia{
				ID:       "Q1744",
				Language: "en",
				Title:    "Madonna",
				Text:     "Madonna Louise Ciccone is an American singer, ...",
				truncate: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Wikipedia{truncate: tt.truncate}

			if err := got.UnmarshalJSON(tt.args.data); err != nil {
				t.Errorf("Wikipedia.UnmarshalJSON() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v; want %+v", got, tt.want)
			}
		})
	}
}

func TestLanguages(t *testing.T) {
	type want struct {
		supported   int
		unsupported int
	}

	for _, c := range []struct {
		name string
		arg  []language.Tag
		want
	}{
		{"basic", []language.Tag{}, want{len(Available), 0}},
		{"en", []language.Tag{language.English}, want{1, 0}},
		{"rejected", []language.Tag{language.Croatian, language.BrazilianPortuguese}, want{1, 1}},
	} {
		t.Run(c.name, func(t *testing.T) {
			supported, unsupported := Languages(c.arg)

			if len(supported) != c.want.supported {
				t.Errorf("got %v, want %v", supported, c.want.supported)
			}

			if len(unsupported) != c.want.unsupported {
				t.Errorf("got %v, want %v", supported, c.want.unsupported)
			}
		})
	}
}
