package waitgroup

import "sync"

type WgWrapper struct {
	sync.WaitGroup
}

func (w *WgWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
