package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/fbaube/fsutils"
	ON "github.com/fbaube/orderednodes"
)

type ContentityFS interface {
	fs.FS
	ON.Norder
	NewContentityRoot(path string)
	GetContentityRoot()
	// The first walk generates PathProps.
	// The second walk reads and processes each file.
	DoPathProps(path string) (interface{}, error)
	ReadContentity(path string) (interface{}, error)
}

func ReadContentity(fsys fs.FS, path string) (interface{}, error) {
	if fsys, ok := fsys.(ContentityFS); ok {
		return fsys.ReadContentity(path)
	}
	return nil, fmt.Errorf("read-contentity %s: operation not supported", path)
}

var theFS fs.FS

func main() {
	var cwd string
	if len(os.Args) > 1 {
		println("No args allowed! Call from target directory.")
		os.Exit(1)
	}
	cwd, _ = os.Getwd()
	fmt.Println("CWD:", cwd)
	theFS = fss.NewFileTreeFS(cwd, nil)
}
