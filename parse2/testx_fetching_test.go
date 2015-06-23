// +build fetching
// go test -tags=fetching

package parse2

import "testing"

func Test2(t *testing.T) {

	url := "www.handelsblatt.com/contentexport/feed/schlagzeilen"
	Fetch(url, 5)

}
