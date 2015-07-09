package aefs

import (
	"fmt"

	"appengine/datastore"
)

// subdirsByPath retrieves a subdirectories of a given directory.
// It is relying on an indexed string property "SKey"
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

	filt := "SKey>="
	if onlyDirectChildren {
		filt = "SKey="
	}

	q := datastore.NewQuery(tdir).
		Filter(filt, path).
		Order("SKey")
	//  Limit(4)

	var children []AeDir
	keys, err := q.GetAll(fs.Ctx(), &children)
	if err != nil {
		fs.Ctx().Errorf("Error getting all children of %v => %v", fs.RootDir(), err)
		return children, err
	}

	if len(children) < 1 {
		return children, fmt.Errorf(
			"Query found no result. The SKey index is only eventual consistent.")
	}

	for k, v := range children {
		v.Fs = fs
		v.Key = keys[k]
	}

	return children, nil

}
