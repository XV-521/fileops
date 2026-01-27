package unzip

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
		return internal.IsThisExt(entry.Name(), "zip")
	}

	handler := func(entry os.DirEntry) error {
		srcPath := filepath.Join(md.SrcDir, entry.Name())
		return internal.Unzip(srcPath, md.DstDir, md.Pwd)
	}

	return internal.DoBatchWrapper(md.SrcDir, bm, filter, handler)
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

	pwd := fs.String(
		"pwd",
		"",
		"Password.",
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
		Pwd:    *pwd,
		Strict: *strict,
	}

	mdPtr, err = internal.Prepare(mdPtr)
	if err != nil {
		return err
	}

	return DoBatch(*mdPtr)
}
