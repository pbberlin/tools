Overview
--------------------
This package serves as universal abstraction
layer for all go programs that require a file system.

The interface, http-fs and the memory-fs were taken 
from [Steve Francia](https://twitter.com/spf13)'s afero:
Afero can be found [here](https:github.com/spf13/afero).

I had three cases in mind:

- Http services on app engine, who manage static resources.

- Image repositories on app engine.

- Content management systems.

Deployments can be made to local machines.
Files are editable with all editors.

### Package fsi - filesystem interface
Contains the minimal
requirements for exchangeable filesystems.

The required interface is intentionally ultra slim.
Contact me, if you think, the interface needs 
additional *mandatory* methods.
Then create a pull-request and add your method to the interface and to all filesystems.

#### Subpackage fsc 
Holds common extensions to all filesystems.
It contains an implementation of filepath.walk,
that can be used for all contained filesystems.
aefs.RemoveAll is built on this walk.

#### Subpackage fstests 
Contains tests for all filesystems.
Tests on file level are taken from [afero](https:github.com/spf13/afero).
I added a directory-tree-walk test suite and an httpfs-wrapper test.

#### osfs 
Osfs is the operating system fs. Replace 
	
	os.Open() and ioutil.WriteFile()
by 
	
	[filesys-instance].Open() and [filesys-instance].WriteFile()

Then you can switch between hard disk and memory filesystem
by changing the instantiation.

#### memfs
Keeps directories and files completely in RAM.
Good for quick testing. Cleanup is included.

#### aefs
With aefs you can write on google's datastore like onto a local hard disk.
See doc.go for details.


#### s3fs
Not yet adapted: A filesystem layer for amazon s3 and ceph

#### httpfs
This fs can wrap any previous filesystems and makes them serveable by a go http fileserver.

## Improvements

- memfs was substantially recoded.

- All filesystems are initialized with variadic option functions, as Dave Cheney [suggested](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).

## Todos

### Locking

- Can the locking approach of memfs be simplified?

- Can the locking approach of memfs at least be explained in comments?

- memfs registerWithParent locking seems neglected.

- aefs needs a locking consideration for RemoveAll and Rename. 
Behold the asynchroneus nature of aefs directories.


Common Remarks:
--------------------
All filesystems need to maintain compatibility
to relative paths of osfs; that is to current working directory prefixing.
Therefore all filesystems must support "." for current working dir.
Currently - in memfs and aefs - working dir always refers to fs.RootDir().

memfs and aefs interpret / or nothing as starting with root.

To access files directly under root, memfs and aefs must use ./filename

The filesystem types are no longer exported.
To access implementation specific functionality, use

	subpck.Unwrap(fsi.FileSystem) SpecificFileSys

Terminology:
--------------------
"name" or "filename" can mean either the basename or the full path of the file,
depending on the actual argument:

	'app1.log'              # simply
	'/tmp/logs/app1.log'    # full

In the first case, it refers to 

	[current dir]/app1.log.

Which is for memfs and aefs        

	[root dir]/app1.log.

Exception: 

	os.FileInfo.Name() # always contains *only* the base name.


Compare [discussion](http:stackoverflow.com/questions/2235173/file-name-path-name-base-name-naming-standard-for-pieces-of-a-path)
