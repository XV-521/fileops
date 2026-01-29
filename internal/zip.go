package internal

import (
	"os/exec"
	"path/filepath"
)

func Zip(srcPath string, dstPath string, pwd string) error {
	args := []string{"-r", "-X"}
	if pwd != "" {
		args = append(args, "-P", pwd)
	}
	args = append(args, dstPath, filepath.Base(srcPath))
	cmd := exec.Command("zip", args...)
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

func SevenZip(srcPath string, dstPath string, pwd string) error {
	args := []string{"a", "-xr!.DS_Store", "-xr!__MACOSX"}

	if pwd != "" {
		args = append(args, "-p"+pwd, "-mhe=on")
	}

	args = append(args, dstPath, filepath.Base(srcPath))

	cmd := exec.Command("7z", args...)
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

func TarZip(srcPath string, dstPath string, pwd string) error {
	_ = pwd
	args := []string{
		"-cz",
		"--exclude=.DS_Store",
		"--exclude=__MACOSX",
		"--no-xattrs",
		"--no-acls",
		"--no-selinux",
		"-f", dstPath,
		filepath.Base(srcPath),
	}
	cmd := exec.Command("tar", args...)
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

func Unzip(srcPath string, dstDir string, pwd string) error {
	var args []string
	if pwd != "" {
		args = append(args, "-P", pwd)
	}
	args = append(args, srcPath, "-d", dstDir)
	cmd := exec.Command("unzip", args...)
	return CmdWrapper(cmd)
}

func SevenUnzip(srcPath string, dstDir string, pwd string) error {
	args := []string{"x", srcPath, "-o" + dstDir}
	if pwd != "" {
		args = append(args, "-p"+pwd)
	}
	cmd := exec.Command("7z", args...)
	return CmdWrapper(cmd)
}

func TarUnzip(srcPath string, dstDir string, _ string) error {
	args := []string{"-xf", srcPath, "-C", dstDir}
	cmd := exec.Command("tar", args...)
	return CmdWrapper(cmd)
}

func RarUnzip(srcPath string, dstDir string, pwd string) error {
	args := []string{srcPath, "-o", dstDir}
	if pwd != "" {
		args = append(args, "-p", pwd)
	}

	cmd := exec.Command("unar", args...)
	return CmdWrapper(cmd)
}
