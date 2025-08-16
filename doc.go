/*

Naming the object is a problem. It stores an in-memmory copy of
a small filesystem but it does not provide the API of io/fs.FS .

It provides tree structure for a collection of fileutils/FSItem's
and it records this structure in multiple, redundant ways (enabling
a certain level of error checking).

So in a name for the object, FSItem is unnecessary because a name
involving Tree implicitly has nodes for directories. 

So, not FSItemFS. FSITree? FSIBundle? FileTree? FileMemTree? MemFileTree? 

Package fsutils provides a MemFileTree, which is summarised as:
 - Created from a filesystem path
 - Created using an os/Root, which is saved but unexported 
 - Each node contains a fileutils/FSItem
 - Tree structure is provided for each node by an orderednodes/Nord
 - Items can be selectively omitted from the tree at its creation time 
 - Note that the "next step up" from this is an mcm/ContentityFS 
 - Note that it is not necessary to implement io/fs.FS (and related),
   which do not provide or need hierarchical structure
 - Note that in principle we _could_ provide funcs & methods associated
   with io/fs.FS (and related), like `func Open(path string) (File, error)`;
   however the return value File would have to return the file on-disk
   rather than anything in-memory, and this would create all kinds of
   sync problems; Open and other calls could of course be shadowed /
   overridden), but this would be overambitious, buggy, and unnecessary

A similar technique will be used to create an improved version
of mcm/ContentityFS (and maybe rename it to Contentitree).

MemFileTree will provide read-write capability for individual FSItem's,
but adding and deleting entire FSItems by name may prove very difficult,
depending on how orderednode/Nord implements tree structure. 

In the interest of simplicity and composibility, a new MemFileTree can
be created stepwise, like so: 
 - Walk the os.Root to gather a slice of simple strings of filepaths
 - Use that slice to build a slice of (ptrs to) fileutils/FSItems's
 - Use the information gathered, and input arguments, to filter out entries 
 - (Optional) Provide user interactivity for filtering out additional entries
 - (If needed) Compress the slice, by removing nil entries 
 - Provide the hierarchical/tree structure, by "weaving" the slice together
   (linking parents and children, probably using more than one method as
   implemented by orderednodes/Nord), and provide other means of access,
   such as a map from filepaths 

fileutils/FSItem implements four interfaces:
 - [io/fs.FileInfo]
 - [io/fs.DirEntry]
 - [fileutils.Errer] (via an embed) 
 - [stringutils.Stringser] (Echo, Infos, Debug)

Notes about the usage of os/Root:
 - MemFileTree will store provate copies of the os.Root and of the
   result of `func (r *os.Root) FS() io/fs.FS`
 - The instances of os/Root and io/fs.FS will have to be unexported
 - Note that altho we will be able to filter out entries, the technique
   used in FilteringFS will not be helpful here, because we do not want
   to return the actual on-disk file, because this would create huge
   problems with content sync, because we will be storing the file
   (and manipulating it) in-memory 
 - The embedded io/fs.FS implements the following interfaces, which will 
   be available to MemFileTree itself but not to users of MemFileTree:
 - io/fs.StatFS     : Stat(name string) (FileInfo, error) // *PathError
   (if it is a link,  returns the file it links to) 
 - io/fs.ReadFileFS : ReadFile(name string) ([]byte, error)
 - io/fs.ReadDirFS  : ReadDir(name string) ([]DirEntry, error) // sorted 
 - io/fs.ReadLinkFS : ReadLink(name string) (string, error) // *PathError
 - io/fs.ReadLinkFS : Lstat(name string) (FileInfo, error) // *PathError
   (if it is a link,  returns the link, does not try to follow the link) 
 
https://benjamincongdon.me/blog/2021/01/21/A-Tour-of-Go-116s-iofs-package/

The Go library allows for more complex behavior by providing other file-
system interfaces that can be composed on top of the base fs.FS interface,
such as ReadDirFS, which allows you to list the contents of a directory:

	type ReadDirFS interface {
	    FS
	    ReadDir(name string) ([]DirEntry, error)
	}

The FS.Open function returns the new fs.File “ReadStatCloser” interface type,
which gives you access to some common file functions:

	type File interface {
	    Stat() (FileInfo, error)
	    Read([]byte) (int, error)
	    Close() error
	}

However, one big caveat: conspicuously absent from the fs.File interface is
any ability to write files. The fs package is a R/O interface for filesystems.

https://lobste.rs/s/kixqgi/tour_go_1_16_s_io_fs_package

fstest.TestFS does more than just assert that a few files exist. It walks the
entire file tree in the file system you give it, checking that all the various
methods it can find are well-behaved and diagnosing a bunch of common mistakes
that file system implementers might make. For example it opens every file it
can find and checks that Read+Seek and ReadAt give consistent results. And
lots more. So if you write your own FS implementation, one good test you
should write is a test that constructs an instance of the new FS and then
passes it to fstest.TestFS for inspection.
*/
package fsutils
