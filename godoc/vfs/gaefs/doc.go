// Package gaefs builds a filesystem layer on top of appengine datastore.
//
// Todos:
// Rename members of AeFile
// Implement AeFile.ReadDir

//
// Integrate into Afero

// Add a Walk function, similar to filepath.Walk(root string, walkFunc)
// 	Use this "walker" to implement removals

// SaveFile => optionally create non-existing directories
//   or return at least path.Error
//
// Add a "block"-layer under file, so that more than 1MB byte files can be writtens.
//
// Mem Caching for directories
// Mem Caching for files - beware of cost
//
// Instance Caching with broadcasting instances via http request to instances.
//
// ReadDir and GetFiles should sort resulting directories/files by name.

// Nice to have: FileLinks

package gaefs
