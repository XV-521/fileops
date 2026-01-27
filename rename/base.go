package rename

import (
	"fmt"
	"os"
	"strings"
)

type Namer interface {
	Next(info os.FileInfo) string
}

type nameGen struct {
	basename string
	ext      string
	count    int
}

func (ng *nameGen) Next(_ os.FileInfo) string {
	ng.count += 1
	name := fmt.Sprintf("%v%v", ng.basename, ng.count)
	if ng.ext != "" {
		name = fmt.Sprintf("%v.%v", name, strings.Trim(ng.ext, "."))
	}
	return name
}
