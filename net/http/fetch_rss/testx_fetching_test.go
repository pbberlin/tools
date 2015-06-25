package fetch_rss

import "testing"

func Test2(t *testing.T) {

	Serve()

	url := spf("%v/contentexport/feed/schlagzeilen", hosts[0])
	Fetch(url, 19)

}
