package rename

import (
	"errors"
	"flag"
	"github.com/XV-521/fileops/internal"
	"os"
	"path/filepath"
)

func DoBatch(md Mode) error {

	bm := internal.BatchMode{
		Async:  false,
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

	handler := func(entry os.DirEntry) error {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		err = os.Rename(
			filepath.Join(md.SrcDir, info.Name()),
			filepath.Join(md.SrcDir, md.Namer.Next(info)),
		)
		if err != nil {
			return err
		}
		return nil
	}

	return internal.DoBatchWrapper(md.SrcDir, bm, filter, handler)
}

func DoBatchWithFlags(fs *flag.FlagSet, args []string) error {

	srcDir := fs.String(
		"src",
		"",
		"Source directory.",
	)
	basename := fs.String(
		"basename",
		"",
		"Base filename.",
	)
	ext := fs.String(
		"ext",
		"",
		"Filter files by extension.",
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
		SrcDir:   *srcDir,
		Basename: *basename,
		Ext:      *ext,
		Strict:   *strict,
	}

	mdPtr, err = internal.Prepare(mdPtr)
	if err != nil {
		return err
	}

	return DoBatch(*mdPtr)
}
