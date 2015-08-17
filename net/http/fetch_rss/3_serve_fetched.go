package fetch_rss

import "net/http"

// unused
func serveSingleRootFile(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename) // filename refers to local path; unusable for fsi
	})
}

var msg []byte

func init() {
	msg = []byte(`<p>This is an embedded static http server.</p>
<p>
It serves previously downloaded pages<br>
 i.e. from handelsblatt or economist.
</p>
<p>
	We want to find 
	Longest Common Subsequences - LCS
	in a DOM tree - but with tolerance for 20 characters.
</p>
`)
}
