// Package queue manages the queue for a distributed crawler
package queue

import (
	"errors"
	"time"
)

// Queuer is handles links and our crawling queue
type Queuer interface {
	CountLinks() (int64, error)
	AddLink(lnk string) error
	QueueLink(ttl time.Duration) (string, error)
	ReserveHost(host string, ttl time.Duration) error
	DelayHost(host string, ttl time.Duration) error
}

// ErrNotQueued indicates a link was not queued
var ErrNotQueued = errors.New("link already queued")

// ErrAlreadyReserved indicates another worker reserved the host
var ErrAlreadyReserved = errors.New("host already reserved")
var errNotDelayed = errors.New("host not delayed")
