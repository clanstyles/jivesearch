package queue

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func TestAddLink(t *testing.T) {
	for _, c := range []struct {
		name string
		link string
	}{
		{
			"first", "http://www.example.com",
		},
		{
			"second", "https://www.somelink.com/and/a/path/?for=fun",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := &Redis{}
			conn := redigomock.NewConn()
			conn.Command("SADD", links, c.link).Expect("OK")

			r.RedisPool = &redis.Pool{
				Dial: func() (redis.Conn, error) {
					return conn, nil
				},
			}
			defer r.RedisPool.Close()

			if err := r.AddLink(c.link); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestQueueLink(t *testing.T) {
	for _, c := range []struct {
		name string
		link string
	}{
		{
			"first", "http://www.example.com",
		},
		{
			"second", "https://www.somelink.com/and/a/path/?for=fun",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ttl := 10 * time.Minute

			r := &Redis{}
			conn := redigomock.NewConn()
			conn.Command("SPOP", links).Expect(c.link)
			conn.Command("SET", r.prefixKey(queuePrefix+c.link), "", "EX", int(ttl/time.Second), "NX").Expect("OK")

			r.RedisPool = &redis.Pool{
				Dial: func() (redis.Conn, error) {
					return conn, nil
				},
			}
			defer r.RedisPool.Close()

			got, err := r.QueueLink(ttl)
			if err != nil {
				t.Fatal(err)
			}

			if got != c.link {
				t.Fatalf("got %v; want: %v", got, c.link)
			}
		})
	}
}

func TestReserveHost(t *testing.T) {
	// this does NOT check if the key actually expires
	for _, c := range []struct {
		name  string
		host  string
		delay time.Duration
	}{
		{
			"10m reservation", "http://www.example.com", 10 * time.Minute,
		},
		{
			"30s reservation", "https://api.somewebsite.org", 30 * time.Second,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := &Redis{}
			conn := redigomock.NewConn()
			k := r.prefixKey(hostPrefix + c.host)

			conn.Command("SET", k, "", "EX", int(c.delay)/1e9, "NX").Expect("OK")

			r.RedisPool = &redis.Pool{
				Dial: func() (redis.Conn, error) {
					return conn, nil
				},
			}
			defer r.RedisPool.Close()

			if err := r.ReserveHost(c.host, c.delay); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDelayHost(t *testing.T) {
	for _, c := range []struct {
		name  string
		host  string
		delay time.Duration
	}{
		{"0s delay", "http://www.someexample.com", 0 * time.Second},
		{"10s delay", "http://www.example.com", 10 * time.Second},
		{"30s delay", "https://api.somewebsite.org", 30 * time.Second},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := &Redis{}
			conn := redigomock.NewConn()
			k := r.prefixKey(hostPrefix + c.host)
			conn.Command("DEL", k).Expect(1)
			conn.Command("SET", k, "", "EX", int(c.delay)/1e9, "XX").Expect("OK")

			r.RedisPool = &redis.Pool{
				Dial: func() (redis.Conn, error) {
					return conn, nil
				},
			}
			defer r.RedisPool.Close()

			if err := r.DelayHost(c.host, c.delay); err != nil {
				t.Fatal(err)
			}
		})
	}
}
