// Package memfs offers a filesystem in memory.
//
// It was taken from Steve Francia's afero and extended.
// Most importantly, memMapFs.fos (previously memMapFs.fos)
// now holds the full path of each directory.
// Before that, directory names had to be unique.
// Same for InMemoryFile.memDir.
// The type memDirMap was removed; just use builtin map semantics.
// The entire pathing logic was redone.
// There is no os-dependent filepath anymore;
// everything is unix forward slashed.
//
// Creation happens with New(options...)
//
// fileinfo.Name() now returns the basename of the file;
// just as os.FileInfo.Name() does.
//
// All internal usage of Name() had to be rewritten.
//
// There is one mutex for the central map.
// There is one mutex for each file.
//
// Strangely, Remove() did not unregister with parent. I fixed that.
//
// registerDirs(), registerWithParent() and unRegisterWithParent()
//   all occur *outside* the locking block.
// This is probably because they use Open... which in itself
// requires a lock - leadig to deadlock.
//
// In general, changes to the parent's InMemoryFile.memDir map are *not* synced at all!
// Also creating/deleting a file and removing it from it's parent directory
// should be one atomic transaction.
// I believe the original architect of the filesys left some architectural work todo.
// Would it not be clearer to have one goroutine with a for-select manage all directory changes?
// We would send creates, renames and removes on a channel then...
//
//
package memfs

import (
	"os"

	"github.com/pbberlin/tools/os/fsi"
)

const (
	sep = "/" // No support for windows
)

func init() {

	// forcing our implementations
	// to comply with our interfaces

	f := InMemoryFile{}
	ifa := fsi.File(&f)
	_ = ifa

	fi := InMemoryFileInfo{}
	ifi := os.FileInfo(&fi)
	_ = ifi

	fs := memMapFs{}
	ifs := fsi.FileSystem(&fs)
	_ = ifs

}
