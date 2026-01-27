package cnv

import (
	"fmt"
	"os"
)

type Mode struct {
	SrcDir  string
	DstDir  string
	FromExt string
	ToExt   string
	Strict  bool
}

func (md *Mode) Check() error {
	if md.SrcDir == "" {
		return fmt.Errorf("md.SrcDir is empty")
	}
	_, err := os.Stat(md.SrcDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("md.srcDir %s does not exist", md.SrcDir)
	}
	if md.DstDir == "" {
		return fmt.Errorf("md.DstDir is empty")
	}
	if md.FromExt == "" {
		return fmt.Errorf("md.FromExt is empty")
	}
	if md.ToExt == "" {
		return fmt.Errorf("md.ToExt is empty")
	}
	return nil
}

func (md *Mode) Normalize() (*Mode, error) {
	_, err := os.Stat(md.DstDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(md.DstDir, 0777)
		if err != nil {
			return nil, err
		}
	}
	return md, nil
}
