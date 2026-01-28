package cnv

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/internal"
	"os"
	"path/filepath"
	"strings"
)

func DoBatch(md *Mode) error {
	md, err := internal.Prepare(md)
	if err != nil {
		return err
	}

	getNewName := func(filename string) string {
		basename, _ := internal.GetBasenameAndExt(filename)
		return fmt.Sprintf("%v.%v", basename, strings.Trim(md.ToExt, "."))
	}

	bm := internal.BatchMode{
		Sem:    6,
		Strict: md.Strict,
	}

	filter := func(entry os.DirEntry) bool {
		if entry.IsDir() {
			return false
		}
		if !internal.IsThisExt(entry.Name(), md.FromExt) {
			return false
		}
		return true
	}

	handler := func(entry os.DirEntry) error {
		srcPath := filepath.Join(md.SrcDir, entry.Name())
		dstPath := filepath.Join(md.DstDir, getNewName(entry.Name()))
		return internal.Cnv(srcPath, dstPath)
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

	fromExt := fs.String(
		"fext",
		"",
		"Source file extension.",
	)

	toExt := fs.String(
		"oext",
		"",
		"Output file extension.",
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
		SrcDir:  *srcDir,
		DstDir:  *dstDir,
		FromExt: *fromExt,
		ToExt:   *toExt,
		Strict:  *strict,
	}

	return DoBatch(md)
}
