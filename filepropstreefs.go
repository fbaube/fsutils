package fss

import (
	"fmt"
	"io/fs"
	"os"
	S "strings"

	FU "github.com/fbaube/fileutils"
	ON "github.com/fbaube/orderednodes"
)

type FilePropsTreeFS struct {
	BaseFS
	rootNord *ON.FilePropsNord
	asSlice  []*ON.FilePropsNord
	asMap    map[string]*ON.FilePropsNord // string is Rel.Path
}

// var ptNEXSEQ int

var pFPTFS *FilePropsTreeFS

func NewFilePropsTreeFS(path string, okayFilexts []string) *FilePropsTreeFS {
	// var e error
	pFPTFS = new(FilePropsTreeFS)
	pFPTFS.asSlice = make([]*ON.FilePropsNord, 0)
	pFPTFS.rootAbsPath = path
	fmt.Println("on.newFilePropsTreeFS.cwd:", pFPTFS.rootAbsPath)
	pFPTFS.inputFS = os.DirFS(pFPTFS.rootAbsPath)
	// func WalkDir(fsys FS, root string, wfn WalkDirFunc) error
	fs.WalkDir(pFPTFS.inputFS, ".", wfnBuildFilePropsTree)
	return pFPTFS
}

func (p *FilePropsTreeFS) Open(path string) (fs.File, error) {
	return p.Open(path)
}

// type wfnBuilwfnBuildFilePropsTreedPPtree func(path string, d DirEntry, err error) error
func wfnBuildFilePropsTree(path string, d fs.DirEntry, err error) error {
	var pNode *ON.FilePropsNord
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
	pNode = new(ON.FilePropsNord)
	// func NewPathProps(rfp string) *PathProps {
	fmt.Printf("\t")
	pPP = FU.NewPathProps(path)
	pNode.PathProps = *pPP
	// ROOT ?
	if path == "." {
		if pFPTFS.rootNord != nil {
			panic("pFPTFS.root not nil")
		}
		pFPTFS.rootNord = pNode
		return nil
	}
	if pFPTFS.rootNord == nil {
		panic("pFPTFS.root is nil")
	}
	// If Parent is Root
	if !S.Contains(path, "/") {
		pFPTFS.rootNord.AddKid(&pNode.Nord)
		return nil
	}
	// Find Parent
	return nil
}
