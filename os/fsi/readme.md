Overview
--------------------
This package serves as universal abstraction
layer for all go programs that require a file system.

The interface, http-fs and the memory-fs were taken 
from [Steve Francia](https://twitter.com/spf13)'s afero:
Afero can be found [here](https:github.com/spf13/afero).

I had these purposes in mind:

- Deploying content management systems on app engine.<br>
Combining any apps with static pages and static HTML UI wrappers.

- Any HTTP services on app engine, managing static resources.


A similar effort was made by [rainycape](https://github.com/rainycape/vfs).

Deployments can be made to local machines.
Files are editable with all editors.

### Package fsi - filesystem interface
Contains the minimal
requirements for exchangeable filesystems.

The required interface did not remain ultra slim,
since apps need accessing all functions *via interface*.
They cannot access conditional methods easily.

All Fileystems provide a pck.Unrwap(fsi)(fsImpl,bool) to get access to
the underlying fileystem. But we dont want to unwrap, only to access fsImpl.Chmod(...).
We need fsiX.Chmod() directly.


#### Subpackage common
Holds common extensions to all filesystems.
It contains an implementation of filepath.walk,
that can be used for all contained filesystems.

dsfs.RemoveAll is built on this walk.

Contains standardization logic for paths.


#### Subpackage tests 
Contains tests for all filesystems.

File level tests are taken from [afero](https:github.com/spf13/afero).

I added a directory-tree test suite, with lots of subdirectories and files, and file-retrieval and removal.

I also added a httpfs-wrapper test (incomplete).

Tests for path standardization logic.

#### osfs 
Osfs is the wrapped operating system. Replace 
	
	os.Open() and ioutil.WriteFile()
by 
	
	[filesys-instance].Open() and [filesys-instance].WriteFile()

Then you can switch between hard disk and memory filesystem
by changing the instantiation.

osfs expects unix style paths, even under windows.
Give it c:/dir/file.


#### memfs
Keeps directories and files completely in RAM.
Good for quick testing and fast cleanup.

You can change the MountPoint of a memfs after its creation,
thus storing multiple trees consecutively, but not concurrently.


#### stacking filesystems
You can add a "shadow" filesystem to memfs.
Memfs will now act as cache for the shadow filesystem.

Its underimplemented though. Only Open() is looking.


#### dsfs
With dsfs you can write on google's datastore like onto a local hard disk.
See doc.go for details.


#### s3fs
Not yet adapted: A filesystem layer for amazon s3 or ceph


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


Fat architectural distinction to Afero
--------------------
fileinfo.Name() should return the filename without path.

I think this keeps the file independent from its current location in a directory tree.

It also means, the file has no method to access its current folder.
We have to know its location, in order to open it.

This is in striking contrast to the Afero implementation, where fileinfo.Name() returns the *full* path.
But the documentation for os.FileInfo.Name() is on my side.
Even Microsoft C# and C++ implementation [agree](https://msdn.microsoft.com/de-de/library/system.io.fileinfo.name(v=vs.110).aspx) and provide a separate FullName method.
