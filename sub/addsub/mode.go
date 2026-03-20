package addsub

import (
	"fmt"
	"github.com/XV-521/fileops/v2/core/out"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"slices"
	"strings"
)

type Mode struct {
	SrcDir  string
	DstDir  string
	SubExt  string
	IsSub   func(filename string) bool
	ST      SubType
	Style   SimpleStyle
	MoveSub bool
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
		return fmt.Errorf("md.dstDir is empty")
	}

	if md.DstDir == md.SrcDir && md.MoveSub {
		return fmt.Errorf("md.DstDir == md.SrcDir but md.MoveSub is true")
	}

	if md.ST != SubH && md.ST != SubS {
		return fmt.Errorf("md.st %v is invalid", md.ST)
	}

	if md.ST == SubS {
		out.Warn.Println("Note: When selecting soft sub, styles are supported only in MKV.")
	}

	subExt := strings.ToLower(util.GetClearExt(md.SubExt))
	if !slices.Contains([]string{".srt", ".ass", ".vtt"}, subExt) {
		return fmt.Errorf("sub.ext %s is invalid", md.SubExt)
	}

	if subExt == ".ass" {
		ok := md.Style.check()
		if !ok {
			return fmt.Errorf("md.style %s is invalid", md.Style)
		}
	} else {
		out.Warn.Println("Note: Style take effect only when md.SubExt == '.ass'.")
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
	md.SubExt = util.GetClearExt(md.SubExt)
	if md.IsSub == nil {
		md.IsSub = func(filename string) bool {
			return util.IsThisExt(filename, md.SubExt)
		}
	}
	return md, nil
}
