package offcnv

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/v2/core/out"
	"github.com/XV-521/fileops/v2/internal/impl"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"path/filepath"
	"strings"
)

func DoBatch(md *Mode) error {

	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	var tasks []func() error

	handle := func(srcDir string) error {
		dstDir, err := util.MapTwoDir(md.SrcDir, md.DstDir, srcDir)
		if err != nil {
			return err
		}
		return util.OffCnv(srcDir, dstDir, md.FromExt, md.ToExt)
	}

	check := func(srcDir string) error {
		var missing []string
		var skipped []string

		entries, err := os.ReadDir(srcDir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			path := filepath.Join(srcDir, entry.Name())

			mappedPath, err := util.MapTwoDir(md.SrcDir, md.DstDir, path)
			if err != nil {
				return err
			}
			mappedPath = strings.TrimSuffix(mappedPath, md.FromExt) + md.ToExt

			_, err = os.Stat(mappedPath)
			if err == nil {
				continue
			}
			if !os.IsNotExist(err) {
				return err
			}

			if util.IsThisExt(entry.Name(), md.FromExt) {
				missing = append(missing, entry.Name())
			} else {
				skipped = append(skipped, entry.Name())
			}
		}

		if len(missing) > 0 || len(skipped) > 0 {
			out.Norm.Printf("%v:\n", srcDir)

			if len(missing) > 0 {
				out.Warn.Println("missing:")
				for _, name := range missing {
					out.Warn.Println(name)
				}
			}

			if len(skipped) > 0 {
				out.Norm.Println("skipped:")
				for _, name := range skipped {
					out.Norm.Println(name)
				}
			}
			fmt.Print("\n")
		}

		if len(missing) > 0 {
			return fmt.Errorf("partial conversion failed")
		}

		return nil
	}

	task := func(srcDir string) error {
		err := handle(srcDir)
		if err != nil {
			return fmt.Errorf("handle failed %v: %v", srcDir, err)
		}
		err = check(srcDir)
		if err != nil {
			return fmt.Errorf("check failed %v: %v", srcDir, err)
		}
		return nil
	}

	tasks = append(
		tasks,
		func() error { return task(md.SrcDir) },
	)

	var walk func(srcDir string) error
	walk = func(srcDir string) error {

		entries, err := os.ReadDir(srcDir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {

				path := filepath.Join(srcDir, entry.Name())

				tasks = append(
					tasks,
					func() error { return task(path) },
				)

				err := walk(path)
				if err != nil {
					if md.Strict {
						return err
					}
				}
			}
		}
		return nil
	}

	err = walk(md.SrcDir)
	if err != nil {
		return err
	}

	var firstErr error
	for _, t := range tasks {
		err := t()
		if err != nil {
			if md.Strict {
				return err
			}
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
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
		FromExt: *fromExt,
		Rec:     *rec,
		ToExt:   *toExt,
		Strict:  *strict,
	}

	return DoBatch(md)
}
