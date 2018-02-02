package log

import (
	"log"
	"testing"
)

func TestSetDefaults(t *testing.T) {
	for _, c := range []struct {
		name   string
		logger *log.Logger
		want   string
	}{
		{
			"Info", Info, "INFO ",
		},
		{
			"Debug", Debug, "DEBUG ",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := c.logger.Prefix(); got != c.want {
				t.Fatalf("got %q; want %q", got, c.want)
			}
		})
	}

}
