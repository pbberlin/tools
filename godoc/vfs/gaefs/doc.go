// Package gaefs builds a filesystem layer on top of appengine datastore.
//
// Todos:

// Complete the usage of the walker function, similar to filepath.Walk(root string, walkFunc)
// 	Use this "walker" to implement removals

//
// Integrate into Afero

// SaveFile => optionally create non-existing directories
//   or return at least path.Error
//
// Add a "block"-layer under file, so that more than 1MB byte files can be writtens.
//
// Mem Caching for directories
// Mem Caching for files - beware of cost
// Instance Caching with broadcasting instances via http request to instances.
//

// Nice to have: FileLinks

package gaefs
