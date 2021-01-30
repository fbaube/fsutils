package fss

import (
  FU "github.com/fbaube/fileutils"
  ON "github.com/fbaube/orderednodes"
)

type pathNord struct {
	ON.Nord
	relFP  string
	absFP  FU.AbsFilePath
}
