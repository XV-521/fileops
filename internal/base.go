package internal

import (
	"os/exec"
	"path/filepath"
	"strings"
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

func Unzip(srcPath string, dstDir string, pwd string) error {
	var args []string
	if pwd != "" {
		args = append(args, "-P", pwd)
	}
	args = append(args, srcPath, "-d", dstDir)
	cmd := exec.Command("unzip", args...)
	return CmdWrapper(cmd)
}

func Cnv(srcPath string, dstPath string) error {
	cmd := exec.Command("ffmpeg", "-i", srcPath, dstPath)
	return CmdWrapper(cmd)
}

func IsThisExt(filename string, ext string) bool {
	result := strings.Split(filename, ".")
	if len(result) < 2 {
		return ext == ""
	}
	return strings.ToLower(result[len(result)-1]) == strings.ToLower(ext)
}

func GetBasenameAndExt(filename string) (baseName string, ext string) {
	result := strings.Split(filename, ".")
	length := len(result)
	if length < 2 {
		return filename, ""
	}
	return strings.Join(result[0:length-1], "."), result[length-1]
}
