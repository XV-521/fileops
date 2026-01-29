package zip

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/internal"
	"github.com/XV-521/fileops/public"
	"os"
	"path/filepath"
)

func DoBatch(md *Mode) error {

	md, err := internal.Prepare(md)
	if err != nil {
		return err
	}

	getZipName := func(fileName string) string {
		basename, _ := internal.GetBasenameAndExt(fileName)
		return fmt.Sprintf("%v.%v", basename, "zip")
	}

	bm := internal.BatchMode{
		Sem:    10,
		Strict: md.Strict,
	}

	filter := func(entry os.DirEntry) bool { return true }

	zipFn, err := public.GetZipFn(md.ZT)
	if err != nil {
		return err
	}

	handler := func(entry os.DirEntry) error {
		srcPath := filepath.Join(md.SrcDir, entry.Name())
		dstPath := filepath.Join(md.DstDir, getZipName(entry.Name()))
		return zipFn(srcPath, dstPath, md.Pwd)
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

	zt := fs.Int(
		"zt",
		int(public.ZipUn),
		fmt.Sprintf(
			"Zip type: { %v: zip, %v: 7z, %v: tar }",
			public.ZipB, public.ZipS, public.ZipT,
		),
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
		ZT:     public.ZipType(*zt),
		Pwd:    *pwd,
		Strict: *strict,
	}

	return DoBatch(md)
}
