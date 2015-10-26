package memcachepb

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"

	aeOrig "appengine"
)

func InitHandlers() {
	http.HandleFunc("/memcache/flush", flushMemcache)
}

func flushMemcache(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c2 := aeOrig.NewContext(r)
	errMc := memcache.Flush(c)
	if errMc != nil {
		c2.Errorf("Error flushing memache: %v", errMc)
		return
	}
	w.Write([]byte("ok"))
}
