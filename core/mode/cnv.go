package mode

import (
	"errors"
	"github.com/XV-521/fileops/internal/util"
	"strings"
)

type CnvType int

const (
	CnvUn CnvType = iota
	CnvV
	CnvA
	CnvI
)

var UnsupportedCnvTypeErr = errors.New("unsupported cnv type")

func GetCnvType(name string) CnvType {
	name = strings.ToLower(name)

	switch {
	case strings.HasSuffix(name, ".mp4"),
		strings.HasSuffix(name, ".mov"),
		strings.HasSuffix(name, ".avi"),
		strings.HasSuffix(name, ".mkv"),
		strings.HasSuffix(name, ".flv"):
		return CnvV

	case strings.HasSuffix(name, ".mp3"),
		strings.HasSuffix(name, ".wav"),
		strings.HasSuffix(name, ".flac"),
		strings.HasSuffix(name, ".aac"),
		strings.HasSuffix(name, ".m4a"):
		return CnvA

	case strings.HasSuffix(name, ".jpg"),
		strings.HasSuffix(name, ".jpeg"),
		strings.HasSuffix(name, ".png"),
		strings.HasSuffix(name, ".webp"),
		strings.HasSuffix(name, ".avif"),
		strings.HasSuffix(name, ".gif"):
		return CnvI

	default:
		return CnvUn
	}
}

type CnvFn func(srcPath string, dstPath string) error

func GetCnvFn(ct CnvType) (CnvFn, error) {

	switch ct {
	case CnvV:
		return util.CnvForVideo, nil
	case CnvA:
		return util.CnvForAudio, nil
	case CnvI:
		return util.CnvForImage, nil

	default:
		return nil, UnsupportedCnvTypeErr
	}
}
