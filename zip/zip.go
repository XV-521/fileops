package zip

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/internal"
	"os"
	"path/filepath"
)

func DoBatch(md Mode) error {
	getZipName := func(fileName string) string {
		basename, _ := internal.GetBasenameAndExt(fileName)
		return fmt.Sprintf("%v.%v", basename, "zip")
	}

	bm := internal.BatchMode{
		Async:  true,
		Strict: md.Strict,
	}

	filter := func(entry os.DirEntry) bool { return true }

	handler := func(entry os.DirEntry) error {
		srcPath := filepath.Join(md.SrcDir, entry.Name())
		dstPath := filepath.Join(md.DstDir, getZipName(entry.Name()))
		return internal.Zip(srcPath, dstPath)
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
		Strict: *strict,
	}

	mdPtr, err = internal.Prepare(mdPtr)
	if err != nil {
		return err
	}

	return DoBatch(*mdPtr)
}
