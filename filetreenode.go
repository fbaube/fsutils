package fsutils

import(
	NOrK "github.com/fbaube/nork"
	FU "github.com/fbaube/fileutils"
)

// FileTreeNork is TBS.
type FileTreeNork struct {
     NOrK.Nork
     FSO FU.FSObject
}

func NewFileTreeNork(absPath string) *FileTreeNork {
     p := new(FileTreeNork)
     p.FSO = *FU.NewFSObject(absPath)
     return p
}

