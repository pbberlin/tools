package aefs

import (
	"fmt"
	"strings"

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

	var q *datastore.Query

	if onlyDirectChildren {
		if !strings.HasSuffix(path, sep) {
			path += sep
		}
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
		return children, fmt.Errorf(
			"Query found no result. The Dir index is only eventual consistent.")
	}

	for k, v := range children {
		v.fSys = fs
		v.Key = keys[k]
	}

	return children, nil

}
