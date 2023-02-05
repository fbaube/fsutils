package fsutils

import (
	"fmt"
	"io/fs"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	ON "github.com/fbaube/orderednodes"
)

// ## var lastNodePerDirLevel []*ON.FileNord

// pFTFS is a singleton and NOT RE-ENTRANT!
var pFTFS *FileTreeFS

// NewFileTreeFS is duh.
func NewFileTreeFS(path string, okayFilexts []string) *FileTreeFS {
	pFTFS = new(FileTreeFS)
	// Initialize embedded baseFS
	pFTFS.BaseFS.rootAbsPath = path
	pFTFS.BaseFS.inputFS = os.DirFS(path)
	println("fss.newFileTreeFS:", pFTFS.BaseFS.rootAbsPath)
	// Initialize slice & map
	pFTFS.asSlice = make([]*ON.Nord, 0)
	pFTFS.asMap = make(map[string]*ON.Nord)

	// FIRST PASS
	// Load slice & map
	// Skip dodgy filenames
	e := fs.WalkDir(pFTFS.inputFS, ".", wfnBuildFileTree)
	if e != nil {
		panic("fss.newFileTreeFS: " + e.Error())
	}
	fmt.Printf("fss.newFileTreeFS: got %d nords \n", len(pFTFS.asSlice))

	// SECOND PASS
	// Go down slice to identify parent nords and link together.
	println("WARN'G: drop ZERO-LENGTHs! utils/fsutils/filetreefsnew.go L48")
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
			par.AddKid(n)
		}
	}
	/* more debugging
	println("DUMP LIST")
	for _, n := range pFTFS.asSlice {
		println(n.LinePrefixString(), n.LineSummaryString())
	}
	println("DUMP MAP")
	for k, v := range pFTFS.asMap {
		fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
	}
	*/
	// println("=== FileTreeFS TREE ===")
	// pFTFS.rootNord.PrintAll(os.Stdout)
	return pFTFS
}
