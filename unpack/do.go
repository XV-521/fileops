package unpack

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/core"
	"github.com/XV-521/fileops/core/mode"
	"github.com/XV-521/fileops/internal/impl"
	"github.com/XV-521/fileops/internal/util"
	"os"
	"path/filepath"
	"slices"
)

type ent struct {
	name  string
	asDir bool
}

func createUniqueName(dir string, e ent) (string, error) {
	basename := e.name
	ext := ""
	if !e.asDir {
		basename, ext = util.GetBasenameAndExt(basename)
	}

	result, err := util.IsContainTheFile(dir, e.name)
	if err != nil {
		return e.name, err
	}
	if !result {
		return e.name, nil
	}
	for i := 2; ; i++ {
		newName := fmt.Sprintf("%v(%v)%v", basename, i, ext)
		r, err := util.IsContainTheFile(dir, newName)
		if err != nil {
			return e.name, err
		}
		if !r {
			return newName, nil
		}
	}
}

func undress(dir string) error {

	dirDir := filepath.Dir(dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		uniqueName, err := createUniqueName(dirDir, ent{name, entry.IsDir()})
		if err != nil {
			return err
		}
		err = os.Rename(filepath.Join(dir, name), filepath.Join(dirDir, uniqueName))

		if err != nil {
			return err
		}
	}
	return os.Remove(dir)
}

type queueSet struct {
	queue []string
}

func (t *queueSet) append(dir string) {
	if slices.Contains(t.queue, dir) {
		return
	}
	t.queue = append(t.queue, dir)
}

func (t *queueSet) remove(dir string) {
	var targets []string
	for _, d := range t.queue {
		if d != dir {
			targets = append(targets, d)
		}
	}
	t.queue = targets
}

func DoBatch(md *Mode) error {

	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	err = util.Copy(md.SrcDir, md.DstDir)
	if err != nil {
		return err
	}

	bm := impl.BatchMode{
		Sem:    10,
		Rec:    md.Rec,
		Strict: md.Strict,
	}

	filter := func(ei core.EntryInfo) bool {
		if mode.GetPackType(ei.Name()) == mode.PackUn {
			return false
		}
		return true
	}

	dstDir := filepath.Join(md.DstDir, filepath.Base(md.SrcDir))

	qs := &queueSet{
		queue: []string{dstDir},
	}

	handler := func(ei core.EntryInfo) error {

		unzipFn, err := mode.GetUnpackFn(mode.GetPackType(ei.Name()))
		if err != nil {
			return err
		}

		path := ei.Path()

		dstName, err := createUniqueName(ei.Dir, ent{name: ei.Basename(), asDir: true})
		if err != nil {
			return err
		}

		dst := filepath.Join(ei.Dir, dstName)

		err = os.Mkdir(dst, 0777)
		if err != nil {
			return err
		}

		err = unzipFn(path, dst, md.Pwd)
		if err != nil {
			return err
		}

		entries, err := os.ReadDir(dst)
		if err != nil {
			return err
		}

		numEntries := len(entries)
		switch numEntries {
		case 0:
			err := undress(dst)
			if err != nil {
				return err
			}
		case 1:
			entry := entries[0]
			name := entry.Name()
			packType := mode.GetPackType(name)
			if packType != mode.PackUn {
				qs.append(dst)
				break
			}
			err = undress(dst)
			if err != nil {
				return err
			}
		default:
			qs.append(dst)
		}

		return os.Remove(path)
	}

	var walk func() error

	walk = func() error {

		for _, dir := range qs.queue[:] {

			err := impl.DoBatchWrapper(dir, bm, filter, handler)
			if err != nil {
				return err
			}

			qs.remove(dir)
		}

		if len(qs.queue) == 0 {
			return nil
		}

		return walk()
	}

	err = walk()
	if err != nil {
		return err
	}

	return undress(dstDir)
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
		SrcDir: *srcDir,
		DstDir: *dstDir,
		Pwd:    *pwd,
		Rec:    *rec,
		Strict: *strict,
	}

	return DoBatch(md)
}
