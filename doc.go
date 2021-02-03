/*
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
package fss
