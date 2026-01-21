package internal

type Preparer[T any] interface {
	Check() error
	Normalize() (T, error)
}

func Prepare[T Preparer[T]](p T) (T, error) {
	err := p.Check()
	if err != nil {
		return p, err
	}
	return p.Normalize()
}
