package pack

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/core/mode"
	"github.com/XV-521/fileops/internal/impl"
	"github.com/XV-521/fileops/internal/util"
	"path/filepath"

	"github.com/XV-521/fileops/core"
)

func DoBatch(md *Mode) error {

	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	ext, err := mode.CreatePackExt(md.PT)

	getZipName := func(fileName string) string {
		basename, _ := util.GetBasenameAndExt(fileName)
		return fmt.Sprintf("%v%v", basename, ext)
	}

	bm := impl.BatchMode{
		Sem:    10,
		Rec:    false,
		Strict: md.Strict,
	}

	filter := func(ei core.EntryInfo) bool { return ei.Name() != ".DS_Store" }

	zipFn, err := mode.GetPackFn(md.PT)
	if err != nil {
		return err
	}

	handler := func(ei core.EntryInfo) error {
		dstPath := filepath.Join(md.DstDir, getZipName(ei.Name()))
		return zipFn(ei.Path(), dstPath, md.Pwd)
	}

	return impl.DoBatchWrapper(md.SrcDir, bm, filter, handler)
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

	pt := fs.Int(
		"pt",
		int(mode.PackUn),
		fmt.Sprintf(
			"Pack type: { %v: zip, %v: 7z, %v: tar }",
			mode.PackZ, mode.PackS, mode.PackT,
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
		PT:     mode.PackType(*pt),
		Pwd:    *pwd,
		Strict: *strict,
	}

	return DoBatch(md)
}
