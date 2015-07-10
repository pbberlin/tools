// Package aefs builds a fully distributed filesystem layer on top of appengine datastore.
package aefs

//
// Todos:
// CreateFile and WriteFile: Only one save operation
// Mem Caching for directories
// Integrate memfs from Afero
// Standardize the root stuff - make "/" synonymous with the RootDir.
//
//
// "path" is internally automatically prefixed with RootDir()
// Subdirs may exist only virtually.
//
// "name" can mean either the basename or the full path of the file,
// depending on the actual argument - /tmp/logs/app1.log or simply app1.log
// In the latter case, it refers to [current dir]/app1.log.
// Btw: golang os file structures have no internal "current dir",
// they save full path into "name".
// Compare // http://stackoverflow.com/questions/2235173/file-name-path-name-base-name-naming-standard-for-pieces-of-a-path
//
// According to http://www.cidrdb.org/cidr2011/Papers/CIDR11_Paper32.pdf
// we must chose the granularity of our entity groups.
//
// We decided on using fictitious directory paths.
// The directory structure must stomach massive updates and inserts.
// Only *one* directory can be an entity group.
// Applications are forced to partition directories,
// if files are changed too frequently.
// Thus, the entire filesystem is extremely parallel.
//
// How to traverse? How to implement ReadDir()?
// Answer: Use one global index of the Dir property.
//
// If we create all directories of all paths,
// such an index could be queried for equality (direct children).
// Otherwise we only can retrieve entire subtree.
//
// Worst disadvantage: Move operations, esp. in high level directories become expensive.
// Advantage: The directory "tree" can be sparse; only lowest dir must exist.
//
// Todo/To consider:
// Add a "block"-layer under file,
// so that more than 1MB byte files can be writtens?
// At least throw an error before the file is saved?
//
// Mem Caching for files - beware of cost
// Instance Caching with broadcasting instances via http request to instances.
//
// Locking the filesys upon RemoveAll and Rename?
//
// Nice to have: FileLinks
