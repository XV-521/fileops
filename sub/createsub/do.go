package createsub

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
)

type group struct {
	videoPath      string
	relatedTempDir string
}

func (g *group) normalize(tempDir string) error {
	if g.relatedTempDir != "" {
		return nil
	}
	for {
		relatedTempDir := filepath.Join(tempDir, util.GetRand(12))
		_, err := os.Stat(relatedTempDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(relatedTempDir, 0777)
			if err != nil {
				return err
			}
			g.relatedTempDir = relatedTempDir
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (g *group) videoFilename() string {
	return filepath.Base(g.videoPath)
}

func (g *group) relatedTempFilepath(ext string) string {
	basename, _ := util.GetBasenameAndExt(g.videoFilename())
	relatedFilename := fmt.Sprintf("%v%v", basename, util.GetClearExt(ext))
	return filepath.Join(g.relatedTempDir, relatedFilename)
}

func doBatchHandle(md *Mode, tempDir string) error {

	var groups []group

	bm := impl.BatchMode{
		Sem:    1,
		Rec:    md.Rec,
		Strict: md.Strict,
	}

	filter := func(ei core.EntryInfo) bool {
		if ei.IsDir() {
			return false
		}
		if mode.GetCnvType(ei.Name()) != mode.CnvV {
			return false
		}
		return true
	}

	handler := func(ei core.EntryInfo) error {

		g := group{videoPath: ei.Path()}
		err := g.normalize(tempDir)
		if err != nil {
			return err
		}

		tempWavPath := g.relatedTempFilepath(".wav")
		err = util.CnvToAudio(g.videoPath, tempWavPath)
		if err != nil {
			return err
		}
		err = whisper(
			tempWavPath,
			filepath.Join(modelDir, md.Model),
			md.Lang,
			md.OSrt || md.OAss,
			md.OVtt,
			md.OTxt,
		)
		if err != nil {
			return err
		}

		if md.OAss {
			tempSrtPath := tempWavPath + ".srt"
			tempAssPath := tempWavPath + ".ass"
			err = util.CnvToSub(tempSrtPath, tempAssPath)
			if err != nil {
				return err
			}
			if !md.OSrt {
				err = os.Remove(tempSrtPath)
				if err != nil {
					return err
				}
			}
		}

		err = os.Remove(tempWavPath)
		if err != nil {
			return err
		}

		groups = append(groups, g)
		return nil
	}

	err := impl.DoBatchWrap(md.SrcDir, bm, filter, handler)
	if err != nil {
		if md.Strict {
			return err
		}
	}
	for _, g := range groups {

		var dstDir string

		creator := func(uniquePath string) error {
			dstDir = uniquePath
			return os.Rename(g.relatedTempDir, uniquePath)
		}

		basename, _ := util.GetBasenameAndExt(g.videoFilename())

		err := util.CreateUniqueDir(filepath.Dir(g.videoPath), basename, creator)
		if err != nil {
			return err
		}

		videoDstPath := filepath.Join(dstDir, g.videoFilename())
		err = os.Rename(g.videoPath, videoDstPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func DoBatch(md *Mode) error {
	md, err := impl.Prepare(md)
	if err != nil {
		return err
	}

	handle := func(dir string) error {
		return doBatchHandle(md, dir)
	}

	return util.MkdirTempWrap(handle)
}

func DoBatchWithFlags(fs *flag.FlagSet, args []string) error {
	srcDir := fs.String(
		"src",
		"",
		"Source directory.",
	)

	model := fs.String(
		"model",
		"ggml-large-v3-turbo.bin",
		"Model filename (e.g., ggml-large-v3-turbo.bin).",
	)

	oSrt := fs.Bool(
		"osrt",
		false,
		"Output result in a srt file.",
	)

	oVtt := fs.Bool(
		"ovtt",
		false,
		"Output result in a vtt file.",
	)

	oAss := fs.Bool(
		"oass",
		false,
		"Output result in a ass file.",
	)

	oTxt := fs.Bool(
		"otxt",
		false,
		"Output result in a txt file.",
	)

	lang := fs.String(
		"lang",
		"en",
		"Language.",
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
		Model:  *model,
		Lang:   *lang,
		OSrt:   *oSrt,
		OVtt:   *oVtt,
		OAss:   *oAss,
		OTxt:   *oTxt,
		Rec:    *rec,
		Strict: *strict,
	}

	return DoBatch(md)
}
