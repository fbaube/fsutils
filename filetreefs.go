package fsutils

import (
	ON "github.com/fbaube/orderednodes"
	"io/fs"
	FP "path/filepath"
	S "strings"
)

/*
var SF ON.StringFunc
func init() {
	SF = ON.NordSummaryString
}
*/

type FileTreeFS struct {
	BaseFS
	rootNord *ON.Nord
	asSlice  []*ON.Nord
	asMap    map[string]*ON.Nord // string is Rel.Path
}

// Open is a dummy function, just here to satisfy an interface.
func (p *FileTreeFS) Open(path string) (fs.File, error) {
	return p.inputFS.Open(path)
}

func mustInitFtfsRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(pFTFS.asSlice) == 0 && len(pFTFS.asMap) == 0)
	didDoInit = (len(pFTFS.asSlice) > 0 && len(pFTFS.asMap) > 0)
	if !(needsInit || didDoInit) {
		panic("mustInitFtfsRoot: illegal state")
	}
	return needsInit
}

// wfnBuildFileTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
//
// It filters out several file types:
// - (TODO:) zero-length file (no content to analyse)
// - hidden (esp'ly .git directory)
// - emacs backup (myfile~)
// - this app's debug files: *_(echo,tkns,tre)
// - filenames without dot
// .
func wfnBuildFileTree(path string, d fs.DirEntry, err error) error {
	var p *ON.Nord
	// ROOT NODE ?
	if mustInitFtfsRoot() {
		if path != "." {
			println("wfnBuildFileTree:",
				"root path is not dot but instead:", path)
		}
		p = ON.NewRootNord(pFTFS.rootAbsPath, nil) // ON.NordSummaryString)
		pFTFS.rootNord = p
		println("wfnBuildFileTree: root node abs.FP:", p.AbsFP())
	} else {
		// Filter out some file types as described in the func commment
		if S.HasPrefix(path, ".") || // hidden
			S.Contains(path, "/.") || // hidden
			S.HasSuffix(path, "~") || // emacs backup
			S.Contains(path, "/.git/") || // git repo
			(len(path) >= 5 && // debug file via "-t" flag
				path[len(path)-5] == '_') ||
			S.Index(FP.Base(path), ".") == -1 { // untyped file
			// Don't print TOO much!
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
