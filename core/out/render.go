package out

import "fmt"

type Render string

func (r Render) wrap(s string) string {
	return fmt.Sprintf("%v%v\033[0m", r, s)
}

func (r Render) Sprint(a ...any) string {
	s := fmt.Sprint(a...)
	return r.wrap(s)
}

func (r Render) Print(a ...any) {
	s := fmt.Sprint(a...)
	fmt.Print(r.wrap(s))
}

func (r Render) Sprintln(a ...any) string {
	s := fmt.Sprintln(a...)
	return r.wrap(s)
}

func (r Render) Println(a ...any) {
	s := fmt.Sprintln(a...)
	fmt.Print(r.wrap(s))
}

func (r Render) Sprintf(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return r.wrap(s)
}

func (r Render) Printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Print(r.wrap(s))
}

var (
	Skip Render = "\033[37m"
	Warn Render = "\033[33m"
	Fail Render = "\033[31m"
)
