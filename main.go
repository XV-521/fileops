package main

import (
	"flag"
	"github.com/XV-521/fileops/cnv"
	"github.com/XV-521/fileops/epub"
	"github.com/XV-521/fileops/rename"
	"github.com/XV-521/fileops/zip"
	"os"
)

type FlagFn func(fs *flag.FlagSet, args []string) error

func WithFlagSet(args []string, flagFn FlagFn) error {
	fs := flag.NewFlagSet(args[1], flag.ContinueOnError)
	return flagFn(fs, args[2:])
}

func main() {
	args := os.Args

	switch args[1] {
	case "highlight":
		err := WithFlagSet(args, epub.HighlightWithFlags)
		if err != nil {
			panic(err)
		}
	case "rename":
		err := WithFlagSet(args, rename.DoBatchWithFlags)
		if err != nil {
			panic(err)
		}
	case "cnv":
		err := WithFlagSet(args, cnv.DoBatchWithFlags)
		if err != nil {
			panic(err)
		}
	case "zip":
		err := WithFlagSet(args, zip.DoBatchWithFlags)
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown command")
	}
}
