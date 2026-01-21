package epub

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
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
		basename := "inner"
		for range 6 {
			num := rand.Intn(10)
			basename += strconv.Itoa(num)
		}
		md.DstDir = filepath.Join(dir, basename)
	}

	if md.Style == "" {
		md.Style = "tango"
	}

	if md.CssBasename == "" {
		basename := "highlight"
		for range 6 {
			num := rand.Intn(10)
			basename += strconv.Itoa(num)
		}
		md.CssBasename = fmt.Sprintf("%v%v", basename, ".css")
	}

	return md, nil
}
