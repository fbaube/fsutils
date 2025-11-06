package fsutils

import (
	"io/fs"
	FU "github.com/fbaube/fileutils"
)

// WalkFSforFilepathSlice takes an [fs.FS] (assumed to be valid) 
// and returns a (filtered!) slice of strings that are filepaths.
// 
// About the results:
//  - in the returned slice, a filepath that is a directory or symlink 
//    will always end with a slash (not with [os.Separator])
//  - pass an os.Root(path).FS() to avoid security issues with symlinks
//  - results of the walk are returned WITH filtering (using defaults for the
//    m5 app), therefore omitting (e.g.) any path element starting with ".git"
// 
// TODO: Describe results for abs v rel
// 
// Maybe the path argument should be an absolute filepath, 
// because a relative filepath might cause problems. Altho
// this is the opposite of the advice for lower-level items.
//
// [os.Root] has not been exercised, but note tho that when using
// [os.DirFS], it appeared to make no difference whether path 
//  - is relative or absolute
//  - ends with a trailing slash or not
//  - is a directory or a symlink to a directory
// .
func WalkFSforFilepathSlice(anFS fs.FS) ([]string) {
     	var FPs []string
	// NOTE that rel.path "." seems to be necessary 
	// here or else really weird errors occur.
	e := fs.WalkDir(anFS, ".",
	    // 
	    //   fs.WalkDirFunc ("WD-func"): 
	    // type WalkDirFunc func(path string, d DirEntry, err error) error
	    //
	    // There are two errors involved: returned and input-arg. 
	    // Handling them: https://pkg.go.dev/io/fs#WalkDirFunc
	    // 
	    // The RETURNED error controls how WalkDir continues:
	    // If WD-func returns SkipAll, WalkDir skips ALL files and dirs.
	    // If WD-func returns SkipDir, WalkDir skips the current dir:
	    //    arg-path if d.IsDir(), else the arg-path's parent dir.
	    // If WD-func returns some other non-nil error (*fs.PathError ?),
	    //    WalkDir stops entirely and returns that error.
	    // So WD-func should do some error handling of its own if it
	    //    wants to return a nil error so that walking continues. 
	    // 
	    // The INPUT-ARG error reports an error related to path, a dir:
	    // WalkDir will not walk into the dir. WD-func must decide how to
	    // handle the error; as described above, returning the error (or
	    // any non-nil, non-Skip error) makes WalkDir stop walking entirely.
	    // 
	    // WalkDir calls WD-func with a non-nil error ARG in two cases:
	    // 1) If the initial Stat on the root dir fails, WalkDir calls
	    // WD-func with path set to root, d set to nil, and err set to
	    // the error from fs.Stat. (This should never happen if we have
	    // pre-validated the path, unless there is a race condition.)
	    // 2) If ReadDir on a dir fails (see [ReadDirFile]), WalkDir
	    // calls WD-func with: path is the dir's path, d is a DirEntry
	    // describing the dir, and err is ReadDir's error. In this 2nd
	    // case, WD-func is called twice with the path of the dir.
	    // The 1st call is before the dir read is attempted and has err
	    // set to nil(!), giving WD-func a chance to return SkipDir or
	    // SkipAll and avoid the ReadDir entirely.
	    // The 2nd call is after a failed ReadDir and reports the error
	    // from ReadDir. (In the normal case, ReadDir succeeds and there
	    // is no 2nd call.)
	    // . 
	    func(inPath string, inDE fs.DirEntry, inErr error) error {
	    // 
	    // This func filters out a default set of unwanted values
	    // (and if a dir is unwanted, it returns [fs.SkipDir]):
	    //  - hidden (esp'ly .git directory)
	    //  - leading underbars ("_")
	    //  - emacs backup ("myfile~")
	    //  - this app's info  files: "*gtk,*gtr"
	    //  - this app's debug files: "*_(echo,tkns,tree)"
	    //  - filenames without a dot (indicating no file extension)
	    //  - NOTE that zero-length files (no content to analyse)
	    //    should NOT be filtered out 
	    //
	    // As path separator, "/" is usually assumed, not [os.PathSep]. 
	    // 
	    // --------------------------
	    // Were we passed an error ?
	    // If this is the first call to WD-func then the 
	    // input path for the walk was bad, so panic.
	    // Else it's a dir and the call to ReadDir() failed 
	    // and we are getting the error from that call.
	    // --------------------------
	    if inErr != nil {
	       	if len(FPs) == 0 {
		   panic("WalkFSforFPslice(root): " + inErr.Error())
		   }
		if !inDE.IsDir() {
		   panic("WalkFSforFPslice(non-root): " + inErr.Error())
		   }
		return fs.SkipDir
	    }
	    // Make sure a dir has a trailing path separator 
	    if inDE.IsDir() {
	       inPath = FU.EnsureTrailingPathSep(inPath)
	    }
	    // -----------------------------------
	    //  Filter out unwanted stuff; should
	    //  this work on the first call too ? 
	    // -----------------------------------
	    bad, _ := FU.ExcludeFilepath_m5(inPath)
	    if bad {
		// L.L.Debug("Rejecting (%s): %s", inPath, rsn)
		if inDE.IsDir() { return fs.SkipDir } 
		return nil
	    }
	    FPs = append(FPs, inPath)
	    return nil
})
	if e != nil { return nil }
	return FPs
}

