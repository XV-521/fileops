package zip

import (
	"fmt"
	"github.com/XV-521/fileops/public"
	"os"
)

type Mode struct {
	SrcDir string
	DstDir string
	ZT     public.ZipType
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

	if md.ZT == public.ZipUn || md.ZT == public.ZipR {
		return public.UnsupportedZipTypeErr
	}
	if md.ZT == public.ZipT && md.Pwd != "" {
		return fmt.Errorf("md.ZT is ZipT (tar), but md.pwd is not empty")
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
