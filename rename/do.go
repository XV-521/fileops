package rename

import (
	"errors"
	"flag"
	"github.com/XV-521/fileops/v2/core"
	"github.com/XV-521/fileops/v2/internal/impl"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"path/filepath"
)

func DoBatch(md *Mode) error {

	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	bm := impl.BatchMode{
		Sem:    1,
		Rec:    md.Rec,
		Strict: md.Strict,
	}

	filter := func(ei core.EntryInfo) bool {
		if ei.IsDir() {
			return false
		}
		if !util.IsThisExt(ei.Name(), md.Ext) {
			return false
		}
		return true
	}

	handler := func(ei core.EntryInfo) error {
		return os.Rename(
			ei.Path(),
			filepath.Join(ei.Dir, md.Namer.Next(ei)),
		)
	}

	return impl.DoBatchWrapper(md.SrcDir, bm, filter, handler)
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
		SrcDir:   *srcDir,
		Basename: *basename,
		Ext:      *ext,
		Rec:      *rec,
		Strict:   *strict,
	}

	return DoBatch(md)
}
