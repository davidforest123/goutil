package gwg

import (
	"sync"
	"time"
)

type (
	WaitSwitch struct {
		ch chan struct{}
		mu sync.Mutex
	}
)

func NewWaitSwitch() *WaitSwitch {
	return &WaitSwitch{ch: make(chan struct{}, 1)}
}

func (w *WaitSwitch) SetBlock() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for len(w.ch) > 0 {
		<-w.ch
	}
}

func (w *WaitSwitch) SetUnBlock() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.ch) == 0 {
		w.ch <- struct{}{}
	}
}

// Wait waits until condition satisfied.
// If waitTimeout == nil: waits until: 1, SetUnBlock() called.
// If waitTimeout != nil: waits until: 1, SetUnBlock() called, or 2, timeout reached.
func (w *WaitSwitch) Wait(waitTimeout *time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if waitTimeout != nil {
		ticker := time.NewTicker(*waitTimeout)
		select {
		case <-ticker.C:
		case <-w.ch:
		}
	} else {
		<-w.ch
	}
}
