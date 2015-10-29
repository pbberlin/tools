// Package fulltext is just an
// exercise in Google fulltext API.
//
package fulltext

import (
	"net/http"
	"time"

	"appengine"
	"appengine/search"

	"fmt"

	"github.com/pbberlin/tools/net/http/loghttp"
)

type User struct {
	Customer  string
	Comment   search.HTML
	Visits    float64
	LastVisit time.Time
	Birthday  time.Time
}

func searchPut(w http.ResponseWriter, r *http.Request) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = b

	id := "PA6-5001"
	user := &User{
		Customer:  "Carl Corral",
		Comment:   "I am <em>riled up</em> text",
		Visits:    1,
		LastVisit: time.Now(),
		Birthday:  time.Date(1968, time.May, 19, 0, 0, 0, 0, time.UTC),
	}

	c := appengine.NewContext(r)

	index, err := search.Open("users")
	lg(err)

	ret_id, err := index.Put(c, id, user)
	lg(err)

	fmt.Fprint(w, "OK, saved "+ret_id+"\n\n")

	var u2 User
	err = index.Get(c, ret_id, &u2)
	lg(err)
	fmt.Fprint(w, "Retrieved document: ", u2)

}

func searchRetrieve(w http.ResponseWriter, r *http.Request) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = b

	c := appengine.NewContext(r)
	index, err := search.Open("users")
	lg(err)

	for t := index.Search(c, "Comment:riled", nil); ; {
		var res User
		id, err := t.Next(&res)
		fmt.Fprintf(w, "\n-- ")
		if err == search.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "Search error: %v\n", err)
			break
		}
		fmt.Fprintf(w, "%s -> %#v\n", id, res)
	}
}

func init() {

	http.HandleFunc("/fulltext-search/put", searchPut)
	http.HandleFunc("/fulltext-search/get", searchRetrieve)

}
