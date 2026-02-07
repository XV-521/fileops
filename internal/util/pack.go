package util

import (
	"os/exec"
	"path/filepath"
)

func Zip(srcPath string, dstPath string, pwd string) error {
	args := []string{"a", "-tzip", "-xr!.DS_Store", "-xr!__MACOSX"}

	if pwd != "" {
		args = append(args, "-p"+pwd)
	}

	args = append(args, dstPath, filepath.Base(srcPath))

	cmd := exec.Command("7z", args...)
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

func Seven(srcPath string, dstPath string, pwd string) error {
	args := []string{"a", "-xr!.DS_Store", "-xr!__MACOSX"}

	if pwd != "" {
		args = append(args, "-p"+pwd, "-mhe=on")
	}

	args = append(args, dstPath, filepath.Base(srcPath))

	cmd := exec.Command("7z", args...)
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

func Tar(srcPath string, dstPath string, _ string) error {
	args := []string{
		"-czf", dstPath,
		"--exclude=.DS_Store",
		"--exclude=__MACOSX",
		filepath.Base(srcPath),
	}
	cmd := exec.Command("tar", args...)
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

// Deprecated: use SevenUnzip instead.
func Unzip(srcPath string, dstDir string, pwd string) error {
	var args []string
	if pwd != "" {
		args = append(args, "-P", pwd)
	}
	args = append(args, srcPath, "-d", dstDir)
	cmd := exec.Command("unzip", args...)
	return CmdWrapper(cmd)
}

func UnSeven(srcPath string, dstDir string, pwd string) error {
	args := []string{"x", srcPath, "-o" + dstDir}
	if pwd != "" {
		args = append(args, "-p"+pwd)
	}
	cmd := exec.Command("7z", args...)
	return CmdWrapper(cmd)
}

func UnTar(srcPath string, dstDir string, _ string) error {
	args := []string{"-xf", srcPath, "-C", dstDir}
	cmd := exec.Command("tar", args...)
	return CmdWrapper(cmd)
}

func UnRar(srcPath string, dstDir string, pwd string) error {
	args := []string{srcPath, "-o", dstDir}
	if pwd != "" {
		args = append(args, "-p", pwd)
	}

	cmd := exec.Command("unar", args...)
	return CmdWrapper(cmd)
}
