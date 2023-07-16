package errutil

import "errors"

type Wrapper struct {
	name string
}

func (w *Wrapper) Wrap(err error) error {
	return Wrap(errors.New(w.name), err)
}

func NewWrapper(name string) *Wrapper {
	return &Wrapper{
		name: name,
	}
}
