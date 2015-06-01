package dsu_ancestored_urls

import (
	"appengine"
	ds "appengine/datastore"
	"bytes"
	"fmt"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_err"
	"io"
	"net/http"
)

type LastURL struct {
	Value string
}

const keyNoAncInt = 0
const kindUrl = "classUrl"

//               5623589659213824
const keyUrl = "constant__KeyUrl"

const kindUrlParent = "classUrlParent"
const keyUrlParent = "keyUrlParent"

func ancKey(c appengine.Context) *ds.Key {
	return ds.NewKey(c, kindUrlParent, keyUrlParent, 0, nil)
}

// saving some data by kind and key
//   without ancestor key
func saveURLNoAnc(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	k := ds.NewKey(c, kindUrl, keyUrl, 0, nil)
	e := new(LastURL)
	err := ds.Get(c, k, e)
	if err == ds.ErrNoSuchEntity {
		util_err.Err_log(err)
	} else {
		util_err.Err_http(w, r, err, false)
	}

	old := e.Value
	e.Value = r.URL.Path + r.URL.RawQuery

	_, err = ds.Put(c, k, e)
	util_err.Err_http(w, r, err, false)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("old=" + old + "\n"))
	w.Write([]byte("new=" + e.Value + "\n"))

}

func saveURLWithAncestor(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	k := ds.NewKey(c, kindUrl, "", 0, ancKey(c))

	s := util.TimeMarker()
	ls := len(s)
	lc := len("09-26 17:29:25")
	lastURL_fictitious_1 := LastURL{"with_anc " + s[ls-lc:ls-3]}
	_, err := ds.Put(c, k, &lastURL_fictitious_1)
	util_err.Err_http(w, r, err, false)

}

// get all URLs
//  not just by ancestor
//  results might be delayed
func listURLNoAnc(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	c := appengine.NewContext(r)

	b1 := new(bytes.Buffer)
	q := ds.NewQuery(kindUrl).
		Filter("Value>", `/save`).
		Order("-Value")

	i := -1
	for t := q.Run(c); ; {
		i++
		var lu LastURL
		key, err := t.Next(&lu)
		if err == ds.Done {
			b1.WriteString("\nds.Done\n")
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(b1, "q loop ", i, "\n")
		fmt.Fprintf(b1, "\tKey %64s  \n\tVal %64s\n", key, lu.Value)
	}

	w.Write(b1.Bytes())

}

// get all ancestor urls
//  ordering possible
//  distinction to above: *always* consistent
func listURLWithAncestors(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	c := appengine.NewContext(r)

	q := ds.NewQuery(kindUrl).
		Ancestor(ancKey(c)).
		Order("-Value")
	var vURLs []LastURL
	keys, err := q.GetAll(c, &vURLs)
	util_err.Err_http(w, r, err, false)

	for i, v := range vURLs {
		io.WriteString(w, fmt.Sprint("q loop ", i, "\n"))
		io.WriteString(w, fmt.Sprintf("\tKey %64s  \n\tVal %64s\n", keys[i], v.Value))
	}
	io.WriteString(w, "\nds.Done\n")

}
