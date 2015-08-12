package fetch_rss

import (
	"testing"
	"time"
)

func Test2(t *testing.T) {

	Serve()

	for _, config := range hosts {
		Fetch(config, "/politik/international/aa/bb", 12)
	}

	time.Sleep(220 * time.Second)
}
