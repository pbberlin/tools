This package might serve as universal abstraction
layer for all go programs that require a file system.

I had three cases in mind:
* Large image repositories.
* Classic content management systems.
* Any http server applications, that need static html resources.

Test-Deployments can be made to local machines.
Files should be editable like local file system files.

The required interface is intentionally ultra slim.
Contact me, if you think, the interface needs 
additional *mandatory* methods.
Then via pull-request add your method to all filesystems.
