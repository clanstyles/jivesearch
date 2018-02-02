package wikipedia

import (
	"bytes"
	"compress/gzip"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/dsnet/compress/bzip2"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/afero"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	"golang.org/x/text/language"
)

var enwikiFile = &File{
	language: language.English,
	Base:     "enwiki-20171218-cirrussearch-content.json.gz",
}

func TestFile_Download(t *testing.T) {
	tests := []struct {
		name string
		u    string
		f    *File
	}{
		{
			"enwiki",
			"https://dumps.wikimedia.org/other/cirrussearch/current/enwiki-20171218-cirrussearch-content.json.gz",
			enwikiFile,
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	oldFs := fs

	mfs := afero.NewMemMapFs()
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	for _, tt := range tests {
		responder := httpmock.NewStringResponder(200, `<html><body></body></html>`)
		httpmock.RegisterResponder("GET", tt.u, responder)

		var err error
		tt.f.URL, err = url.Parse(tt.u)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Download(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFile_Parse(t *testing.T) {
	type args struct {
		truncate int
	}

	tests := []struct {
		name string
		path string
		lang language.Tag
		args args
	}{
		{
			"enwiki",
			"enwiki-20171218-cirrussearch-content.json.gz",
			language.English,
			args{10},
		},
		{
			"wikidata",
			"latest-all.json.bz2",
			language.English,
			args{},
		},
	}

	oldFs := fs
	mfs := afero.NewMemMapFs()
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFile(&url.URL{Path: tt.path}, tt.lang)

			md := &mockDumper{}
			f.SetABS("/path/to/nowhere/").SetDumper(md)

			// just need to minimize for the tests
			m := minify.New()
			m.AddFunc("text/javascript", js.Minify)
			s, err := m.String("text/javascript", shaqRawClaims)
			if err != nil {
				t.Fatal(err)
			}

			var b bytes.Buffer

			switch strings.HasSuffix(tt.path, ".bz2") {
			case true:
				// go still doesn't have an "official" .bz2 writer https://github.com/golang/go/issues/4828.
				// until then we use github.com/dsnet/compress/bzip2 (which will most likely at some point be merged into standard lib).
				bz2w, err := bzip2.NewWriter(&b, nil)
				if err != nil {
					t.Fatal(err)
				}
				bz2w.Write([]byte(s))
				bz2w.Close()

			default:
				w := gzip.NewWriter(&b)
				w.Write([]byte(s))
				w.Close()
			}

			fs.MkdirAll(f.Dir, 0755)
			afero.WriteFile(mfs, f.ABS, b.Bytes(), 0644)

			if err := f.Parse(tt.args.truncate); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCirrusLinks(t *testing.T) {
	type args struct {
		supported []language.Tag
	}

	tests := []struct {
		name string
		args args
		urls []string
		want []*File
	}{
		{
			"basic",
			args{[]language.Tag{language.English}},
			[]string{
				"https://dumps.wikimedia.org/other/cirrussearch/current/enwiki-20171218-cirrussearch-content.json.gz",
			},
			[]*File{
				enwikiFile,
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
			httpmock.RegisterResponder("GET", CirrusURL.String(), responder)

			got, err := CirrusLinks(tt.args.supported)
			if err != nil {
				t.Fatal(err)
			}

			for i, u := range tt.urls {
				tt.want[i].URL, err = url.Parse(u)
				if err != nil {
					t.Fatal(err)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

type mockDumper struct{}

func (md *mockDumper) Dump(wikidata bool, lang language.Tag, rows chan interface{}) error {
	for _ = range rows {
	}

	return nil
}
