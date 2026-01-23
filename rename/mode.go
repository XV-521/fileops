package rename

import (
	"fmt"
	"os"
)

type Mode struct {
	SrcDir   string
	Basename string
	Namer    Namer
	Ext      string
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
	if md.Namer == nil {
		md.Namer = &NameGen{
			Basename: md.Basename,
			Ext:      md.Ext,
		}
	}
	return md, nil
}
