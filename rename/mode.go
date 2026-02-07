package rename

import (
	"fmt"
	"github.com/XV-521/fileops/v2/core"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
)

type Mode struct {
	SrcDir   string
	Basename string
	Namer    core.Namer
	Ext      string
	Rec      bool
	Strict   bool
}

func (md *Mode) Check() error {
	if md.SrcDir == "" {
		return fmt.Errorf("md.SrcDir is empty")
	}

	_, err := os.Stat(md.SrcDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("md.srcDir %s does not exist", md.SrcDir)
	}

	if md.Basename == "" && md.Namer == nil {
		return fmt.Errorf("md.Basename is empty and md.Namer is nil")
	}

	return nil
}

func (md *Mode) Normalize() (*Mode, error) {
	md.Ext = util.GetClearExt(md.Ext)
	if md.Namer == nil {
		md.Namer = &nameGen{
			basename: md.Basename,
			ext:      md.Ext,
		}
	}
	return md, nil
}
