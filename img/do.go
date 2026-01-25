package img

import (
	"errors"
	"flag"
	"github.com/XV-521/fileops/internal"
	"os"
	"path/filepath"
)

func DoBatch(md Mode) error {

	bm := internal.BatchMode{
		Async:  true,
		Strict: md.Strict,
	}

	filter := func(entry os.DirEntry) bool {
		if entry.IsDir() {
			return false
		}
		if !internal.IsThisExt(entry.Name(), md.Ext) {
			return false
		}
		return true
	}

	resizeHandler := func(entry os.DirEntry) error {
		srcPath := filepath.Join(md.SrcDir, entry.Name())
		dstPath := filepath.Join(md.DstDir, entry.Name())
		return resize(srcPath, dstPath, md.Rto)
	}
	err := internal.DoBatchWrapper(md.SrcDir, bm, filter, resizeHandler)
	if err != nil {
		return err
	}

	changeDpiHandler := func(entry os.DirEntry) error {
		srcPath := filepath.Join(md.DstDir, entry.Name())
		dstPath := filepath.Join(md.DstDir, entry.Name())
		return changeDpi(srcPath, dstPath, md.DPI)
	}
	return internal.DoBatchWrapper(md.DstDir, bm, filter, changeDpiHandler)
}

func DoBatchWithFlags(fs *flag.FlagSet, args []string) error {
	srcDir := fs.String(
		"src",
		"",
		"Source directory.",
	)

	dstDir := fs.String(
		"dst",
		"",
		"Destination directory.",
	)

	ext := fs.String(
		"ext",
		"",
		"Filter files by extension.",
	)

	rto := fs.Float64(
		"rto",
		0,
		"Ratio of the new image to the old image.",
	)

	dpi := fs.Float64(
		"dpi",
		0,
		"DPI of the new image.",
	)

	strict := fs.Bool(
		"strict",
		false,
		"Stop processing on the first error.",
	)

	err := fs.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	mdPtr := &Mode{
		SrcDir: *srcDir,
		DstDir: *dstDir,
		Ext:    *ext,
		Rto:    *rto,
		DPI:    *dpi,
		Strict: *strict,
	}

	mdPtr, err = internal.Prepare(mdPtr)
	if err != nil {
		return err
	}

	return DoBatch(*mdPtr)
}
