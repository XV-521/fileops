package rename

import (
	"fmt"
	"strings"
)

type EntryNameGen struct {
	Basename string
	Ext      string
	count    int
}

func (ng *EntryNameGen) Next() string {
	name := fmt.Sprintf("%v%v", ng.Basename, ng.count)
	if ng.Ext != "" {
		name = fmt.Sprintf("%v.%v", name, strings.Trim(ng.Ext, "."))
	}
	ng.count += 1
	return name
}
