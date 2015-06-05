package transposablematrix

import "net/http"

func serveSingleRootFile(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func init() {
	// http.HandleFunc("/", singlePage)

	// static resources - Mandatory root-based
	serveSingleRootFile("/msg.html", "c:/docroot/msg.html")
	// static resources - other
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("c:/docroot"))))

	go func() {
		// fmt.Println("listening on 4000")
		http.ListenAndServe("localhost:4000", nil)
	}()

}
