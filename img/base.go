package img

import (
	"fmt"
	"github.com/XV-521/fileops/internal"
	"os/exec"
)

func resize(srcPath string, dstPath string, rto float64) error {
	cmd := exec.Command("vips", "resize", srcPath, dstPath, fmt.Sprintf("%f", rto))
	return internal.CmdWrapper(cmd)
}

func changeDpi(srcPath string, dstPath string, dpi float64) error {
	dpmStr := fmt.Sprintf("%.4f", dpi/25.4)
	cmd := exec.Command("vips", "copy", srcPath, dstPath, "--xres", dpmStr, "--yres", dpmStr)
	return internal.CmdWrapper(cmd)
}
