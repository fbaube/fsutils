package fsutils

import(
	FU "github.com/fbaube/fileutils"
	"os"
)

type MemFileTree struct {
     *os.Root
     RootPaths   *FU.Filepaths
     RootNode    *FileTreeNork // but.. DRY ?!
     AsSlice   []*FileTreeNork // [0] pts to RootNode
     AsMapOfAbsFP map[string]*FileTreeNork 
     AsMapOfRelFP map[string]*FileTreeNork 
     *FU.FSObjectSummaryStats 
}

