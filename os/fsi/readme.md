Overview
--------------------
This package serves as universal abstraction
layer for all go programs that require a file system.

The interface, http-fs and the memory-fs were taken 
from [Steve Francia](https://twitter.com/spf13)'s afero:
Afero can be found [here](https:github.com/spf13/afero).

I had these ideas in mind:

- Content management systems on app engine.<br>
Combining any apps with static pages and static HTML UI wrappers.

- Any HTTP services on app engine, managing static resources.


There is a similar effort going on by [rainycape](https://github.com/rainycape/vfs).

Deployments can be made to local machines.
Files are editable with all editors.

### Package fsi - filesystem interface
Contains the minimal
requirements for exchangeable filesystems.

The required interface did not remain ultra slim,
since apps need accessing all functions *via interface*.
They cannot access conditional methods easily.

All Fileystems provide an pck.Unrwap(fsi)(fsImpl,bool) to get access to
the underlying fileystem. But we dont want to unwrap, only to access fsImpl.Chmod(...).
We need fsiX.Chomod() directly.


#### Subpackage common
Holds common extensions to all filesystems.
It contains an implementation of filepath.walk,
that can be used for all contained filesystems.
dsfs.RemoveAll is built on this walk.

#### Subpackage tests 
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

You can change the MountPoint of a memfs after its creation,
thus storing multiple trees consecutively, but not concurrently.

You can add a "shadow" filesystem to memfs.
Memfs will now act as cache for the shadow filesystem.


#### dsfs
With dsfs you can write on google's datastore like onto a local hard disk.
See doc.go for details.


#### s3fs
Not yet adapted: A filesystem layer for amazon s3 and ceph


#### httpfs
httpfs can wrap any previous filesystem and make it serveable by a go http fileserver.

## Improvements

- memfs was substantially recoded.

- All filesystems are initialized with variadic option functions, as Dave Cheney [suggested](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).

- (Deterioration) appengine/aetest is now needed for all tests.

## Todos

#### Locking

- The locking approach of memfs is incomplete.<br>
memfs/0_init.go has details and solutions.<br>
It is also impossible to improve, since you constantly run into deadlocks, when calling Open, Close, Create ... since they often nest.

- memfs registerWithParent locking seems neglected.

- I probably sync all directory tree structures with a go-routine-for-select.

- dsfs needs a locking consideration for RemoveAll and Rename. 
Behold the asynchroneus nature of dsfs directories.

#### Tests

- A concurrent access test suite is missing.


Common Remarks
--------------------
All filesystems need to maintain compatibility
to relative paths of osfs; that is to current working directory prefixing.
Therefore all filesystems must support "." for current working dir.
Currently - in memfs and dsfs - working dir always refers to fs.RootDir().

memfs and dsfs interpret "", "/" and "." as starting with root.

To access files directly under root, memfs and dsfs one must use ./filename

The filesystem types are no longer exported.
To access implementation specific functionality, use

	subpck.Unwrap(fsi.FileSystem) SpecificFileSys

Terminology
--------------------
"name" or "filename" can mean either the basename or the full path of the file,
depending on the actual argument:

	'app1.log'              # simply
	'/tmp/logs/app1.log'    # full

In the first case, it refers to 

	[current dir]/app1.log.

Which is for memfs and dsfs        

	[root dir]/app1.log.

Exception: 

	os.FileInfo.Name() # always contains *only* the base name.


Compare [discussion on stackoverflow](http:stackoverflow.com/questions/2235173/file-name-path-name-base-name-naming-standard-for-pieces-of-a-path)
