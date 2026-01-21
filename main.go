package main

import (
	"github.com/XV-521/fileops/cnv"
	"github.com/XV-521/fileops/epub"
	"github.com/XV-521/fileops/internal"
	"github.com/XV-521/fileops/rename"
	"github.com/XV-521/fileops/zip"
	"os"
)

func main() {
	args := os.Args

	switch args[1] {
	case "highlight":
		err := internal.FlagWrapper(args, epub.HighlightWithFlags)
		if err != nil {
			panic(err)
		}
	case "rename":
		err := internal.FlagWrapper(args, rename.DoBatchWithFlags)
		if err != nil {
			panic(err)
		}
	case "cnv":
		err := internal.FlagWrapper(args, cnv.DoBatchWithFlags)
		if err != nil {
			panic(err)
		}
	case "zip":
		err := internal.FlagWrapper(args, zip.DoBatchWithFlags)
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown command")
	}
}
