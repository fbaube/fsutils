package fss

import (
	"io/fs"
	S "strings"

	ON "github.com/fbaube/orderednodes"
)

/*
var SF ON.StringFunc
func init() {
	SF = ON.NordSummaryString
}
*/

type FileTreeFS struct {
	baseFS
	rootNord *ON.Nord
	asSlice  []*ON.Nord
	asMap    map[string]*ON.Nord // string is Rel.Path
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
