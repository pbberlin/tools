package aefs

import (
	"strings"

	"github.com/pbberlin/tools/os/fsi/fsc"
	"github.com/pbberlin/tools/stringspb"

	"appengine/datastore"
)

// subdirsByPath retrieves a subdirectories of a given directory.
// It is relying on an indexed string property "Dir"
// containing a string representation of the full path.
//
// It might be fast for deep, uncached directory subtrees,
// that have been saved in nested manner.
//
// However, it might not find recently added directories.
//
// If a range scan over a huge directory tree is neccessary,
// the func could easily be enhanced chunked scanning.
func (fs *AeFileSys) subdirsByPath(path string, onlyDirectChildren bool) ([]AeDir, error) {

	path = cleanseLeadingSlash(path)
	if !strings.HasPrefix(path, fs.RootName()) {
		path = fs.RootDir() + path
	}
	if !strings.HasSuffix(path, sep) {
		path += sep
	}

	var q *datastore.Query

	if onlyDirectChildren {
		q = datastore.NewQuery(tdir).
			Filter("Dir=", path).
			Order("Dir")
		//  Limit(4)
	} else {
		pathInc := stringspb.IncrementString(path)
		q = datastore.NewQuery(tdir).
			Filter("Dir>=", path).
			Filter("Dir<", pathInc).
			Order("Dir")
	}

	// logif.Pf("%v", q)

	var children []AeDir
	keys, err := q.GetAll(fs.Ctx(), &children)
	if err != nil {
		fs.Ctx().Errorf("Error getting all children of %v => %v", fs.RootDir(), err)
		return children, err
	}

	if len(children) < 1 {
		return children, fsc.EmptyIndexQueryResult
	}

	// Very evil: We filter out root node, since it's
	// has the same dir as the level-1 directories.
	keyRoot := datastore.NewKey(fs.Ctx(), tdir, fs.mount, 0, nil)
	// keySelf := datastore.NewKey(fs.Ctx(), tdir, path, 0, nil)
	idxRoot := -1
	for k, v := range children {
		v.fSys = fs
		v.Key = keys[k]
		// if keys[k].Equal(keyRoot) || keys[k].Equal(keySelf) {
		if keys[k].Equal(keyRoot) {
			idxRoot = k
		}
	}

	if idxRoot > -1 {
		// logif.Pf("self excluded")
		children = append(children[:idxRoot], children[idxRoot+1:]...)
	}

	return children, nil

}
