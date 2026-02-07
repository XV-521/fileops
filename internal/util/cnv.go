package util

import (
	"os/exec"
)

func CnvForVideo(srcPath, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-hwaccel", "videotoolbox",
		"-i", srcPath,
		"-c:v", "h264_videotoolbox",
		"-c:a", "aac",
		"-movflags", "+faststart",
		dstPath,
	)
	return CmdWrap(cmd)
}

func CnvForAudio(srcPath, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", srcPath,
		"-vn",
		"-c:a", "aac",
		dstPath,
	)
	return CmdWrap(cmd)
}

func CnvForImage(srcPath string, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", srcPath,
		dstPath,
	)
	return CmdWrap(cmd)
}
