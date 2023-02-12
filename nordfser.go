package fsutils

import (
	"io/fs"

	ON "github.com/fbaube/orderednodes"
)

// An mcfile.Contentity (including those that represent directories)
// should implement four methods listed below: Open, Readfile, Stat,
// and ReadDir. The a Contentity can be treated like an fs.File, and
// the latter three methods can be delegated to by the FS itself.

// NordFSer is implemented by all types that assemble a tree of Nords.
// For working with actual files and directories, use FileNordFSer instead.
// The godoc for this interface describes methods common to both.
//
// Tipicly a NordFSer is a tree of tags (e.g. an AST) parsed from an XML file.
// The tag name is the last tag element of the abs.path and/or the rel.path,
// but possibly the relative path is just the tag name. Furthermore, calling
// Open (the only method specified in the fs.FS interface) returns the tag's
// body (from opening tag to closing tag, inclusive) of that tag only, without
// the context of the larger tag tree.
type NordFSer interface {
	// Interface fs.FS is the minimum required of an fs file system.
	// The Open(path) method is its only method. Iy opens the named
	// file. An error should be of type *PathError with the Op "open",
	// Path set to the path, and Err describing the problem.
	// Open should refuse to open a path that fails ValidPath(path),
	// returning a *PathError with Err set to ErrInvalid or ErrNotExist.
	Open(path string) (fs.File, error)
	// Size returns (for a FileTree) the combined number of files and dirs,
	// or (for a TagTree) the total number of tags.
	Size() int
	// RootAbsPath can return the the abs.FP of the markup file if the FS is a TagTree.
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
// like in a fs.DirTreeFS), or (b) listed individually (like in a DITA map).
// The file or directory name is the last path element of either of (but
// maybe both) the absolute path and/or the relative path.
//
// Note that fs.File is an interface that providess just three methods:
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
