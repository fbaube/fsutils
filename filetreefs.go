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
var SF ON.StringFunc
func init() {
	SF = ON.NordSummaryString
}
*/

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
	baseFS
	rootNord *ON.Nord
	asSlice  []*ON.Nord
	asMap    map[string]*ON.Nord // string is Rel.Path
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
	pFTFS.asSlice = make([]*ON.Nord, 0)
	pFTFS.asMap = make(map[string]*ON.Nord)

	// FIRST PASS
	// Load slice & map
	e := fs.WalkDir(pFTFS.inputFS, ".", wfnBuildFileTree)
	if e != nil {
		panic("fss.newFileTreeFS: " + e.Error())
	}
	fmt.Printf("fss.newFileTreeFS: got %d nords \n", len(pFTFS.asSlice))

	// SECOND PASS
	// Go down slice to identify parent nords and link together.
	for i, n := range pFTFS.asSlice {
		if i == 0 {
			continue
		}
		// Is child of root ?
		if !S.Contains(n.Path(), FU.PathSep) {
			pFTFS.rootNord.AddKid(n)
			// ON.AddKid2(pFTFS.rootNord, n)
		} else {
			itsDir := FP.Dir(n.Path())
			// println(n.Path, "|cnex2|", itsDir)
			var par *ON.Nord
			var ok bool
			if par, ok = pFTFS.asMap[itsDir]; !ok {
				panic(n.Path)
			}
			if itsDir != par.Path() {
				panic(itsDir + " != " + par.Path())
			}

			// ON.AddKid3(&(par.Nord), &(n.Nord))
			par.AddKid(n)

			// fmt.Printf("ftfs: ptrs? n<%T> par<%T> \n", n, par)
			// fmt.Printf("ftfs: ptrs? n<%T> par<%T> \n", &(n.Nord), &(par.Nord))
			/*
				plk := par.LastKid()
				plk2 := plk.(*ON.Nord)
				if uintptr(unsafe.Pointer(n)) == uintptr(unsafe.Pointer(plk2)) {
					println("EQUAL!!??")
				}
				plk = n.Parent()
				plk2 = plk.(*ON.Nord)
				if uintptr(unsafe.Pointer(par)) == uintptr(unsafe.Pointer(plk2)) {
					println("EQUAL!!??")
				}
				if par.LastKid() == n {
					fmt.Printf("**** OK LINK 1??? %p,%p \n", n, par.LastKid())
				}
				if par.LastKid() != n {
					fmt.Printf("**** FAILED LINK 1??? %p,%p \n", n, par.LastKid())
				}
				if n.Parent() == par {
					fmt.Printf("**** OK LINK 2??? %p,%p \n", par, n.Parent())
				}
				if n.Parent() != par {
					fmt.Printf("**** FAILED LINK 2??? %p,%p \n", par, n.Parent())
				}
			*/
		}
	}
	/*
		println("DUMP LIST")
		for _, n := range pFTFS.asSlice {
			println(n.LinePrefixString(), n.LineSummaryString())
		}
		println("DUMP MAP")
		for k, v := range pFTFS.asMap {
			fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
		}
	*/
	println("=== TREE ===")
	pFTFS.rootNord.PrintAll(os.Stdout)
	return pFTFS
}

// Open is a dummy function, just here to satisfy an interface.
func (p *FileTreeFS) Open(path string) (fs.File, error) {
	return p.inputFS.Open(path)
}

func mustInitRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(pFTFS.asSlice) == 0 && len(pFTFS.asMap) == 0)
	didDoInit = (len(pFTFS.asSlice) > 0 && len(pFTFS.asMap) > 0)
	if !(needsInit || didDoInit) {
		panic("mustInitRoot: illegal state")
	}
	return needsInit
}

// wfnBuildFileTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
func wfnBuildFileTree(path string, d fs.DirEntry, err error) error {
	var p *ON.Nord
	// ROOT NODE ?
	if mustInitRoot() {
		if path != "." {
			println("wfnBuildFileTree: root path is not dot but instead:", path)
		}
		p = ON.NewRootNord(pFTFS.rootPath, nil) // ON.NordSummaryString)
		pFTFS.rootNord = p
		println("wfnBuildFileTree: root node abs.FP:", p.AbsFP())
	} else {
		// Filter out hidden (esp'ly .git) and emacs backup.
		// Note that "/" is assumed, not os.Sep
		if S.HasPrefix(path, ".") || S.Contains(path, "/.") ||
			S.HasSuffix(path, "~") || S.Contains(path, "/.git/") {
			if !S.Contains(path, "/.git/") {
				println("Path rejected:", path)
			}
			return nil
		}
		p = ON.NewNord(path)
	}
	pFTFS.asSlice = append(pFTFS.asSlice, p)
	pFTFS.asMap[path] = p
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
