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

type TreeFS struct {
	inputFS  fs.FS
	rootPath string
	sync.Mutex
	isLocked bool
	root     *pathNord
	asSlice  []*pathNord
}

// type TreeBldrState struct {
var lastNodePerDirLevel []*pathNord

// }

var pTFS *TreeFS

func (tfs *TreeFS) Lock() (success bool) {
	if tfs.isLocked {
		return false
	}
	tfs.isLocked = true
	tfs.Mutex.Lock()
	return true
}
func (tfs *TreeFS) Unlock() {
	if !tfs.isLocked {
		panic("Unlock failed")
	}
	tfs.Mutex.Unlock()
	tfs.isLocked = false
}
func (tfs *TreeFS) IsLocked() bool { return tfs.isLocked }

func NewTreeFS(path string, okayFilexts []string) *TreeFS {
	// var e error
	pTFS = new(TreeFS)
	pTFS.asSlice = make([]*pathNord, 0)
	pTFS.rootPath = path
	fmt.Println("on.newTreeFS.cwd:", pTFS.rootPath)
	pTFS.inputFS = os.DirFS(pTFS.rootPath)
	// func WalkDir(fsys FS, root string, wfn WalkDirFunc) error
	fs.WalkDir(pTFS.inputFS, ".", wfnBuildTree)
	return pTFS
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
	needsInit = (len(pTFS.asSlice) == 0 && len(lastNodePerDirLevel) == 0)
	didDoInit = (len(pTFS.asSlice) > 0 && len(lastNodePerDirLevel) > 0)
	if !(needsInit || didDoInit) {
		panic("doRoot: illegal state")
	}
	return needsInit
}

// type WalkDirFunc func(path string, d DirEntry, err error) error
func wfnBuildTree(path string, d fs.DirEntry, err error) error {
	// ROOT NODE ?
	if mustInitRoot() {
		pN := new(pathNord)
		pN.absFP = FU.AbsFP(pTFS.rootPath)
		pN.relFP = ""
		pTFS.root = pN
		pTFS.asSlice = append(pTFS.asSlice, pN)
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
	pNord := new(pathNord)
	pNord.relFP = path
	pNord.absFP = FU.AbsFP(FP.Join(pTFS.rootPath, path))
	// COUNT SLASHES
	nSlashes := S.Count(path, "/")

	// If Parent is Root
	if 0 == nSlashes {
		pTFS.root.AddKid(pNord) // (&pNode.Nord)
		return nil
	}
	// Parent is not root, so will have to locate parent.

	// Check length
	lenLNPDL := len(lastNodePerDirLevel)
	fmt.Printf("%d\n", lenLNPDL)
	// Find Parent

	return nil // FIXME
}
