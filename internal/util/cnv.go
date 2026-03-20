package util

import (
	"os/exec"
)

func CnvToVideo(srcPath, dstPath string) error {
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

func CnvToAudio(srcPath, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", srcPath,
		"-vn",
		dstPath,
	)
	return CmdWrap(cmd)
}

func CnvToImage(srcPath string, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", srcPath,
		"-vframes", "1",
		dstPath,
	)
	return CmdWrap(cmd)
}

func CnvToSub(srcPath string, dstPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", srcPath,
		dstPath,
	)
	return CmdWrap(cmd)
}
