package util

import (
	"os"
	"path/filepath"
)

func RecMkdirBySrc(srcDir string, dstDir string) error {

	_, err := os.Stat(srcDir)
	if err != nil {
		return err
	}

	_, err = os.Stat(dstDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dstDir, 0777)
		if err != nil {
			return err
		}
	}

	var fn func(srcDir string, dstDir string) error

	fn = func(srcDir string, dstDir string) error {
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			path := filepath.Join(srcDir, name)
			dstPath := filepath.Join(dstDir, name)
			err := os.Mkdir(dstPath, 0777)
			if err != nil {
				return err
			}
			err = fn(path, dstPath)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return fn(srcDir, dstDir)
}

func MapTwoDir(srcDir string, dstDir string, srcPath string) (string, error) {
	rel, err := filepath.Rel(srcDir, srcPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(dstDir, rel), nil
}

func IsContainTheFile(dir string, filename string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.Name() == filename {
			return true, nil
		}
	}
	return false, nil
}
