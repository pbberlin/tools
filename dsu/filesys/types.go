package filesys

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"appengine"

	ds "appengine/datastore"
)

// Filesystem
type FileSys struct {
	w http.ResponseWriter `datastore:"-"`
	r *http.Request       `datastore:"-"`
	c appengine.Context   `datastore:"-"`

	RootDir FSysObj
}

// Filesystem Object - a directory or a file
type FSysObj struct {
	fs    *FileSys
	Name  string
	IsDir bool
	Mod   time.Time

	Key  *ds.Key
	SKey string // from *ds.Key.Encode()
}

var t string

func init() {
	fo := FSysObj{}
	t = fmt.Sprintf("%T", fo) // "kind"
	t = "fso"
}

func NewFileSys(w http.ResponseWriter, r *http.Request, root string) FileSys {
	fs := FileSys{}
	fs.w = w
	fs.r = r
	fs.c = appengine.NewContext(r)
	if strings.Contains(root, "/") {
		panic("root can't have slash in it")
	}
	fs.RootDir = fs.newFsoByParentKey(root, nil, true)
	return fs
}
