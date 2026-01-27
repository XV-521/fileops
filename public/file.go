package public

import (
	"fmt"
	"github.com/XV-521/fileops/internal"
	"os"
	"syscall"
	"time"
)

type FileData struct {
	Name    string
	Ext     string
	Size    int64
	CrtTime time.Time
	ModTime time.Time
}

func GetFileData(info os.FileInfo) FileData {
	name, ext := internal.GetBasenameAndExt(info.Name())

	var crtTime time.Time

	stat, ok := info.Sys().(*syscall.Stat_t)
	if ok {
		crtTime = time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	} else {
		fmt.Println("can not get crtTime")
	}

	return FileData{
		Name:    name,
		Ext:     ext,
		Size:    info.Size(),
		CrtTime: crtTime,
		ModTime: info.ModTime(),
	}
}
