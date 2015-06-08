package parse2

import "net/http"

func serveSingleRootFile(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

const docRoot = "c:/docroot/"

func init() {
	// http.HandleFunc("/", singlePage)

	// static resources - Mandatory root-based
	serveSingleRootFile("/msg.html", docRoot+"msg.html")
	// static resources - other
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(docRoot))))

	go func() {
		// fmt.Println("listening on 4000")
		http.ListenAndServe("localhost:4000", nil)
	}()

}
