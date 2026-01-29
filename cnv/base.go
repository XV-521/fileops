package cnv

import (
	"github.com/XV-521/fileops/internal"
	"os/exec"
)

func cnv(srcPath string, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-hwaccel", "videotoolbox",
		"-i", srcPath,
		"-c:v", "h264_videotoolbox",
		dstPath,
	)
	return internal.CmdWrapper(cmd)
}
