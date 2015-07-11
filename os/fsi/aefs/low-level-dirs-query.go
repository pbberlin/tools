package aefs

import (
	"strings"

	"github.com/pbberlin/tools/os/fsi"

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
// Upon finding nothing, it therefore returns the
// "warning" fsi.EmptyQueryResult
//
// The func could easily be enhanced chunked scanning.
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
		pathInc := IncrementString(path)
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
		return children, fsi.EmptyQueryResult
	}

	// Very evil: We filter out root node, since it's
	// has the same dir as the level-1 directories.
	keyRoot := datastore.NewKey(fs.Ctx(), tdir, fs.mount, 0, nil)
	idxRoot := -1

	for i := 0; i < len(children); i++ {
		children[i].fSys = fs
		children[i].Key = keys[i]
		if keys[i].Equal(keyRoot) {
			idxRoot = i
		}
	}

	if idxRoot > -1 {
		children = append(children[:idxRoot], children[idxRoot+1:]...)
	}

	return children, nil

}
