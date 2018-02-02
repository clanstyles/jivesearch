// Package robots handles caching robots.txt files
package robots

import (
	"io"
	"io/ioutil"
	"time"
)

const dateFormat = "200601021504"

// Robots is a single robots.txt response
type Robots struct {
	SchemeHost string `json:"-"`
	StatusCode int    `json:"status"`
	Body       string `json:"body"`
	Expires    string `json:"expires"`
	Cached     bool   `json:"-"`
}

// Cacher handles the caching backend for robots.txt files
type Cacher interface {
	IndexExists() (bool, error)
	Setup() error
	Put(*Robots)
	Get(host string) (*Robots, error)
}

// this makes testing our SetExpires and Expired methods easier
var now = func() time.Time {
	return time.Now().UTC()
}

// New creates a new Robot & sets the ID to the host
// Robots are handled on a per-host basis
// https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
func New(sh string) *Robots {
	return &Robots{SchemeHost: sh}
}

// SetExpires marks the date the robots.txt file should expire from cache
func (r *Robots) SetExpires() *Robots {
	// Googlebot will cache the robots.txt file for up to 1 day
	// In the case of 5xx errors on fetching robots.txt we cache it
	// for less time.
	ttl := 24 * time.Hour
	if r.StatusCode >= 500 && r.StatusCode < 600 {
		ttl = 1 * time.Hour
	}
	expires := now().Add(ttl)
	r.Expires = expires.Format(dateFormat)
	return r
}

// Expired tells us if a robots.txt file needs to be refetched.
func (r *Robots) Expired() (bool, error) {
	expires, err := time.Parse(dateFormat, r.Expires)
	if err != nil {
		return true, err
	}
	return expires.Sub(now()) < 0, nil
}

// SetStatusCode sets the HTTP response code
func (r *Robots) SetStatusCode(code int) *Robots {
	r.StatusCode = code
	return r
}

// SetBody saves the body of a robots.txt file
// Since Group has unexported fields we must save the entire body, not just the group
func (r *Robots) SetBody(b io.Reader) error {
	htmlData, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}

	r.Body = string(htmlData)

	return err
}
