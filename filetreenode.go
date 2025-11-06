package fsutils

import(
	ON "github.com/fbaube/orderednodes"
	FU "github.com/fbaube/fileutils"
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

