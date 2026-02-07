package cnv

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/v2/core"
	"github.com/XV-521/fileops/v2/core/mode"
	"github.com/XV-521/fileops/v2/internal/impl"
	"github.com/XV-521/fileops/v2/internal/util"
	"path/filepath"
)

func DoBatch(md *Mode) error {
	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	getNewName := func(filename string) string {
		basename, _ := util.GetBasenameAndExt(filename)
		return fmt.Sprintf("%v%v", basename, md.ToExt)
	}

	bm := impl.BatchMode{
		Sem:    6,
		Rec:    md.Rec,
		Strict: md.Strict,
	}

	filter := func(ei core.EntryInfo) bool {
		if ei.IsDir() {
			return false
		}
		name := ei.Name()
		if md.CT != mode.CnvUn && mode.GetCnvType(name) != md.CT {
			return false
		}
		if md.FromExt != "" && !util.IsThisExt(name, md.FromExt) {
			return false
		}
		return true
	}

	handler := func(ei core.EntryInfo) error {
		name := ei.Name()

		dstDir, err := util.MapTwoDir(md.SrcDir, md.DstDir, ei.Dir)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, getNewName(name))

		var ct mode.CnvType

		if md.CT != mode.CnvUn {
			ct = md.CT
		} else {
			ct = mode.GetCnvType(name)
		}

		cvFn, err := mode.GetCnvFn(ct)
		if err != nil {
			return err
		}
		return cvFn(ei.Path(), dstPath)
	}

	return impl.DoBatchWrap(md.SrcDir, bm, filter, handler)
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
	ct := fs.Int(
		"ct",
		int(mode.CnvUn),
		fmt.Sprintf("Zip type: { %v: video, %v: audio, %v: image }", mode.CnvV, mode.CnvA, mode.CnvI),
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

	rec := fs.Bool(
		"rec",
		false,
		"Recursive.",
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
		CT:      mode.CnvType(*ct),
		FromExt: *fromExt,
		Rec:     *rec,
		ToExt:   *toExt,
		Strict:  *strict,
	}

	return DoBatch(md)
}
