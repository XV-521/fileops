package mapdir

import (
	"errors"
	"flag"
	"github.com/XV-521/fileops/v2/internal/impl"
	"github.com/XV-521/fileops/v2/internal/util"
)

func Do(md *Mode) error {
	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}
	return util.RecMkdirBySrc(md.SrcDir, md.DstDir)
}

func DoWithFlags(fs *flag.FlagSet, args []string) error {
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
	}

	return Do(md)
}
