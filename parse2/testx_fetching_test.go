// +build fetching
// go test -tags=fetching

package parse2

import "testing"

func Test2(t *testing.T) {

	url := spf("%v/contentexport/feed/schlagzeilen", hosts[0])
	Fetch(url, 8)

}
