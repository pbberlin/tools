package ancestored_gb_entries

import (
	"net/http"
	"time"

	"appengine"
	ds "appengine/datastore"
	"appengine/user"

	"bytes"
	"fmt"
	"strings"
	//"reflect"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"
)

/*
	default_guestbook[T:Guestbook]
		\
		 - Entry_1
		 - Entry_2



*/

const (
	kind_guestbk string = "class_parent_gb"    // "classname" of the parent
	key_guestbk  string = "instance_parent_gb" // string Key

	DSKindGBEntry string = "gbEntry" // "classname" of a guestbookk entry

)

type GbsaveEntry struct {
	Author            string
	Content           string `datastore:"Content,noindex" json:"content"`
	Date              time.Time
	unsaved           string `datastore:"-"`
	TypeShiftingField int
	Comment1          string
}

// just for experiment
// we retrieve into a different structure then we saved into
//
type GbEntryRetr struct {
	Author            string
	Content           string
	Date              time.Time
	Field2            string
	TypeShiftingField uint8
	Comment1          string
}

// returns entity group - or parent - key
//   to store and retrieve all guestbook entries.
//   the content of this parent is nil
//   it only servers as umbrella for the entries
func key_entity_group_key_parent(c appengine.Context) (r *ds.Key) {
	// key_guestbk could be varied for multiple guestbooks.
	// Either key_guestbk XOR key_int_guestbk must be zero
	var key_int_guestbk int64 = 0
	var key_parent *ds.Key = nil
	r = ds.NewKey(c, kind_guestbk, key_guestbk, key_int_guestbk, key_parent)
	return
}

// implementing the "stringer" interface
func (g GbEntryRetr) String() string {

	b1 := new(bytes.Buffer)

	b1.WriteString(g.Author + "<br>\n")
	b1.WriteString(g.Content + "<br>\n")
	f2 := "2006-01-02 (Jan 02)"
	s2 := g.Date.Format(f2)
	b1.WriteString(s2 + "<br>\n")

	return b1.String()
}

func SaveEntry(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	contnt, ok := m["content"].(string)
	loghttp.E(w, r, ok, false, "need a key 'content' with a string")

	c := appengine.NewContext(r)

	g := GbsaveEntry{
		Content:           contnt,
		Date:              time.Now(),
		TypeShiftingField: 2,
	}
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	g.Comment1 = "comment"

	/* We set the same parent key on every GB entry entity
	   to ensure each GB entry is in the same entity group.

	   Queries across the single entity group will be consistent.
	   However, we should limit write rate to single entity group ~1/sec.

		NewIncompleteKey(appengine.Context, kind string , parent *Key)
			has neither a string key, nor an integer key
		 	only a "kind" (classname) and a parent
		Upon usage the datastore generates an integer key
	*/

	incomplete := ds.NewIncompleteKey(c, DSKindGBEntry, key_entity_group_key_parent(c))
	concreteNewKey, err := ds.Put(c, incomplete, &g)
	loghttp.E(w, r, err, false)
	_ = concreteNewKey // we query entries via key_entity_group_key_parent - via parent

}

func ListEntries(w http.ResponseWriter,
	r *http.Request) (gbEntries []GbEntryRetr, report string) {

	c := appengine.NewContext(r)
	/* High Replication Datastore:
	Ancestor queries are strongly consistent.
	Queries spanning MULTIPLE entity groups are EVENTUALLY consistent.
	If .Ancestor was omitted from this query, there would be slight chance
	that recent GB entry would not show up in a query.
	*/
	q := ds.NewQuery(DSKindGBEntry).Ancestor(key_entity_group_key_parent(c)).Order("-Date").Limit(10)
	gbEntries = make([]GbEntryRetr, 0, 10)
	keys, err := q.GetAll(c, &gbEntries)

	if fmt.Sprintf("%T", err) == fmt.Sprintf("%T", new(ds.ErrFieldMismatch)) {
		//s := fmt.Sprintf("%v %T  vs %v %T <br>\n",err,err,ds.ErrFieldMismatch{},ds.ErrFieldMismatch{})
		loghttp.E(w, r, err, true)
		err = nil // ignore this one - it's caused by our deliberate differences between gbsaveEntry and gbEntrieRetr
	}
	loghttp.E(w, r, err, false)

	// for investigative purposes,
	// we
	var b1 bytes.Buffer
	var sw string
	var descrip []string = []string{"class", "path", "key_int_guestbk"}
	for i0, v0 := range keys {
		sKey := fmt.Sprintf("%v", v0)
		v1 := strings.Split(sKey, ",")
		sw = fmt.Sprintf("key %v", i0)
		b1.WriteString(sw)
		for i2, v2 := range v1 {
			d := descrip[i2]
			sw = fmt.Sprintf(" \t %v:  %q ", d, v2)
			b1.WriteString(sw)
		}
		b1.WriteString("\n")
	}
	report = b1.String()

	for _, gbe := range gbEntries {
		s := gbe.Comment1
		if len(s) > 0 {
			if pos := strings.Index(s, "0300"); pos > 1 {
				i1 := util.Max(pos-4, 0)
				i2 := util.Min(pos+24, len(s))
				s1 := s[i1:i2]
				s1 = strings.Replace(s1, "3", "E", -1)
				report = fmt.Sprintf("%v -%v", report, s1)
			}
		}
	}

	return
}
