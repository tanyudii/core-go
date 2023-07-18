package errutil

import (
	"errors"
	"strings"
)

type Wrapper struct {
	name string
	fn   string
	err  error
}

func (w *Wrapper) Error() string {
	return w.name
}

func (w *Wrapper) SetFn(fn string) *Wrapper {
	w.fn = fn
	return w
}

func (w *Wrapper) Wrap(err error, msg ...string) error {
	messages := []string{w.name}
	if w.fn != "" {
		messages = append(messages, w.fn)
	}
	messages = append(messages, msg...)
	w.err = Wrap(err, errors.New(strings.Join(messages, ": ")))
	return w.err
}

func NewWrapper(name string) *Wrapper {
	return &Wrapper{
		name: name,
	}
}
