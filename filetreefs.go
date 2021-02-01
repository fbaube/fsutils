package fss

import (
	"fmt"
	"io/fs"
	"os"
	FP "path/filepath"
	S "strings"
	"sync"

	FU "github.com/fbaube/fileutils"
)

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

type DirTreeFS struct {
	inputFS  fs.FS
	rootPath string
	sync.Mutex
	isLocked bool
	root     *dirPathNord
	asSlice  []*dirPathNord
}

// type TreeBldrState struct {
var lastNodePerDirLevel []*dirPathNord

// }

var pDTFS *DirTreeFS

func (dtfs *DirTreeFS) Lock() (success bool) {
	if dtfs.isLocked {
		return false
	}
	dtfs.isLocked = true
	dtfs.Mutex.Lock()
	return true
}
func (dtfs *DirTreeFS) Unlock() {
	if !dtfs.isLocked {
		panic("Unlock failed")
	}
	dtfs.Mutex.Unlock()
	dtfs.isLocked = false
}
func (dtfs *DirTreeFS) IsLocked() bool { return dtfs.isLocked }

func NewDirTreeFS(path string, okayFilexts []string) *DirTreeFS {
	// var e error
	pDTFS = new(DirTreeFS)
	pDTFS.asSlice = make([]*dirPathNord, 0)
	pDTFS.rootPath = path
	fmt.Println("on.newTreeFS.cwd:", pDTFS.rootPath)
	pDTFS.inputFS = os.DirFS(pDTFS.rootPath)
	// func WalkDir(fsys FS, root string, wfn WalkDirFunc) error
	fs.WalkDir(pDTFS.inputFS, ".", wfnBuildTree)
	return pDTFS
}

/* // Open is a dummy function, just here to satisfy an interface.
func (p *TreeFS) Open(path string) (fs.File, error) {
	return nil, nil } */

/* type DirEntry interface {
    IsDir() bool
    Name()  string   // the final elm of the path (the base name)
    Type()  FileMode // those FileMode bits ret'd by FileMode.Type()
    Info() (FileInfo, error)
} */

func mustInitRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(pDTFS.asSlice) == 0 && len(lastNodePerDirLevel) == 0)
	didDoInit = (len(pDTFS.asSlice) > 0 && len(lastNodePerDirLevel) > 0)
	if !(needsInit || didDoInit) {
		panic("doRoot: illegal state")
	}
	return needsInit
}

// type WalkDirFunc func(path string, d DirEntry, err error) error
func wfnBuildTree(path string, d fs.DirEntry, err error) error {
	// ROOT NODE ?
	if mustInitRoot() {
		pN := new(dirPathNord)
		pN.absFP = FU.AbsFP(pDTFS.rootPath)
		pN.argPath = ""
		pDTFS.root = pN
		pDTFS.asSlice = append(pDTFS.asSlice, pN)
		// len is 0, but...
		lastNodePerDirLevel = append(lastNodePerDirLevel, pN)
		// len is now 1
		println("Did root node; cwd:", pN.absFP)
		return nil
	}
	// Filter out non-content
	if S.HasPrefix(path, ".") || S.Contains(path, "/.") || S.HasSuffix(path, "~") {
		println("Path rejected:", path)
		return nil
	}
	// ALLOCATE AND INIT Nord
	pNord := new(dirPathNord)
	pNord.argPath = path
	pNord.absFP = FU.AbsFP(FP.Join(pDTFS.rootPath, path))
	// COUNT SLASHES
	nSlashes := S.Count(path, "/")

	// If Parent is Root
	if 0 == nSlashes {
		pDTFS.root.AddKid(pNord) // (&pNode.Nord)
		return nil
	}
	// Parent is not root, so will have to locate parent.

	// Check length
	lenLNPDL := len(lastNodePerDirLevel)
	fmt.Printf("%d\n", lenLNPDL)
	// Find Parent

	return nil // FIXME
}
