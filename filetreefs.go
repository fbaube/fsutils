package fss

import (
	"fmt"
	"io/fs"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	ON "github.com/fbaube/orderednodes"
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

type FileTreeFS struct {
	baseFS /*
		inputFS  fs.FS
		rootPath string
		sync.Mutex
		isLocked bool */
	root    *ON.FileNord
	asSlice []*ON.FileNord
	asMap   map[string]*ON.FileNord // string is Rel.Path
}

// ## var lastNodePerDirLevel []*ON.FileNord

var pFTFS *FileTreeFS

// NewFileTreeFS is duh.
func NewFileTreeFS(path string, okayFilexts []string) *FileTreeFS {
	pFTFS = new(FileTreeFS)
	// Initialize embedded baseFS
	pFTFS.baseFS.rootPath = path
	pFTFS.baseFS.inputFS = os.DirFS(path)
	println("fss.newFileTreeFS:", pFTFS.baseFS.rootPath)
	// Initialize slice & map
	pFTFS.asSlice = make([]*ON.FileNord, 0)
	pFTFS.asMap = make(map[string]*ON.FileNord)

	// FIRST PASS
	// Load slice & map
	e := fs.WalkDir(pFTFS.inputFS, ".", wfnBuildFileTree)
	if e != nil {
		panic("fss.newFileTreeFS: " + e.Error())
	}
	fmt.Printf("fss.newFileTreeFS: got %d nords \n", len(pFTFS.asSlice))
	// SECOND PASS
	// Go down slice to identify parent nords and link together.
	for _, n := range pFTFS.asSlice {
		// Is child of root ?
		if !S.Contains(n.Path, FU.PathSep) {
			pFTFS.root.AddKid(n)
		} else {
			itsDir := FP.Dir(n.Path)
			var par *ON.FileNord
			var ok bool
			if par, ok = pFTFS.asMap[itsDir]; !ok {
				panic(n.Path)
			}
			par.AddKid(n)
		}
	}

	println("DUMP LIST")
	for _, n := range pFTFS.asSlice {
		println(n.LinePrefixString(), n.LineSummaryString())
	}
	println("DUMP MAP")
	for k, v := range pFTFS.asMap {
		fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
	}
	println("DUMP TREE")
	pFTFS.root.PrintAll(os.Stdout)
	return pFTFS
}

// Open is a dummy function, just here to satisfy an interface.
func (p *FileTreeFS) Open(path string) (fs.File, error) {
	return p.inputFS.Open(path)
}

/* type DirEntry interface {
    IsDir() bool
    Name()  string   // the final elm of the path (the base name)
    Type()  FileMode // those FileMode bits ret'd by FileMode.Type()
    Info() (FileInfo, error)
} */

func mustInitRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(pFTFS.asSlice) == 0 && len(pFTFS.asMap) == 0) // && len(lastNodePerDirLevel) == 0)
	didDoInit = (len(pFTFS.asSlice) > 0 && len(pFTFS.asMap) > 0)   // && len(lastNodePerDirLevel) > 0)
	if !(needsInit || didDoInit) {
		panic("mustInitRoot: illegal state")
	}
	return needsInit
}

// wfnBuildFileTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
func wfnBuildFileTree(path string, d fs.DirEntry, err error) error {
	var pN *ON.FileNord
	// ROOT NODE ?
	if mustInitRoot() {
		pN = new(ON.FileNord)
		pN.SetIsRoot(true)
		pN.Path = ""
		pN.AbsFilePath = FU.AbsFP(pFTFS.rootPath)
		pFTFS.root = pN
		pFTFS.asSlice = append(pFTFS.asSlice, pN)
		pFTFS.asMap[path] = pN
		// len is 0, but...
		// ## lastNodePerDirLevel = append(lastNodePerDirLevel, pN)
		// len is now 1
		fmt.Printf("Root node FP: rel<%s> abs<%s> \n", path, pN.AbsFilePath)
		return nil
	}
	// Filter out hidden and emacs backup
	if S.HasPrefix(path, ".") || S.Contains(path, "/.") || S.HasSuffix(path, "~") {
		// println("Path rejected:", path)
		return nil
	}
	// ALLOCATE AND INIT Nord
	pN = new(ON.FileNord)
	pN.Path = path
	pN.AbsFilePath = FU.AbsFP(FP.Join(pFTFS.rootPath, path))
	pFTFS.asSlice = append(pFTFS.asSlice, pN)
	pFTFS.asMap[path] = pN
	// println("Path OK:", pN.AbsFilePath)

	/*
		// If Parent is Root
		if 0 == nSlashes {
			pFTFS.root.AddKid(pN) // (&pNode.Nord)
			return nil
		}
	*/
	// Parent is not root, so will have to locate parent.

	// Check length
	// ## lenLNPDL := len(lastNodePerDirLevel)
	// fmt.Printf("%d\n", lenLNPDL)
	// Find Parent

	return nil // FIXME
}
