// Package aefs builds a fully distributed filesystem layer on top of appengine datastore.
package aefs

import (
	"fmt"
	"path"
)

//
// Todos:

// Standardize path - everywhere
//    cases
// 		"/"					mntX/	     ""
//    	mntX				mntX/	     ""
//    	mntX/				mntX/	     ""
//    	mntX/dir1			mntX/	     dir1
//    	mntX/dir1/			mntX/	     dir1
//    	mntX/dir1/file2		mntX/dir1/	 file2
//
//  We always want split.
//
// Solve Remaining Test failures
//
// Remove imports from aefs
//   // logif

// http-handlers:
// 	"github.com/pbberlin/tools/logif"
// "github.com/pbberlin/tools/net/http/loghttp"
// "github.com/pbberlin/tools/net/http/tplx"
// "github.com/pbberlin/tools/os/fsi/fsc"
// "github.com/pbberlin/tools/util"
//
// low-level-dirs-query.go
// 	"github.com/pbberlin/tools/stringspb"
//
//
// Common Remarks:
// ==============================
// "path" is prefixed with RootDir(),
// This is done by all internal methods.
// We don't allow intermittent "virtual" directories.
// All directories are explicitly created. Otherwise traversal is impossible.
//
// "name" or "filename" can mean either the basename or the full path of the file,
// depending on the actual argument - '/tmp/logs/app1.log' or simply 'app1.log'
// In the latter case, it refers to [current dir]/app1.log.
// Btw: golang os file structures have no internal "current dir",
// they save full path into "name".
// Compare http://stackoverflow.com/questions/2235173/file-name-path-name-base-name-naming-standard-for-pieces-of-a-path
//
// Architecture
// ==============================
// According to http://www.cidrdb.org/cidr2011/Papers/CIDR11_Paper32.pdf
// we must chose the granularity of our entity groups.
//
// We decided on using weakly consistent directory paths.
// The directory structure can stomach massive updates and inserts.
// But its indexing on property 'dir' may be delayed.
//
// Direct directory reads are not affected
// Only traversals might miss newest directories.
//
// Only each *one* directory is an entity group.
// Applications are forced to partition directories,
// if files are changed too frequently.
//
// In summary: The entire filesystem is extremely parallel.
// and heavily writeable. But it's structural changes
// are not instantly visible to everyone.
//
// Again: Traversal - meaning ReadDir is done
// using one global index of the Dir property.
// This index can be queried for equality (direct children),
// or for retrieval of entire subtrees.
//
//
// Todo/Consider:
// Add a "block"-layer under file,
// so that more than 1MB byte files can be written?
// At least throw an error before the file is saved?
//
// Mem Caching for files; not just directories - but beware of cost.
// Instance Caching with broadcasting instances via http request to instances.
//
// Locking the filesys upon RemoveAll and Rename?
//
// Nice to have: Links

func splitIsWhatWeWant() {

	s := "rt/"

	dir1, f1 := path.Split(s)

	dir2 := path.Dir(s)
	f2 := path.Base(s)

	fmt.Printf("%q %q \n", dir1, f1) //  "rt/"   ""
	fmt.Printf("%q %q \n", dir2, f2) //  "rt"    "rt"

}
