package core

import (
	"fmt"
	"github.com/XV-521/fileops/internal/util"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type entryView interface {
	Name() string
	IsDir() bool
	Info() (os.FileInfo, error)
}

type FileInfoView struct {
	Fi os.FileInfo
}

func (v FileInfoView) Name() string {
	return v.Fi.Name()
}

func (v FileInfoView) IsDir() bool {
	return v.Fi.IsDir()
}

func (v FileInfoView) Info() (os.FileInfo, error) {
	return v.Fi, nil
}

type EntryInfo struct {
	Dir  string
	View entryView
}

func (e EntryInfo) Path() string {
	return filepath.Join(e.Dir, e.View.Name())
}

func (e EntryInfo) Name() string {
	return e.View.Name()
}

func (e EntryInfo) IsDir() bool {
	return e.View.IsDir()
}

func (e EntryInfo) Basename() string {
	basename, _ := util.GetBasenameAndExt(e.View.Name())
	return basename
}

func (e EntryInfo) Ext() string {
	_, ext := util.GetBasenameAndExt(e.View.Name())
	return ext
}

func (e EntryInfo) Size() (int64, error) {
	info, err := e.View.Info()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (e EntryInfo) ModTime() (time.Time, error) {
	info, err := e.View.Info()
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func (e EntryInfo) CrtTime() (time.Time, error) {
	info, err := e.View.Info()
	if err != nil {
		return time.Time{}, err
	}

	stat, ok := info.Sys().(*syscall.Stat_t)

	if ok {
		return time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec), nil
	}
	return time.Time{}, fmt.Errorf("birth time not supported on this OS")
}
