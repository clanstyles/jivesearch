package bangs

import "testing"

// TestDefault tests that each !bang has a default location
func TestDefault(t *testing.T) {
	b := New()
	for trigger, bng := range b.M {
		if _, ok := bng[def]; !ok {
			t.Fatalf("!%v bang needs a default region", trigger)
		}
	}
}

func TestDetect(t *testing.T) {
	type data struct {
		loc string
		ok  bool
	}

	for _, c := range []struct {
		q    string
		r    string
		l    string
		want data
	}{
		{
			q: "!g bob", r: "en", l: "fr",
			want: data{
				loc: "https://encrypted.google.com/search?hl=fr&q=bob",
				ok:  true,
			},
		},
		{
			q: "!g bob french", r: "fr", l: "en",
			want: data{
				loc: "https://www.google.fr/search?hl=en&q=bob french",
				ok:  true,
			},
		},
		{
			q: "!gfr something french", r: "fr", l: "en",
			want: data{
				loc: "https://www.google.fr/search?hl=en&q=something french",
				ok:  true,
			},
		},
		{
			q: "nonexistent! some query", r: "en", l: "fr",
			want: data{
				loc: "",
				ok:  false,
			},
		},
		{
			q: "this is not a bang", r: "en", l: "en",
			want: data{
				loc: "",
				ok:  false,
			},
		},
		{
			q: "this is not a bang g", r: "en", l: "en",
			want: data{
				loc: "",
				ok:  false,
			},
		},
		{
			q: "this is not a bang google", r: "en", l: "en",
			want: data{
				loc: "",
				ok:  false,
			},
		},
	} {
		t.Run(c.q, func(t *testing.T) {
			b := New()

			var got = data{}
			got.loc, got.ok = b.Detect(c.q, c.r, c.l)
			if got != c.want {
				t.Fatalf("got %+v; want %+v", got, c.want)
			}
		})
	}
}
