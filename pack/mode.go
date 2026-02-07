package pack

import (
	"fmt"
	"github.com/XV-521/fileops/core/mode"
	"os"
)

type Mode struct {
	SrcDir string
	DstDir string
	PT     mode.PackType
	Pwd    string
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

	if md.PT == mode.PackUn || md.PT == mode.PackR {
		return mode.UnsupportedPackTypeErr
	}
	if md.PT == mode.PackT && md.Pwd != "" {
		return fmt.Errorf("md.PT is ZipT (tar), but md.pwd is not empty")
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
