package util

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func OffCnv(srcDir string, dstDir string, fext string, oext string) error {
	files, err := filepath.Glob(filepath.Join(srcDir, "*"+fext))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	const limit = 200

	runBatch := func(batch []string) error {
		if len(batch) == 0 {
			return nil
		}

		args := []string{
			"--headless",
			"--convert-to", strings.TrimPrefix(oext, "."),
			"--outdir", dstDir,
		}
		args = append(args, batch...)

		cmd := exec.Command("soffice", args...)
		return CmdWrap(cmd)
	}

	var subFiles []string
	for _, file := range files {
		subFiles = append(subFiles, file)
		if len(subFiles) == limit {
			if err := runBatch(subFiles); err != nil {
				return err
			}
			subFiles = nil
		}
	}

	if err := runBatch(subFiles); err != nil {
		return err
	}

	return nil
}
