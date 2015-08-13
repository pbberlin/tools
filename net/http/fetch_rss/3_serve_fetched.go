package fetch_rss

import (
	"log"
	"net/http"
	"path"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/httpfs"
)

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
 from handelsblatt or economist.
</p>
<p>
	We want to find 
	Longest Common Subsequences - LCS
	in a DOM tree - but with tolerance for 20 characters.
</p>
`)
}

func Serve(fs fsi.FileSystem) (baseUrl string, topDirs []string) {

	fs.WriteFile(path.Join(docRoot, "msg.html"), msg, 0644)

	httpFSys := &httpfs.HttpFs{SourceFs: fs}
	mux := http.NewServeMux()
	fileserver1 := http.FileServer(httpFSys.Dir(docRoot))
	mux.Handle("/", fileserver1)
	mux.Handle("/static2/", http.StripPrefix("/static2/", fileserver1)) // same

	go func() {
		log.Fatal(http.ListenAndServe("localhost:4000", mux))
	}()

	topDirs = make([]string, 0, len(hosts))
	for k, _ := range hosts {
		topDirs = append(topDirs, k)
	}
	return
}
