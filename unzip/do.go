package unzip

import (
	"errors"
	"flag"
	"github.com/XV-521/fileops/internal"
	"os"
	"path/filepath"
)

func DoBatch(md *Mode) error {

	md, err := internal.Prepare(md)
	if err != nil {
		return err
	}

	bm := internal.BatchMode{
		Sem:    10,
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

	md := &Mode{
		SrcDir: *srcDir,
		DstDir: *dstDir,
		Pwd:    *pwd,
		Strict: *strict,
	}

	return DoBatch(md)
}
