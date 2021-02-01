package fss

import (
	FU "github.com/fbaube/fileutils"
	ON "github.com/fbaube/orderednodes"
)

type dirPathNord struct {
	ON.Nord
	argPath string // relFP
	absFP   FU.AbsFilePath
}
