package filesys

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pbberlin/tools/stringspb"

	"appengine"
	ds "appengine/datastore"
)

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

	fs.RootDir = fs.newFso(root, nil, true)

	return fs
}

func (fs FileSys) newFso(name string, parent *ds.Key, isDir bool) FSysObj {

	fo := FSysObj{}
	fo.IsDir = isDir
	fo.Name = name
	fo.Mod = time.Now()
	fo.fs = &fs

	suggKey := ds.NewKey(fs.c, t, name, 0, parent)
	fo.Key = suggKey
	fo.SKey = spf("%v", suggKey) // not effKey.Encode()

	effKey, err := ds.Put(fs.c, suggKey, &fo)

	if err != nil {
		fs.c.Errorf("%v  [%s]", err, fo.Key)
	} else {
		fs.c.Infof("fso ds-saved - key %v [%v]", effKey, effKey.Encode())
	}

	fo.Key = effKey
	fo.SKey = spf("%v", effKey) // not effKey.Encode()

	fo.MemCacheSet()

	return fo
}

func (fs FileSys) GetFso(path string) (FSysObj, error) {

	fo := FSysObj{}

	rootKey := ds.NewKey(fs.c, t, fs.RootDir.Name, 0, nil)
	pathInc := stringspb.IncrementString(path)

	q := ds.NewQuery(t).
		Ancestor(rootKey).
		Filter("SKey >=", path).
		Filter("SKey <", pathInc).
		Order("SKey").
		Limit(1)

	if !appengine.IsDevAppServer() {
		q = ds.NewQuery(t).
			Ancestor(rootKey).
			Filter("SKey >=", path).
			Order("SKey").
			Limit(10)
	}

	var children []FSysObj
	_, err := q.GetAll(fs.c, &children)
	if err != nil {
		fs.c.Errorf("Error getting all children of %v => %v", fs.RootDir.Name, err)
		return fo, err
	} else {
		fs.c.Infof(" got %v fso's between %v and %v", len(children), path, pathInc)
	}

	for k, v := range children {
		fs.c.Infof("%-4v => %v", k, v.SKey)
	}

	return fo, nil
}
