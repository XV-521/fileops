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

func DoBatch(md Mode) error {
	getNewName := func(filename string) string {
		basename, _ := internal.GetBasenameAndExt(filename)
		return fmt.Sprintf("%v.%v", basename, strings.Trim(md.ToExt, "."))
	}

	bm := internal.BatchMode{
		Async:  true,
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
		iptPath := filepath.Join(md.SrcDir, entry.Name())
		optPath := filepath.Join(md.DstDir, getNewName(entry.Name()))
		return internal.Cnv(iptPath, optPath)
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

	mdPtr := &Mode{
		SrcDir:  *srcDir,
		DstDir:  *dstDir,
		FromExt: *fromExt,
		ToExt:   *toExt,
		Strict:  *strict,
	}

	mdPtr, err = internal.Prepare(mdPtr)
	if err != nil {
		return err
	}

	return DoBatch(*mdPtr)
}
