package unpack

import (
	"fmt"
	"os"
)

type Mode struct {
	SrcDir string
	DstDir string
	Pwd    string
	Rec    bool
	Strict bool
}

func (md *Mode) Check() error {
	if md.SrcDir == "" {
		return fmt.Errorf("md.srcDir is empty")
	}
	_, err := os.Stat(md.SrcDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("md.srcDir %s does not exist", md.SrcDir)
	}
	if md.DstDir == "" {
		return fmt.Errorf("md.dstDir is empty")
	}
	if md.DstDir == md.SrcDir {
		return fmt.Errorf("md.DstDir is the same as md.SrcDir")
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
