package dsfs

import (
	"strings"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"

	"appengine/datastore"
)

// subtreeByPath retrieves a subdirectories of a given directory.
// It is relying on an indexed string property "Dir"
// containing a string representation of the full path.
//
// It might be fast for deep, uncached directory subtree,
// that have been saved in nested manner.
//
// However, it might not find recently added directories.
// Upon finding nothing, it therefore returns the
// "warning" fsi.EmptyQueryResult
//
// The func could easily be enhanced chunked scanning.
//
// It is currently used by ReadDir and by the test package.
// It is public for the test package.
func (fs *dsFileSys) SubtreeByPath(name string, onlyDirectChildren bool) ([]DsDir, error) {

	dir, bname := fs.SplitX(name)
	name = dir + common.Filify(bname)
	if !strings.HasSuffix(name, sep) {
		name += sep
	}

	var q *datastore.Query

	if onlyDirectChildren {
		q = datastore.NewQuery(tdir).
			Filter("Dir=", name).
			Order("Dir")
		//  Limit(4)
	} else {
		pathInc := IncrementString(name)
		q = datastore.NewQuery(tdir).
			Filter("Dir>=", name).
			Filter("Dir<", pathInc).
			Order("Dir")
	}

	// log.Printf("%v", q)

	var children []DsDir
	keys, err := q.GetAll(fs.Ctx(), &children)
	if err != nil {
		fs.Ctx().Errorf("Error getting all children of %v => %v", dir+bname, err)
		return children, err
	}

	if len(children) < 1 {
		return children, fsi.EmptyQueryResult
	}

	// Very evil: We filter out root node, since it's
	// has the same dir as the level-1 directories.
	keyRoot := datastore.NewKey(fs.Ctx(), tdir, fs.RootDir(), 0, nil)
	idxRoot := -1

	for i := 0; i < len(children); i++ {
		children[i].fSys = fs
		children[i].Key = keys[i]
		if keys[i].Equal(keyRoot) {
			// log.Printf("root idx %v", i)
			idxRoot = i
		}
	}

	if idxRoot > -1 {
		children = append(children[:idxRoot], children[idxRoot+1:]...)
	}

	return children, nil

}
