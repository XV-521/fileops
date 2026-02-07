package util

import (
	"fmt"
	"path/filepath"
	"strings"
)

func IsThisExt(filename string, ext string) bool {
	if ext == "" {
		return filepath.Ext(filename) == ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return strings.EqualFold(filepath.Ext(filename), ext)
}

func GetClearExt(ext string) string {
	ext = strings.TrimSpace(ext)
	if ext == "" || ext == "." {
		return ""
	}
	ext = strings.TrimPrefix(ext, ".")
	return fmt.Sprintf(".%v", ext)
}

func GetBasenameAndExt(filename string) (basename string, ext string) {
	ext = filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext), ext
}
