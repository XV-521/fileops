package mode

import (
	"errors"
	"github.com/XV-521/fileops/internal/util"
	"strings"
)

type PackType int

const (
	PackUn PackType = iota
	PackZ
	PackS
	PackT
	PackR
)

var UnsupportedPackTypeErr = errors.New("unsupported pack type")
var UnsupportedUnpackTypeErr = errors.New("unsupported unpack type")

func GetPackType(name string) PackType {

	name = strings.ToLower(name)
	switch {
	case strings.HasSuffix(name, ".zip"),
		strings.HasSuffix(name, ".zipx"):
		return PackZ
	case strings.HasSuffix(name, ".7z"):
		return PackS
	case strings.HasSuffix(name, ".tar"),
		strings.HasSuffix(name, ".tar.gz"),
		strings.HasSuffix(name, ".tgz"),
		strings.HasSuffix(name, ".tar.xz"),
		strings.HasSuffix(name, ".tar.bz2"):
		return PackT
	case strings.HasSuffix(name, ".rar"):
		return PackR
	default:
		return PackUn
	}
}

type PackFn func(srcPath string, dstPath string, pwd string) error

func GetPackFn(pt PackType) (PackFn, error) {

	switch pt {
	case PackZ:
		return util.Zip, nil
	case PackS:
		return util.Seven, nil
	case PackT:
		return util.Tar, nil
	default:
		return nil, UnsupportedPackTypeErr
	}
}

type UnpackFn func(srcPath string, dstDir string, pwd string) error

func GetUnpackFn(pt PackType) (UnpackFn, error) {

	switch pt {
	case PackZ:
		return util.UnSeven, nil
	case PackS:
		return util.UnSeven, nil
	case PackT:
		return util.UnTar, nil
	case PackR:
		return util.UnRar, nil
	default:
		return nil, UnsupportedUnpackTypeErr
	}
}

func CreatePackExt(pt PackType) (string, error) {
	switch pt {
	case PackZ:
		return ".zip", nil
	case PackS:
		return ".7z", nil
	case PackT:

		return ".tar.gz", nil
	default:
		return "", UnsupportedPackTypeErr
	}
}
