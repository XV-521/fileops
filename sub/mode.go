package sub

import (
	"github.com/XV-521/fileops/v2/sub/addsub"
)

type Mode struct {
	SrcDir    string
	DstDir    string
	Model     string
	Lang      string
	SubExt    string
	ST        addsub.SubType
	Style     addsub.SimpleStyle
	RetainSub bool
	Rec       bool
	Strict    bool
}
