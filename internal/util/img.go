package util

import (
	"fmt"
	"os/exec"
)

func Resize(srcPath string, dstPath string, rto float64) error {
	cmd := exec.Command("vips", "resize", srcPath, dstPath, fmt.Sprintf("%f", rto))
	return CmdWrapper(cmd)
}

func ChangeDpi(srcPath string, dstPath string, dpi float64) error {
	dpmStr := fmt.Sprintf("%.4f", dpi/25.4)
	cmd := exec.Command("vips", "copy", srcPath, dstPath, "--xres", dpmStr, "--yres", dpmStr)
	return CmdWrapper(cmd)
}
