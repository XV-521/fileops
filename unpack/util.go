package unpack

import (
	"fmt"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"path/filepath"
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
