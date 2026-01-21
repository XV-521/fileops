package internal

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func Zip(srcPath string, dstPath string) error {
	cmd := exec.Command("zip", "-r", dstPath, filepath.Base(srcPath))
	cmd.Dir = filepath.Dir(srcPath)
	return CmdWrapper(cmd)
}

func Unzip(srcPath string, dstDir string) error {
	cmd := exec.Command("unzip", srcPath, "-d", dstDir)
	return CmdWrapper(cmd)
}

func Cnv(srcPath string, dstPath string) error {
	cmd := exec.Command("ffmpeg", "-i", srcPath, dstPath)
	err := CmdWrapper(cmd)
	if err != nil {
		return err
	}

	return nil
}

func IsThisExt(filename string, ext string) bool {
	result := strings.Split(filename, ".")
	if len(result) < 2 {
		return ext == ""
	}
	return result[len(result)-1] == ext
}

func GetBasenameAndExt(filename string) (baseName string, ext string) {
	result := strings.Split(filename, ".")
	length := len(result)
	if length < 2 {
		return filename, ""
	}
	return strings.Join(result[0:length-1], "."), result[length-1]
}
