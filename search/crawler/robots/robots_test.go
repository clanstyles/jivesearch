package robots

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	for _, c := range []struct {
		sh string
	}{
		{
			"http://example.com",
		},
		{
			"https://www.example.com",
		},
	} {
		t.Run(c.sh, func(t *testing.T) {
			got := New(c.sh)

			if got.SchemeHost != c.sh {
				t.Fatalf("got %v; want %v", got.SchemeHost, c.sh)
			}
		})
	}
}

// NOTE: In the absence of a time zone indicator, time.Parse returns a time in UTC
func TestSetExpires(t *testing.T) {
	for _, c := range []struct {
		name       string
		host       string
		statusCode int
		now        time.Time
		want       string
	}{
		{
			name:       "http://example.com",
			host:       "http://example.com",
			statusCode: 200,
			now:        time.Date(2015, 12, 15, 4, 35, 35, 10, time.UTC),
			want:       "201512160435",
		},
		{
			name:       "https://www.example.com",
			host:       "https://www.example.com",
			statusCode: 500,
			now:        time.Date(2019, 12, 31, 23, 6, 47, 6, time.UTC),
			want:       "202001010006",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			now = func() time.Time { return c.now }
			got := New(c.host).SetStatusCode(c.statusCode).SetExpires()

			if !reflect.DeepEqual(got.Expires, c.want) {
				t.Fatalf("got %q; want %q", got.Expires, c.want)
			}
		})
	}
}

func TestExpired(t *testing.T) {
	for _, c := range []struct {
		name    string
		host    string
		expires string
		now     time.Time
		want    bool
		err     error
	}{
		{
			name:    "not expired",
			host:    "https://example.com",
			expires: "201410020359",
			now:     time.Date(2014, 10, 2, 3, 17, 18, 31, time.UTC),
			want:    false,
			err:     nil,
		},
		{
			name:    "expired",
			host:    "http://www.example.com",
			expires: "201905311657",
			now:     time.Date(2019, 05, 31, 16, 58, 01, 0, time.UTC),
			want:    true,
			err:     nil,
		},
		{
			name:    "parsing error",
			host:    "http://www.example.com",
			expires: "20191230",
			now:     time.Date(2019, 12, 31, 23, 6, 47, 6, time.UTC),
			want:    true,
			err: &time.ParseError{
				Value:      "20191230",
				Layout:     "200601021504",
				ValueElem:  "",
				LayoutElem: "15",
				Message:    "",
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			now = func() time.Time { return c.now }
			rbts := New(c.host)
			rbts.Expires = c.expires
			got, err := rbts.Expired()

			if !reflect.DeepEqual(err, c.err) {
				t.Fatalf("got err %q; want %q", err, c.err)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %t; want %t", got, c.want)
			}
		})
	}
}

func TestSetStatusCode(t *testing.T) {
	for _, c := range []struct {
		name       string
		host       string
		statusCode int
	}{
		{
			name:       "ok",
			host:       "http://www.example.com",
			statusCode: 200,
		},
		{
			name:       "server error",
			host:       "https://api.example.com",
			statusCode: 500,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := New(c.host).SetStatusCode(c.statusCode)

			if got.StatusCode != c.statusCode {
				t.Fatalf("got %q; want %q", got.StatusCode, c.statusCode)
			}
		})
	}
}

func TestSetBody(t *testing.T) {
	for _, c := range []struct {
		name       string
		host       string
		statusCode int
		html       string
		want       string
	}{
		{
			name:       "ok",
			host:       "https://www.example.com",
			statusCode: 200,
			html:       "<html><body>A happy webpage</body></html>",
			want:       "<html><body>A happy webpage</body></html>",
		},
		{
			name:       "server error",
			host:       "https://api.example.com",
			statusCode: 500,
			html:       "<html><body>Server Error</body></html>",
			want:       "<html><body>Server Error</body></html>",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: c.statusCode}
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(c.html)))

			got := New(c.host).SetStatusCode(c.statusCode)
			if err := got.SetBody(resp.Body); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got.Body, c.want) {
				t.Fatalf("got %q; want %q", got.Body, c.want)
			}
		})
	}
}
