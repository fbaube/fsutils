package fsutils

import(
	NOrK "github.com/fbaube/nork"
	FU "github.com/fbaube/fileutils"
)

// FileTreeNork is TBS.
type FileTreeNork struct {
     NOrK.Nork
     Fsi FU.FSItem
}

func NewFileTreeNork(absPath string) *FileTreeNork {
     p := new(FileTreeNork)
     p.Fsi = *FU.NewFSItem(absPath)
     return p
}

