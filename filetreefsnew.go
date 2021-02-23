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

// ## var lastNodePerDirLevel []*ON.FileNord

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
			par.AddKid(n)
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
	// println("=== FileTreeFS TREE ===")
	// pFTFS.rootNord.PrintAll(os.Stdout)
	return pFTFS
}
