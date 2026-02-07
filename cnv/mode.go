package cnv

import (
	"fmt"
	"github.com/XV-521/fileops/v2/core/mode"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
)

type Mode struct {
	SrcDir  string
	DstDir  string
	CT      mode.CnvType
	FromExt string
	ToExt   string
	Rec     bool
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
	if md.FromExt == "" && md.CT == mode.CnvUn {
		return fmt.Errorf("md.FromExt is empty and md.CT is CnvUn")
	}
	if md.FromExt != "" && md.CT != mode.CnvUn {
		return fmt.Errorf("md.FromExt is not empty and md.CT is not CnvUn")
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
	if md.Rec && (md.SrcDir != md.DstDir) {
		err := util.RecMkdirBySrc(md.SrcDir, md.DstDir)
		if err != nil {
			return nil, err
		}
	}
	md.FromExt = util.GetClearExt(md.FromExt)
	md.ToExt = util.GetClearExt(md.ToExt)
	return md, nil
}
