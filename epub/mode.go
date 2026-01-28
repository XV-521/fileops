package epub

import (
	"fmt"
	"github.com/XV-521/fileops/internal"
	"path/filepath"
)

type Mode struct {
	SrcPath        string
	DstDir         string
	Strict         bool
	StrictCodeOnly bool
	Tag            string
	Lang           string
	Style          string
	BgColor        string
	CssBasename    string
}

func (md *Mode) Check() error {
	if md.SrcPath == "" {
		return fmt.Errorf("srcPath is empty")
	}
	if md.Lang == "" {
		return fmt.Errorf("lang is empty")
	}
	if !md.StrictCodeOnly && md.Tag == "" {
		return fmt.Errorf("the StrictCodeOnly is false, and tag is empty")
	}
	return nil
}

func (md *Mode) Normalize() (*Mode, error) {

	if md.DstDir == "" {
		dir := filepath.Dir(md.SrcPath)
		name := "inner" + internal.GetRand(6)
		md.DstDir = filepath.Join(dir, name)
	}

	if md.Style == "" {
		md.Style = "tango"
	}

	if md.CssBasename == "" {
		name := "highlight" + internal.GetRand(6)
		md.CssBasename = fmt.Sprintf("%v%v", name, ".css")
	}

	return md, nil
}
