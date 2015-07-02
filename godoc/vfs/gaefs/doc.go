// Package gaefs builds a filesystem layer on top of appengine datastore.
//
// Todos:
// SaveFile => optionally create non-existing directories
// Completely implement the x/tools/godoc/vfs-Interface
// Add a Walk function, similar to filepath.Walk(root string, walkFunc)
// Add a "block"-layer under file, so that more than 1MB byte files can be used.

// Mem Caching for directories
// Mem Caching for files - beware of cost
//

// Instance Caching with broadcasting instances via http request to instances.
//

// ReadDir, GetFiles should sort resulting directories/files by name.

// Nice to have: Links

package gaefs
