package out

import "fmt"

type Render string

func (r Render) Sprint(format string, a ...any) string {
	content := fmt.Sprintf(format, a...)
	return fmt.Sprintf("%v%v\033[0m", r, content)
}

func (r Render) Printf(format string, a ...any) {
	content := fmt.Sprintf(format, a...)
	content = fmt.Sprintf("%v%v\033[0m", r, content)
	fmt.Print(content)
}

var (
	Skip Render = "\033[37m"
	Warn Render = "\033[33m"
	Fail Render = "\033[31m"
)
