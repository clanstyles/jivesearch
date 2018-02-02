package queue

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	prefix      = "jivesearch:"
	hostPrefix  = "h:"
	queuePrefix = "q:"
	links       = prefix + "links"
)

// Redis implements the Queuer interface
type Redis struct {
	RedisPool *redis.Pool
}

func (r *Redis) prefixKey(key string) string {
	return prefix + key // jivesearch:http://www.example.com
}

// grab connection from pool and do the redis cmd
func (r *Redis) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := r.RedisPool.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// CountLinks counts the number of links in our queue
func (r *Redis) CountLinks() (int64, error) {
	cnt, err := redis.Int64(r.do("SCARD", links))
	return cnt, err
}

// AddLink adds a link to our redis set
func (r *Redis) AddLink(lnk string) error {
	_, err := r.do("SADD", links, lnk)
	return err
}

// QueueLink pops a link from our set
func (r *Redis) QueueLink(ttl time.Duration) (string, error) {
	lnk, err := redis.String(r.do("SPOP", links))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return "", err
	}

	k := r.prefixKey(queuePrefix + lnk)
	set, err := r.do("SET", k, "", "EX", seconds(ttl), "NX")
	if set != "OK" && err == nil { // means it is already queued
		lnk = ""
	}

	return lnk, err
}

// ReserveHost reserves a host for crawling
func (r *Redis) ReserveHost(host string, ttl time.Duration) error {
	k := r.prefixKey(hostPrefix + host)
	s := seconds(ttl)

	set, err := r.do("SET", k, "", "EX", s, "NX")
	if err != nil {
		return err
	}

	if set != "OK" {
		err = ErrAlreadyReserved
	}
	return err
}

// DelayHost is like ReserveHost but makes sure the key is already set.
// If key doesn't exist then something is wrong (eg reservation < crawler's timeout)
func (r *Redis) DelayHost(host string, ttl time.Duration) error {
	k := r.prefixKey(hostPrefix + host)
	s := seconds(ttl)

	if s == 0 {
		rsp, err := r.do("DEL", k)
		if err != nil {
			return err
		}
		if rsp == 0 {
			err = errNotDelayed // indicates key doesn't exist
		}

		return err
	}

	rsp, err := r.do("SET", k, "", "EX", s, "XX")
	if err != nil {
		return err
	}

	if rsp != "OK" {
		err = errNotDelayed // indicates key doesn't exist
	}

	return err
}

func seconds(ttl time.Duration) int {
	return int(ttl / time.Second)
}
