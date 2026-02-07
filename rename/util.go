package rename

import (
	"fmt"
	"github.com/XV-521/fileops/core"
	"github.com/XV-521/fileops/internal/util"
)

type nameGen struct {
	basename string
	ext      string
	count    int
}

func (ng *nameGen) Next(_ core.EntryInfo) string {
	ng.count += 1
	name := fmt.Sprintf("%v%v", ng.basename, ng.count)
	if ng.ext != "" {
		name = fmt.Sprintf("%v%v", name, util.GetClearExt(ng.ext))
	}
	return name
}
