package util

import (
	"os/exec"
)

func Copy(srcPath string, dstPath string) error {
	cmd := exec.Command("cp", "-r", srcPath, dstPath)
	return CmdWrap(cmd)
}
