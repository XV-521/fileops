package addsub

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/v2/core"
	"github.com/XV-521/fileops/v2/core/mode"
	"github.com/XV-521/fileops/v2/internal/impl"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"path/filepath"
	"sync"
)

func DoBatch(md *Mode) error {

	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	subFn, err := getSubFn(md.ST)
	if err != nil {
		return err
	}

	var tasks []func() error
	var mu sync.Mutex

	bm := impl.BatchMode{
		Sem:    6,
		Rec:    md.Rec,
		Strict: md.Strict,
	}

	filter := func(ei core.EntryInfo) bool {
		return mode.GetCnvType(ei.Name()) == mode.CnvV
	}

	handler := func(ei core.EntryInfo) error {
		entries, err := os.ReadDir(ei.Dir)
		if err != nil {
			return err
		}
		var subPaths []string

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			entryName := entry.Name()
			if md.IsSub(entryName) {
				subPaths = append(subPaths, filepath.Join(ei.Dir, entryName))
			}
			if mode.GetCnvType(entryName) == mode.CnvV && entryName != ei.Name() {
				return fmt.Errorf("%v: there is more than 1 video", ei.Dir)
			}
		}
		if len(subPaths) != 1 {
			return fmt.Errorf("%v: expected 1 sub entry, found %v", ei.Dir, len(subPaths))
		}
		subPath := subPaths[0]

		if util.IsThisExt(filepath.Base(subPath), ".ass") {
			err = md.Style.insertToAss(subPath)
			if err != nil {
				return err
			}
		}

		newVideoName := ei.Name()
		if md.SrcDir == md.DstDir {
			videoBasename, videoExt := util.GetBasenameAndExt(newVideoName)
			newVideoName = fmt.Sprintf("%v_new%v", videoBasename, videoExt)
		}

		dstDir, err := util.MapTwoDir(md.SrcDir, md.DstDir, ei.Dir)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, newVideoName)

		err = subFn(ei.Path(), subPath, dstPath)
		if err != nil {
			return err
		}

		if md.SrcDir != md.DstDir {
			var task func() error
			if md.MoveSub {
				task = func() error {
					return os.Rename(subPath, filepath.Join(dstDir, filepath.Base(subPath)))
				}
			} else {
				task = func() error {
					return util.Undress(dstDir)
				}
			}
			mu.Lock()
			tasks = append(tasks, task)
			mu.Unlock()
		}

		return nil
	}

	err = impl.DoBatchWrap(md.SrcDir, bm, filter, handler)
	if err != nil && md.Strict {
		return err
	}
	for _, task := range tasks {
		err := task()
		if err != nil && md.Strict {
			return err
		}
	}

	return err
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

	subExt := fs.String(
		"subext",
		"",
		"Subtitle file extension",
	)

	st := fs.Int(
		"st",
		int(SubUn),
		fmt.Sprintf("Sub type: { %v: soft sub, %v: hard sub}.", SubS, SubH),
	)

	style := fs.String(
		"style",
		"Fontname=Arial,Fontsize=16,PrimaryColour=&H00ffffff,OutlineColour=&H0,BackColour=&H0,BorderStyle=1,Outline=1,Shadow=0,Alignment=2,MarginV=10",
		"Style of additional substitutions.",
	)

	moveSub := fs.Bool(
		"msub",
		false,
		"Move sub file to destination directory.",
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
		SubExt:  *subExt,
		ST:      SubType(*st),
		Style:   SimpleStyle(*style),
		MoveSub: *moveSub,
		Rec:     *rec,
		Strict:  *strict,
	}

	return DoBatch(md)
}
