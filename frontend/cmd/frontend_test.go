package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/jivesearch/jivesearch/frontend"
	"github.com/jivesearch/jivesearch/wikipedia"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func TestSetup(t *testing.T) {
	var parsed bool
	frontend.ParseTemplates = func() {
		parsed = true
	}

	v := viper.New()
	s := setup(v)

	if !parsed {
		t.Fatal("expected templates to be parsed. they weren't.")
	}

	if reflect.DeepEqual(httpClient.Timeout, time.Time{}) {
		t.Fatal("expected http client to have a timeout. it doesn't")
	}

	want := ":8000"
	if s.Addr != want {
		t.Fatalf("got %q; want %q", s.Addr, want)
	}

	if s.Handler == nil {
		t.Fatalf("got nil handler")
	}
}

func TestLanguages(t *testing.T) {
	type want struct {
		supported   []language.Tag
		unsupported int
	}

	tests := []struct {
		name string
		args []string
		want
	}{
		{"default", []string{}, want{[]language.Tag{}, 0}},
		{"en", []string{"en"}, want{[]language.Tag{language.English}, 0}},
		{"custom", []string{"fr", "de"}, want{[]language.Tag{language.French, language.German}, 0}},
		{"unsupported", []string{"en", "def"}, want{[]language.Tag{language.English}, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetDefault("languages", tt.args)

			supported, unsupported := languages(v)
			if len(unsupported) != tt.want.unsupported {
				t.Errorf("got %d unsupported, want %d", len(unsupported), tt.want.unsupported)
			}

			switch len(tt.args) {
			case 0:
				if len(supported) != len(wikipedia.Available) {
					t.Errorf("got %+v, want %+v", supported, tt.want.supported)
				}
			default:
				if !reflect.DeepEqual(supported, tt.want.supported) {
					t.Errorf("got %+v, want %+v", supported, tt.want.supported)
				}
			}
		})
	}
}
