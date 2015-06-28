package nested

// LowLevelArchitecture

import (
	"net/http"
	"time"

	"github.com/pbberlin/tools/dsu/filesys"
	"github.com/pbberlin/tools/stringspb"

	"appengine"
	ds "appengine/datastore"
)

type Nested struct {
	filesys.FileSys // inheriting methods of filesys
}

func NewNestedFileSys(w http.ResponseWriter, r *http.Request, mount string) Nested {
	fs := filesys.NewFileSys(w, r, mount)
	n := Nested{}
	n.FileSys = fs
	return n
}

func (fs Nested) saveDirByPath(path string) filesys.Directory {
	fo := filesys.Directory{}
	return fo
}

func (fs Nested) saveDirectoryUnderParent(name string, parent *ds.Key) filesys.Directory {

	fo := filesys.Directory{}
	fo.IsDir = true
	fo.Name = name
	fo.Mod = time.Now()

	suggKey := ds.NewKey(fs.Ctx(), t, name, 0, parent)
	fo.Key = suggKey
	fo.SKey = spf("%v", suggKey) // not effKey.Encode()

	effKey, err := ds.Put(fs.Ctx(), suggKey, &fo)

	if err != nil {
		fs.Ctx().Errorf("%v  [%s]", err, fo.Key)
	} else {
		fs.Ctx().Infof("fso ds-saved - key %v [%v]", effKey, effKey.Encode())
	}

	fo.Key = effKey
	fo.SKey = spf("%v", effKey) // not effKey.Encode()

	fo.MemCacheSet()

	return fo
}

func (fs Nested) dirByPath(path string) (filesys.Directory, error) {

	fo := filesys.Directory{}

	rootKey := ds.NewKey(fs.Ctx(), t, fs.RootDir.Name, 0, nil)
	pathInc := stringspb.IncrementString(path)

	q := ds.NewQuery(t).
		Ancestor(rootKey).
		Filter("SKey>=", path).
		Filter("SKey<", pathInc).
		Order("SKey").
		Limit(4)

	if appengine.IsDevAppServer() {
		// query variation
	}

	var children []filesys.Directory
	_, err := q.GetAll(fs.Ctx(), &children)
	if err != nil {
		fs.Ctx().Errorf("Error getting all children of %v => %v", fs.RootDir.Name, err)
		return fo, err
	} else {
		fs.Ctx().Infof(" got %v fso's between %v --- %v", len(children), path, pathInc)
	}

	for k, v := range children {
		fs.Ctx().Infof("%-4v => %v", k, v.SKey)
	}
	for _, v := range children {
		fo = v
		break
	}

	return fo, nil
}
