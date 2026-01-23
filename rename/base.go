package rename

import (
	"fmt"
	"strings"
)

type Namer interface {
	Next(oldName string) string
}

type NameGen struct {
	Basename string
	Ext      string
	count    int
}

func (ng *NameGen) Next(_ string) string {
	name := fmt.Sprintf("%v%v", ng.Basename, ng.count)
	if ng.Ext != "" {
		name = fmt.Sprintf("%v.%v", name, strings.Trim(ng.Ext, "."))
	}
	ng.count += 1
	return name
}
