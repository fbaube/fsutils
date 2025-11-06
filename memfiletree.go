package fsutils

import(
	FU "github.com/fbaube/fileutils"
	"os"
)

type MemFileTree struct {
     *os.Root
     RootPaths   *FU.Filepaths
     RootNode    *FileTreeNode // but.. DRY ?!
     AsSlice   []*FileTreeNode // [0] pts to RootNode
     AsMapOfAbsFP map[string]*FileTreeNode 
     AsMapOfRelFP map[string]*FileTreeNode 
     *FU.FSItemSummaryStats 
}

