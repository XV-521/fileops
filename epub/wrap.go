package epub

import (
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"os/exec"
)

func packEpub(srcDir string) error {
	cmd0 := exec.Command("zip", "-X0", "./new.epub", "mimetype")
	cmd0.Dir = srcDir
	if err := util.CmdWrap(cmd0); err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	args := []string{"-r", "./new.epub"}
	for _, e := range entries {
		if e.Name() == "mimetype" {
			continue
		}
		args = append(args, e.Name())
	}

	cmd1 := exec.Command("zip", args...)
	cmd1.Dir = srcDir

	err = util.CmdWrap(cmd1)
	if err != nil {
		return err
	}

	return nil
}

func epubWrap(srcPath string, dstDir string, handler func(srcDir string) error) error {

	wrapper := func(dir string) error {
		err := util.UnSeven(srcPath, dir, "")
		if err != nil {
			return err
		}
		err = handler(dir)
		if err != nil {
			return err
		}
		err = packEpub(dir)
		if err != nil {
			return err
		}
		return os.Rename(dir, dstDir)
	}

	err := util.MkdirTempWrap(wrapper)
	if err != nil {
		return err
	}

	return nil
}
