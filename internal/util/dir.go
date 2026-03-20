package util

import (
	"fmt"
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

var dirLocks = NewLocks()

func CreateUniqueFile(
	dir string,
	filename string,
	creator func(uniquePath string) error,
) error {
	mu := dirLocks.GetOrCreate(dir)
	mu.Lock()
	defer mu.Unlock()
	basename, ext := GetBasenameAndExt(filename)
	uniquePath := filepath.Join(dir, fmt.Sprintf("%v%v", basename, ext))
	for i := 2; ; i++ {
		_, err := os.Stat(uniquePath)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return err
		}
		uniquePath = filepath.Join(dir, fmt.Sprintf("%v(%v)%v", basename, i, ext))
	}
	return creator(uniquePath)
}

func CreateUniqueDir(
	dir string,
	dirname string,
	creator func(uniquePath string) error,
) error {
	mu := dirLocks.GetOrCreate(dir)
	mu.Lock()
	defer mu.Unlock()
	uniquePath := filepath.Join(dir, dirname)
	for i := 2; ; i++ {
		_, err := os.Stat(uniquePath)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return err
		}
		uniquePath = filepath.Join(dir, fmt.Sprintf("%v(%v)", dirname, i))
	}
	return creator(uniquePath)
}

func Undress(dir string) error {
	dirDir := filepath.Dir(dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		var fn func(dir string, name string, creator func(uniquePath string) error) error
		if entry.IsDir() {
			fn = CreateUniqueDir
		} else {
			fn = CreateUniqueFile
		}

		creator := func(uniquePath string) error {
			return os.Rename(filepath.Join(dir, name), uniquePath)
		}

		err = fn(dirDir, name, creator)
		if err != nil {
			return err
		}
	}
	return os.Remove(dir)
}
