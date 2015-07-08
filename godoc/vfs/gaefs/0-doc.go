// Package gaefs builds a filesystem layer on top of appengine datastore.
//
// "name" can mean either the basename or the full path of the file,
// depending on the actual argument - /tmp/logs/app1.log or simply app1.log
// In the latter case, it refers to [current dir]/app1.log.
// Btw: golang os file structures have no internal "current dir",
// they save full path into "name".
// Compare // http://stackoverflow.com/questions/2235173/file-name-path-name-base-name-naming-standard-for-pieces-of-a-path
//
// Todos:
//
// According to http://www.cidrdb.org/cidr2011/Papers/CIDR11_Paper32.pdf
// we must chose the granularity of our entity groups.
// So far, everything is erroneously inside one entity group.
//
// Despite "nesting" of keys, everything is saved under one root.
// Queries may return huge subtrees.
// Direct children need to be filtered out at extra cost.
// Writes are limited to single datastore "tablet" performance.
//
// We have to rewrite everything using fictitious directory paths
// The directory structure must stomach massive updates and inserts.
// Only *one* directory can be an entity group.
// Applications need to partition directories, if needed.
// Thus, the entire filesystem is extremely parallel.
//
// Major problem: How to traverse? How to implement ReadDir()?
// Answer: Use one global index of the Dir property.
// Such an index can be queried for equality.
//
// Worst disadvantage: Move operations, esp. in high level directories become expensive.
// Advantage: The directory "tree" can be sparse; only lowest dir must exist.
//
// Integrate into Afero.
//
// Unify/Extend the interface stuff. ReadDir or Readdir() ???
//
// Split into core package with interfaces only; plus several implementations?
//
//
// Add a "block"-layer under file,
// so that more than 1MB byte files can be writtens?
// At least throw an error before the file is saved?
//
// Mem Caching for directories
// Mem Caching for files - beware of cost
// Instance Caching with broadcasting instances via http request to instances.
//
// Use the walker function to implement removals
//
// SaveFile => optionally create non-existing directories
//   or return at least path.Error
//
// Nice to have: FileLinks

package gaefs
