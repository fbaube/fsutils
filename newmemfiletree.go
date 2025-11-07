package fsutils

import (
	"io/fs"
	"fmt"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	// SU "github.com/fbaube/stringutils"
	// CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
)

// NewMemFileTree proceeds as follows:
//  1. Walk the FS of a new [os.Root] to get a slice of filepath strings
//  2. Use that slice to build a slice of (ptrs to) [FileTreeNode] (via 
//     ptrs to [*fileutils/FSItem])
//  3. Provide the hierarchical/tree structure, by "weaving" the slice 
//     together (i.e. linking parents and children, probably using more 
//     than one method as implemented by [orderednodes/Nord]), and provide 
//     other means of access, such as a map from filepaths
//
// TBD: Maybe the path argument should be an absolute filepath, 
// because a relative filepath might cause problems. Altho
// this is the opposite of the advice for lower-level items.
//
// It isn't yet clear precisely how to use [os.Root]. Note tho that when
// we used [os.DirFS], it appeared to make no difference whether path 
//  - is relative or absolute
//  - ends with a trailing slash or not
//  - is a directory or a symlink to a directory
//
// The only error returns for this func are:
//  - a bad path, rejected by func FU.NewFilepaths
//  - the path is not a directory (altho it can be
//    a symlink to a directory ?)
//  - filepath-specific errors are in interface
//    [fileutils.Errer] and counted up in
//    [MemFileTree.FileSummaryStats]
// MemFileTree does not embed Errer and cannot
// itself return an error. 
//
// TODO: Maybe it needs two boolean arguments:
//  - One to say whether to be stricter about security using
//    funcs [io/fs.ValidPath] and [path/filepath.IsLocal], and
//  - One to say whether to follow symlinks (i.e. symlinks
//    that are legal by having targets under the Root)
//   - These two flags might have some interesting interactions
// .
func NewMemFileTree(aPath string, okayFilexts []string) (*MemFileTree, error) {
	var e error
	// var pFSS = new(FU.FSItemSummaryStats)

	L.L.Info("Making NewMemFileTree: " + aPath)
	// -------------------------------
	// 0. Check the filepath argument 
	// -------------------------------
	pMFT := new(MemFileTree)
	pMFT.Root, pMFT.RootPaths, e = FU.GetRootPaths(aPath)
	if e != nil {
		return nil, &fs.PathError{ Err:e, Path:aPath,
		       Op:"fsu.NewMemFileTree fu.GetRootPaths" } 
	}
	// -----------------------------------------
	// 1. Walk the os.Root's FS to gather a
	//    slice of simple strings of filepaths,
	//    and then use them to create [*FSItem].
	// -----------------------------------------
	var rFPs []string
	// var FSIs []*FU.FSItem
	// var pFSS *FU.FSItemSummaryStats
	rFPs = WalkFSforFilepathSlice(pMFT.FS())
	// FSIs, pFSS := FU.NewFSItemSliceFromFilepathSlice(rFPs)
	// -------------------------------------------
	// 2. Use the slice of FSItem'ss to build a
	//    slice of FileTreeNodes (which are just
	//    [Nord] plus [FSItem]) and the two maps 
	// ---------------------------------------------
	// FSIs is the same length as rFPs and each element
	// of FSIs implements interface [Errer]. So upgrade 
	// FSItems that do not have errors to FileTreeNode's.
	// --------------------------------------------------
	// pMFT.AsSlice   = make(         []*FileTreeNode, 0)
	pMFT.AsMapOfAbsFP = make(map[string]*FileTreeNode)
	pMFT.AsMapOfRelFP = make(map[string]*FileTreeNode)
	// It's a dir IFF it ends in a slash 
	for _, sFP := range rFPs { // range FSIs !!
	    // ----------------------------------
	    //  Form the path of the file-or-dir
	    //   and make the FileTreeNode
	    // ----------------------------------
	    /* absPathToUse := FU.EnsureTrailingPathSep(
		  	       FP.Join(pMFT.RootPaths.AbsFP, inPath)) */
	    pFTN := NewFileTreeNode(sFP) // (absPathToUse)
	    pFSI := &(pFTN.Fsi)
	    if pFSI.HasError() {
	        e = pFSI.GetError()
		L.L.Error("New FileTreeNode(%s) failed: %T %+v", sFP, e, e)
		pMFT.NrErrors++
		continue // keep on truckin' 
	    }
	    /*
	    // ---------------------------------
	    //  Do something based on just 
	    //  what exactly the input DirEntry
	    //  (inDE) is - file, dir, wotevs.
	    // ---------------------------------
	    // This is where bugs can appear when it's a dir.
	    // TODO: Not sure what happens with symlinks. 
	    if pFSI.IsDir() {
	   	if pFSI.TypedRaw == nil {
		   pFSI.TypedRaw = new(CT.TypedRaw)
		   } 
	        pFSI.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE
		pMFT.nDirs++ // just a simple counter
	    } else if pFSI.Type() == 0 { // regular file
		pMFT.nFiles++ // just a simple counter
	    } else if (pFSI.Type() & fs.ModeSymlink) != 0 { // Symlink
	       	if pFSI.TypedRaw == nil {
		   pFSI.TypedRaw = new(CT.TypedRaw)
		}
	        pFSI.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE // OOPS
		pMFT.nMiscs++ // just a simple counter 
		L.L.Okay("Item (SYML) OK: what to do ?!")
	    } else { // Some weirdness in the Mode bits 
	        pFSI.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE
             // pMFT.nMiscs++ // just a simple counter
		pMFT.nErrors++
                L.L.Error("Item (WTF) BAD: what to do ?!")
	    }
	    */
	    pMFT.AsSlice = append(pMFT.AsSlice, pFTN)
	    // Also add it to the maps 
	    pMFT.AsSlice = append(pMFT.AsSlice, pFTN)
	    pMFT.AsMapOfAbsFP[pFSI.FPs.AbsFP] = pFTN
	    pMFT.AsMapOfRelFP[pFSI.FPs.RelFP] = pFTN
	    // L.L.Info("ADDED TO MAP L225: " + sFP)
	}




	// ---------------------------------
	// 4. Here we could do some further
	//    filtering, even interactively
	// ---------------------------------
	
	if e != nil {
		// L.L.Panic("NewpMFT.WalkDir: " + e.Error())
		return nil, &fs.PathError { Op:"NewMemFileTree.Walk",
		       Err:e, Path:aPath } 
	}
	L.L.Okay("NewpMFT: walked OK %d nords from path %s",
		 len(pMFT.AsSlice), aPath)

	// Debuggery
	var ii int
	var ftn *FileTreeNode
	for ii, ftn = range pMFT.AsSlice {
	    if ftn == nil {
	       L.L.Error ("OOPS, pMFT.asSlice[%02d] is NIL", ii)
	       continue
	    }
	    /* if ftn.FSItem == nil || ftn.FSItem.FileMeta == nil {
	       L.L.Error("WTF, man!")
	       continue } */
	    if ftn.Fsi.IsDirlike() {
	        L.L.Debug("[%02d] isDIRLIKE: AbsFP: %s",
			ii, ftn.Fsi.FPs.AbsFP)
	    } else {
		L.L.Debug("[%02d] MarkupType: %s", ii, ftn.Fsi.TypedRaw.Raw_type)
	    }
	}

	// ==============================================
	//      SECOND PASS
	//  Range over the slice, using the materialised
	//  paths in asMapToAbsFS to identify parent/kid 
	//  Nord relationships and link together
	// ==============================================
	// TODO: This needs to be in some generalized 
	// form, such as TreeFromMaterializedPaths
	// =========================================
	var i int
	var pC *FileTreeNode
	for i, pC = range pMFT.AsSlice {
		if i == 0 { // skip over root 
			continue
		}
		// ---------------------------
		//  Shortcut if child of root
		// ---------------------------
		if !S.Contains(pC.RelFP(), FU.PathSep) {
			pMFT.RootNode.AddKid(pC)
			continue
		}
		// --------------------------
		//   Get dir portion of path
		// --------------------------
		itsDir := FP.Dir(pC.RelFP())
		itsDir = FU.EnsureTrailingPathSep(itsDir)
		// println(n.Path, "|cnex2|", itsDir)
		// L.L.Warning("itsDir: " + itsDir)
		// L.L.Warning("theMap: %+v", pMFT.asMap)
		var pPar *FileTreeNode
		var ok bool
		// PROBLEMS HERE ?
		// The parent directory should be in the map.
		// If it's not, then possibly we have messed
		// up with trailing separators. 
		if pPar, ok = pMFT.AsMapOfAbsFP[itsDir]; !ok {
			L.L.Error("findParentInMap: failed for: " +
				itsDir + " of " + pC.AbsFP())
			println(fmt.Sprintf("%+v", pMFT.AsMapOfAbsFP))
			panic(pC.AbsFP())
		}
		/*
		if itsDir != par.AbsFP() { // or, Rel? 
			panic(itsDir + " != " + par.AbsFP())
		}
		*/
		pPar.AddKid(pC)
	}
	// TODO: Look for entries that do not have a parent assigned !
	
	/* more debugging
	println("DUMP LIST")
	for _, n := range pFTFS.AsSlice {
		println(n.LinePrefixString(), n.LineSummaryString())
	}
	println("DUMP MAP")
	for k, v := range pFTFS.AsMap {
		fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
	}
	*/
	// println(SU.Gbg("=== TREE ==="))
	// pMFT.rootNord.PrintAll(os.Stdout)
	return pMFT, nil
}
