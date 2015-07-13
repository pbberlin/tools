// Package aefs builds a fully distributed
// filesystem layer on top of appengine datastore.
package aefs

//
//
// Common Remarks:
// ==============================
// To remain compatible with osfs,
// we support "." for current working dir.
// We could introduce ChDir(), but so far
// current working dir is always RootDir().
//
//
// Terminology:
// ==============================
// "name" or "filename" can mean either the basename or the full path of the file,
// depending on the actual argument - '/tmp/logs/app1.log' or simply 'app1.log'
// In the latter case, it refers to [current dir]/app1.log => [root dir]/app1.log
// Exception: os.FileInfo.Name() contains only the base name.
// Compare http://stackoverflow.com/questions/2235173/file-name-path-name-base-name-naming-standard-for-pieces-of-a-path
//
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
// Direct directory reads are not always consistent.
// Only subdirectory queries are affected.
// Such traversals might miss newest directories.
// Such traversals might report directories already deleted.
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
// Implement rename?
// Locking the filesys during RemoveAll and Rename?
//
// Nice to have: Links
