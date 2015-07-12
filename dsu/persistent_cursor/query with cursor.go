// Package persistent_cursor shows
// how a query iterator has a cursor object;
// that can be serialized into memcache
// and reused to continue the query in
// a different request
package persistent_cursor

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"appengine"
	ds "appengine/datastore"
	"appengine/memcache"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"

	gbp "github.com/pbberlin/tools/dsu/ancestored_gb_entries" // guest book persistence
)

func guestViewCursor(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	q := ds.NewQuery(gbp.GbEntryKind)
	q.Order("-Date")

	b1 := new(bytes.Buffer)

	cur_start, err := memcache.Get(c, "greeting_cursor")
	if err == nil {
		str_curs := string(cur_start.Value)
		if len(cur_start.Value) > 0 {
			cursor, err := ds.DecodeCursor(str_curs) //  inverse is string()
			loghttp.E(w, r, err, false)
			if err == nil {
				b1.WriteString("found cursor from memcache -" + stringspb.Ellipsoider(str_curs, 10) + "-<br>\n")
				q = q.Start(cursor)
			}
		}
	}

	iter := q.Run(c)
	var cntr int = 0
	for {
		var g gbp.GbEntryRetr
		cntr++
		if cntr > 2 {
			b1.WriteString("  batch complete -" + string(cntr) + "-<br>\n")
			break

		}

		_, err := iter.Next(&g)
		if err == ds.Done {
			b1.WriteString("scan complete -" + string(cntr) + "-<br>\n")
			break
		}

		if fmt.Sprintf("%T", err) == fmt.Sprintf("%T", new(ds.ErrFieldMismatch)) {
			err = nil // ignore this one - it's caused by our deliberate differences between gbsaveEntry and gbEntrieRetr
		}

		if err != nil {
			b1.WriteString("error fetching next: " + err.Error() + "<br>\n")
			break
		}

		b1.WriteString("  - " + g.String())
	}

	// Get updated cursor and store it for next time.
	if cur_end, err := iter.Cursor(); err == nil {

		str_c_end := cur_end.String() //  inverse is decode()
		val := []byte(str_c_end)

		mi_save := &memcache.Item{
			Key:        "greeting_cursor",
			Value:      val,
			Expiration: 60 * time.Second,
		}

		if err := memcache.Set(c, mi_save); err != nil {
			b1.WriteString("error adding memcache item " + err.Error() + "<br>\n")
		} else {
			b1.WriteString("wrote cursor to memcache -" + stringspb.Ellipsoider(str_c_end, 10) + "-<br>\n")
		}

	} else {
		b1.WriteString("could not retrieve cursor_end " + err.Error() + "<br>\n")
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(b1.Bytes())

	w.Write([]byte("<br>----<br>"))

}

func init() {
	http.HandleFunc("/guest-view-cursor", loghttp.Adapter(guestViewCursor))
}
