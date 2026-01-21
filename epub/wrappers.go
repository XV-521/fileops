package epub

import (
	"github.com/XV-521/fileops/internal"
	"os"
	"os/exec"
)

func packEpub(srcDir string) error {
	cmd0 := exec.Command("zip", "-X0", "./new.epub", "mimetype")
	cmd0.Dir = srcDir
	if err := internal.CmdWrapper(cmd0); err != nil {
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

	err = internal.CmdWrapper(cmd1)
	if err != nil {
		return err
	}

	return nil
}

func epubWrapper(srcPath string, dstDir string, handler func(srcDir string) error) error {

	wrapper := func(dir string) error {
		err := internal.Unzip(srcPath, dir)
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
		return nil
	}

	err := internal.MkdirTempWrapper(dstDir, wrapper)
	if err != nil {
		return err
	}

	return nil
}
