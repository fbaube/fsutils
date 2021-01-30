package fss

import (
	"fmt"
	"io/fs"
	"os"
	S "strings"

	FU "github.com/fbaube/fileutils"
	ON "github.com/fbaube/orderednodes"
)

type PathPropsTreeFS struct {
	inputFS fs.FS
	root    *ON.PathPropsNord
	asSlice []*ON.PathPropsNord
}

var ptCWD string
var ptROOT *ON.PathPropsNord
var ptNEXSEQ int

func NewPathPropsTreeFS(path string, okayFilexts []string) *PathPropsTreeFS {
	// var e error
	var pFS *PathPropsTreeFS
	ptCWD = path
	fmt.Println("on.newpptfs:", ptCWD)
	pFS = new(PathPropsTreeFS)
	pFS.inputFS = os.DirFS(ptCWD)
	// func WalkDir(fsys FS, root string, fn WalkDirFunc) error
	ptROOT = nil
	ptNEXSEQ = 0
	fs.WalkDir(pFS.inputFS, ".", wfnBuildPPtree)
	return pFS
}

func (p *PathPropsTreeFS) Open(path string) (fs.File, error) {
	return p.Open(path)
}

/*
type DirEntry interface {
    IsDir() bool
    // Name returns the final element of the path (the base name).
    Name() string
    // Type returns a subset of the usual FileMode bits,
    // those returned by FileMode.Type().
    Type() FileMode
    Info() (FileInfo, error)
}
*/

// type wfnBuildPPtree func(path string, d DirEntry, err error) error
func wfnBuildPPtree(path string, d fs.DirEntry, err error) error {
	var pNode *ON.PathPropsNord
	var pPP *FU.PathProps
	// Filter out non-content
	if S.HasPrefix(path, ".") {
		return nil
	} // or SKIPDIR?
	if S.HasSuffix(path, "~") {
		return nil
	} // or SKIPDIR?
	fmt.Println(path)
	println("COUNT UP DIR SEP'RS TO GET LEVEL")
	pNode = new(ON.PathPropsNord)
	// func NewPathProps(rfp string) *PathProps {
	fmt.Printf("\t")
	pPP = FU.NewPathProps(path)
	pNode.PathProps = *pPP
	// ROOT ?
	if path == "." {
		if ptROOT != nil {
			panic("ptROOT not nil")
		}
		if ptNEXSEQ != 0 {
			panic("ptNEXSEQ not 0")
		}
		ptROOT = pNode
		ptNEXSEQ = 1
		return nil
	}
	if ptROOT == nil {
		panic("ptROOT is nil")
	}
	if ptNEXSEQ == 0 {
		panic("ptNEXSEQ is 0")
	}
	// If Parent is Root
	if !S.Contains(path, "/") {
		ptROOT.AddKid(&pNode.Nord)
		return nil
	}
	// Find Parent
	return nil
}
