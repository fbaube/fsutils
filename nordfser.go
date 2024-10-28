package fsutils

import (
	"io/fs"

	ON "github.com/fbaube/orderednodes"
)

// https://pkg.go.dev/io/fs
// func ReadFile(fsys FS, name string) ([]byte, error)
// func WalkDir (fsys FS, root string, fn WalkDirFunc) error
// type         FS interface { Open    (name string) (File, error)
// type  ReadDirFS interface { ReadDir (name string) ([]DirEntry, error)
// type ReadFileFS interface { ReadFile(name string) ([]byte, error)
// type     StatFS interface { Stat    (name string) (FileInfo, error)

// https://benjamincongdon.me/blog/2021/01/21/A-Tour-of-Go-116s-iofs-package/
// FS.Open returns fs.File, a "ReadStatCloser":
// type File interface {
//      Stat()  (FileInfo, error)
//      Read([]byte) (int, error)
//      Close() error }
// These are a small subset of methods on os.File.
// The above-er list adds FS's ReadFile & ReadDir.
// Pkg os has os.ReadFile, and os.File has methods ReadDir & Readdir .
// Oddly, func os.DirFS(dir string) fs.FS
// returns an fs.FS, and "the result implements
// io/fs.StatFS, io/fs.ReadFileFS and io/fs.ReadDirFS",
// but does it accept rel.paths yet also reject a leading "." ? 

// https://github.com/golang/go/issues/47803
// - Note, all fs.FS paths are "unrooted" so there is
//   still a fundamental difference with os.Open, etc.
// - DirFS is unrooted and does not allow access outside
//    of the initial directory.
// - Note that io/fs.FS.Open supports only unrooted paths.
// - fs.ValidPath rejects rooted paths. 
// - os.Open handles either a rooted path or a cwd path.
// - fs.FS doesn't have concepts like "relative to CWD"
//   and "OS root directory". 
// - fs.FS very clearly disallows paths starting with / 
//   or with ../ or even with ./ - they must be rejected 
//   by any valid Open implementation.

// https://pkg.go.dev/path/filepath#Rel
// func Rel(basepath, targpath string) (string, error)
// targpath == Join(basepath, Rel(basepath, targpath))
// An error is returned if
// - targpath can't be made relative to basepath, or
// - knowing the CWD would be necessary to compute it.
// Rel calls Clean on the result.

// NordFSer is implemented by all types that assemble a tree of Nords.
// NOTE: This godoc is clearly out of date. 
// 
// NOTE: For working with actual files and directories, use the superset
// FileNordFSer instead. The godoc for this interface describes methods
// common to both.
//
// NordFSer is what an mcfile.Contentity (including directories) should
// implement: Open, Readfile, Stat, and ReadDir.
// Then a Contentity can be treated like an fs.File, and
// the latter three methods can be delegated to by the FS itself.
//
// Tipicly a NordFSer is a tree of tags (e.g. an AST) parsed from an XML file.
// The tag name is the last tag element of the abs.path and/or the rel.path,
// but possibly the relative path is just the tag name. Furthermore, calling
// Open (the only method specified in the fs.FS interface) returns the tag's
// body (from opening tag to closing tag, inclusive) of that tag only, without
// the context of the larger tag tree.
type NordFSer interface {
	// Interface fs.FS is the minimum required of an fs file system.
	// The Open(path) method is its only method. It opens the named
	// file. An error should be of type *PathError with the Op "open",
	// Path set to the path, and Err describing the problem.
	// Open should refuse to open a path that fails ValidPath(path),
	// returning a *PathError with Err set to ErrInvalid or ErrNotExist.
	Open(path string) (fs.File, error)
	// Size returns (for a FileTree) the combined number of files and dirs,
	// or (for a TagTree) the total number of tags.
	Size() int
	// RootAbsPath can return the the abs.FP of the markup file
	// if the FS is a TagTree.
	RootAbsPath() string
	// AllPaths can be either all absolute paths or all relative paths.
	AllPaths() []string
	RootNord() ON.Norder
	AsSlice() []ON.Norder
	AsMap() map[string]ON.Norder
}

// FileNordFSer is implemented by all types that assemble a tree of
// Nords and embed an fs.FS . For working with tags (or anything else
// that is not actual files and directories), use NordFSer instead.
// The godoc for this interface describes methods unique to FileNordFSer.
//
// A FileTree NordFSer is a collection of files & directories - either
// (1) hierarchically gathered & organized (when a directory tree is walked,
// like in a fs.DirTreeFS), or (b) listed individually (like in a DITA map
// or a MapFS or a list of materialized paths). The file or directory name
// is the last path element of either of (but maybe both) the absolute path
// and/or the relative path.
//
// Note that fs.File is an interface that provides just three methods:
// Stat() (FileInfo, error) ; Read([]byte) (int, error) ; Close() error
type FileNordFSer interface {
	NordFSer
	// InputFS is undefined for a TagTree.
	InputFS() fs.FS
	// DirCount returns zero if the FS is a TagTree.
	DirCount() int
	// FileCount returns zero if the FS is a TagTree.
	FileCount() int

	// Interface ReadFileFS is a file system providing a custom ReadFile(path).
	// ReadFile reads the named file and returns its contents.
	// A successful call returns a nil error, not io.EOF.
	ReadFile(path string) ([]byte, error)
	// Interface StatFS is a file system with a Stat method.
	// Stat returns a FileInfo describing the file.
	// Any error should be type *PathError.
	Stat(path string) (fs.FileInfo, error)
	// Interface ReadDirFS is a file system with a custom ReadDir(path).
	// ReadDir reads the named directory and returns a
	// list of directory entries sorted by filename.
	ReadDir(path string) ([]fs.DirEntry, error)

	// A ReadDirFile is an fs.File for a directory file, whose entries can
	// be read with ReadDir(n). Every directory fs.File should implement
	// this interface. (It is OK for any file to implement this interface,
	// but if so, ReadDir should return an error for non-directories.)
	// See https://tip.golang.org/pkg/io/fs/#ReadDirFile
	// ReadDir(n int) ([]fs.DirEntry, error)
}
