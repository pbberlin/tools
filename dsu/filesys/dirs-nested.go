package filesys

import (
	"time"

	"github.com/pbberlin/tools/stringspb"

	"appengine"
	ds "appengine/datastore"
)

func (fs FileSys) newFsoByParentKey(name string, parent *ds.Key, isDir bool) FSysObj {

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

func (fs FileSys) GetFsoByQuery(path string) (FSysObj, error) {

	fo := FSysObj{}

	rootKey := ds.NewKey(fs.c, t, fs.RootDir.Name, 0, nil)
	pathInc := stringspb.IncrementString(path)

	q := ds.NewQuery(t).
		Ancestor(rootKey).
		Filter("SKey>=", path).
		Filter("SKey<", pathInc).
		Order("SKey").
		Limit(4)

	if appengine.IsDevAppServer() {
		q = ds.NewQuery(t).
			Ancestor(rootKey).
			Filter("SKey>=", path).
			Filter("SKey<", pathInc).
			Order("SKey").
			Limit(4)
	}

	var children []FSysObj
	_, err := q.GetAll(fs.c, &children)
	if err != nil {
		fs.c.Errorf("Error getting all children of %v => %v", fs.RootDir.Name, err)
		return fo, err
	} else {
		fs.c.Infof(" got %v fso's between %v --- %v", len(children), path, pathInc)
	}

	for k, v := range children {
		fs.c.Infof("%-4v => %v", k, v.SKey)
	}
	for _, v := range children {
		fo = v
		break
	}

	return fo, nil
}
