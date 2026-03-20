package createsub

import (
	"fmt"
	"github.com/XV-521/fileops/v2/internal/util"
	"os"
	"path/filepath"
)

type Mode struct {
	SrcDir string
	Model  string
	Lang   string
	OSrt   bool
	OVtt   bool
	OAss   bool
	OTxt   bool
	Rec    bool
	Strict bool
}

func (md *Mode) Check() error {
	if md.SrcDir == "" {
		return fmt.Errorf("md.SrcDir is empty")
	}
	_, err := os.Stat(md.SrcDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("md.srcDir %s does not exist", md.SrcDir)
	}

	if md.Lang == "" {
		return fmt.Errorf("md.Lang is empty")
	}

	if !(md.OSrt || md.OVtt || md.OAss || md.OTxt) {
		return fmt.Errorf("!(md.OSrt || md.OVtt || md.OAss || md.OTxt)")
	}

	return nil
}

func (md *Mode) Normalize() (*Mode, error) {

	_, err := os.Stat(modelDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(modelDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	if md.Model == "" {
		md.Model = "ggml-large-v3-turbo.bin"
	}

	modelPath := filepath.Join(modelDir, md.Model)
	_, err = os.Stat(modelPath)
	if os.IsNotExist(err) {
		modelUrl := baseUrl + md.Model
		err := util.Download(modelUrl, modelPath)
		if err != nil {
			return nil, err
		}
	}

	return md, nil
}
