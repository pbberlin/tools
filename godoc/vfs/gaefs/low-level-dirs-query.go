package gaefs

import (
	"fmt"

	"github.com/pbberlin/tools/stringspb"

	"appengine"
	"appengine/datastore"
)

// GetDirByPathQuery retrieves a directory *at once*.
// It is relying on an indexed string property "SKey"
// containing a string representation of the full pth.
//
// It might be fast for deep, uncached directory trees,
// that have been saved in nested manner.
//
// However, it might not find recently added directories.
//
// It works on nested *and* rooted storage schemes;
// only the path encoding would be different.
//
// If a range scan over a huge directory tree is neccessary,
// the func could easily be enhanced for range scans.
func (fs *AeFileSys) GetDirByPathQuery(path string) (AeDir, error) {

	fo := AeDir{}
	fo.Fs = fs

	rootKey := datastore.NewKey(fs.Ctx(), tdir, fs.RootDir.BName, 0, nil)
	pathInc := stringspb.IncrementString(path)

	q := datastore.NewQuery(tdir).
		Ancestor(rootKey).
		Filter("SKey>=", path).
		Filter("SKey<", pathInc).
		Order("SKey").
		Limit(4)

	if appengine.IsDevAppServer() {
		// query variation
	}

	var children []AeDir
	keys, err := q.GetAll(fs.Ctx(), &children)
	if err != nil {
		fs.Ctx().Errorf("Error getting all children of %v => %v", fs.RootDir.Name, err)
		return fo, err
	}

	// fs.Ctx().Infof(" got %v fso's between %v --- %v", len(children), path, pathInc)
	// for k, v := range children {
	// 	fs.Ctx().Infof("%-4v => %v", k, v.SKey)
	// }

	if len(children) < 1 {
		return fo, fmt.Errorf(
			"Query found no result. The SKey index is only eventual consistent.")
	}

	children[0].Fs = fs
	children[0].Key = keys[0]

	return children[0], nil

}
