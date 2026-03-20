package util

import "os/exec"

func Download(url string, destPath string) error {
	cmd := exec.Command(
		"aria2c",
		"-x", "16",
		"-s", "16",
		"-o", destPath,
		url,
	)
	return CmdWrap(cmd)
}
