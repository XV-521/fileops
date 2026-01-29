package public

import (
	"errors"
	"github.com/XV-521/fileops/internal"
	"strings"
)

type ZipFn func(srcPath string, dstPath string, pwd string) error

type UnzipFn func(srcPath string, dstDir string, pwd string) error

type ZipType int

const (
	ZipUn ZipType = iota
	ZipB
	ZipS
	ZipT
	ZipR
)

var UnsupportedZipTypeErr = errors.New("unsupported zip type")
var UnsupportedUnzipTypeErr = errors.New("unsupported unzip type")

func GetZipType(name string) ZipType {

	name = strings.ToLower(name)
	switch {
	case strings.HasSuffix(name, ".zip"),
		strings.HasSuffix(name, ".zipx"):
		return ZipB
	case strings.HasSuffix(name, ".7z"):
		return ZipS
	case strings.HasSuffix(name, ".tar"),
		strings.HasSuffix(name, ".tar.gz"),
		strings.HasSuffix(name, ".tgz"),
		strings.HasSuffix(name, ".tar.xz"),
		strings.HasSuffix(name, ".tar.bz2"):
		return ZipT
	case strings.HasSuffix(name, ".rar"):
		return ZipR
	default:
		return ZipUn
	}
}

func GetZipFn(zt ZipType) (ZipFn, error) {

	switch zt {
	case ZipB:
		return internal.Zip, nil
	case ZipS:
		return internal.SevenZip, nil
	case ZipT:
		return internal.TarZip, nil
	default:
		return nil, UnsupportedZipTypeErr
	}
}

func GetUnzipFn(zt ZipType) (UnzipFn, error) {

	switch zt {
	case ZipB:
		return internal.Unzip, nil
	case ZipS:
		return internal.SevenUnzip, nil
	case ZipT:
		return internal.TarUnzip, nil
	case ZipR:
		return internal.RarUnzip, nil
	default:
		return nil, UnsupportedUnzipTypeErr
	}
}
