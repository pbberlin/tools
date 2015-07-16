This package should serve as universal abstraction
layer for all go programs that require a file system.

The interface, http-fs and the memory-fs were taken 
and adapted from Steve Francia's afero:
Afero can be found here: https://github.com/spf13/afero.

Files are made editable like local files.

Test-Deployments can be made to local machines.

I had three cases in mind:
* http services on app engine, who manage static resources.
* image repositories on app engine.
* content management systems.

The required interface is intentionally ultra slim.
Contact me, if you think, the interface needs 
additional *mandatory* methods.
Then via pull-request add your method to all filesystems.
