package main

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/jivesearch/jivesearch/wikipedia"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
	}{
		{"basic", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()

			_, err := setup(v)
			if err != nil {
				t.Fatal(err)
			}
		})
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

func TestFiles(t *testing.T) {
	type args struct {
		wikipedia bool
		wikidata  bool
		supported []language.Tag
	}

	tests := []struct {
		name string
		args
		urls map[language.Tag]string
	}{
		{
			"wikipedia",
			args{
				true, false, []language.Tag{language.English},
			},
			map[language.Tag]string{
				language.English: "enwiki-20171218-cirrussearch-content.json.gz",
			},
		},
		{
			"wikidata",
			args{
				false, true, []language.Tag{language.English},
			},
			map[language.Tag]string{
				language.English: wikipedia.WikiDataURL.String(),
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responder := httpmock.NewStringResponder(
				200,
				`<html>
					<body>
					<a href="aawiki-20171218-cirrussearch-content.json.gz">aawiki-20171218-cirrussearch-content.json.gz</a>18-Dec-2017 16:15 2158
					<a href="aawiki-20171218-cirrussearch-general.json.gz">aawiki-20171218-cirrussearch-general.json.gz</a>18-Dec-2017 16:15 164287
					<a href="enwiki-20171218-cirrussearch-content.json.gz">enwiki-20171218-cirrussearch-content.json.gz</a>19-Dec-2017 10:33 25078247011
					<a href="enwiki-20171218-cirrussearch-general.json.gz">enwiki-20171218-cirrussearch-general.json.gz</a>19-Dec-2017 15:25 43605620413
					<a href="usabilitywiki-20171218-cirrussearch-content.json.gz">usabilitywiki-20171218-cirrussearch-content.jso..&gt;</a>20-Dec-2017 12:56 386462
					<a href="usabilitywiki-20171218-cirrussearch-general.json.gz">usabilitywiki-20171218-cirrussearch-general.jso..&gt;</a>20-Dec-2017 12:56 813441
					</body>
				</html>`,
			)
			httpmock.RegisterResponder("GET", wikipedia.CirrusURL.String(), responder)

			v := viper.New()
			v.SetDefault("wikipedia.text", tt.args.wikipedia)
			v.SetDefault("wikipedia.data", tt.args.wikidata)

			want := []*wikipedia.File{}

			for k, v := range tt.urls {
				u, err := url.Parse(v)
				if err != nil {
					t.Fatal(err)
				}

				want = append(want, wikipedia.NewFile(wikipedia.CirrusURL.ResolveReference(u), k))
			}

			got, err := files(v, tt.args.supported)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %+v, want %+v", got, want)
			}
		})
	}
}
