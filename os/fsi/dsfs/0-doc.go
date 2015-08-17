// Package dsfs builds a fully distributed
// filesystem layer on top of google datastore.
package dsfs

//
// Architecture
// ==============================
// According to http://www.cidrdb.org/cidr2011/Papers/CIDR11_Paper32.pdf
// we must choose the granularity of our entity groups.
//
// We decided on using weakly consistent directory paths.
// Thus, the directory structure can stomach massive updates and inserts.
// But its indexing on property 'dir' may be delayed.
//
// Only each *one* directory is an entity group.
// Applications are forced to partition directories,
// if files *per directory* are changed too frequently.
//
// We cannot allow intermittent "virtual" directories.
// All directories must be explicitly created. Otherwise traversal is impossible.
//
// Direct directory reads are always consistent.
// Only subdirectory queries might be slightly stale.
// Traversals might miss newest directories.
// Traversals might report directories already deleted.
//
//
// In summary: The entire filesystem is extremely parallel,
// and heavily writeable. But it's structural changes
// are not instantly visible to everyone.
//
// Again: Traversal - meaning ReadDir() - is done
// using one global index of the Dir property.
// This index can be queried for equality (direct children),
// or for retrieval of entire subtrees.
//
//
// Todo/Consider:
// Add a "block"-layer under file layer,
// so that more than 1MB files can be written?
// At least throw an error before the file is saved?
//
// Mem Caching for files; not just directories - but beware of cost.
//
// Combine with memfs?
// Usage of instance caching with broadcasting instances
// via http request to instances?
//
// Rename is not implemented.
// Rename can be an expensive operation.
//
// RemoveAll and Rename might have to lock
// parts of the filesystem.
// See memfs, for how this could be done.
//
// Nice to have: Links
