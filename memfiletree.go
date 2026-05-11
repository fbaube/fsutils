package fsutils

import(
	FU "github.com/fbaube/fileutils"
	N "github.com/fbaube/nork"
	"os"
)

type MemFileTree struct {
     *os.Root
     RootPaths   *FU.Filepaths
     RootNode    *N.FSONork // but.. DRY ?!
     AsSlice   []*N.FSONork // [0] pts to RootNode
     AsMapOfAbsFP map[string]*N.FSONork 
     AsMapOfRelFP map[string]*N.FSONork 
     *FU.FSObjectSummaryStats 
}

