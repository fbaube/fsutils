package fsutils

import(
	ON "github.com/fbaube/orderednodes"
	FU "github.com/fbaube/fileutils"
	"io/fs"
)

// FileTreeNode is not named "..Nord" because we don't need
// or use the ordering functionality of the embedded Nord.
type FileTreeNode struct {
     ON.Nord
     Fsi FU.FSItem
}

func NewFileTreeNode(absPath string) *FileTreeNode {
     p := new(FileTreeNode)
     p.Fsi = *FU.NewFSItem(absPath)
     return p
}

type MemFileTree struct {
     fs.FS 
     RootAbsPath  string 
     Root         FileTreeNode
     AsSlice  []*FileTreeNode // [0] pts to Root 
     AsMapOfAbsFP map[string]*FileTreeNode 
     AsMapOfRelFP map[string]*FileTreeNode 
     nItems, nFiles, nDirs, nMiscs, nErrors int
}

/*
type ContentityFS struct {
        // FS will be set from func [os.DirFS]
        fs.FS
        rootAbsPath string
        rootNord    *RootContentity
        asSlice     []*Contentity
        // For the maps, the string key USED TO be the relative filepath 
        // w.r.t. the rootAbsPath. Now we simplify it to AbsFP. It's not 
        // really crucial one way or the other cos this map is discarded
        // when the ContentityFS is saved to disk. But do the ptrs point
        // into the tree of Nord's or into the slice of Nords ? 
        asMapOfAbsFP   map[string]*Contentity
        asMapOfRelFP   map[string]*Contentity // TODO! 
        nItems, nFiles, nDirs, nMiscs, nErrors int
}
*/

