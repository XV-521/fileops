package createsub

import (
	"github.com/XV-521/fileops/v2/internal/util"
	"os/exec"
)

func whisper(
	audioPath string,
	modelPath string,
	lang string,
	oSrt bool,
	oVtt bool,
	oTxt bool,
) error {
	args := []string{
		"-m", modelPath,
		"-f", audioPath,
		"-l", lang,
	}

	if oSrt {
		args = append(args, "-osrt")
	}
	if oVtt {
		args = append(args, "-ovtt")
	}
	if oTxt {
		args = append(args, "-otxt")
	}

	cmd := exec.Command("whisper-cli", args...)
	return util.CmdWrap(cmd)
}
