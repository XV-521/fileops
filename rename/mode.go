package rename

import (
	"fmt"
	"os"
)

type Mode struct {
	SrcDir   string
	Basename string
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

	if md.Basename == "" {
		return fmt.Errorf("md.Basename is empty")
	}

	return nil
}

func (md *Mode) Normalize() (*Mode, error) {
	return md, nil
}
