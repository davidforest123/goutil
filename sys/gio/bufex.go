package gio

// from https://go.dev/play/p/CX5ojtTOm4
// FIXME: it seems not correct, bufio.Peek will block too, please test it.

import (
	"bufio"
	"github.com/davidforest123/goutil/basic/gerrors"
	"io"
	"runtime"
	"time"
)

const (
	bufferSize = 4096
)

type (
	BufEx struct {
		b         *bufio.Reader
		rdTimeout time.Duration
		ch        <-chan error
	}
)

func NewBufEx(r io.Reader) *BufEx {
	return &BufEx{b: bufio.NewReaderSize(r, bufferSize), rdTimeout: -1}
}

// SetReadTimeout sets the timeout for all future Read calls as follows:
//
//	 t:
//		t<0:  block
//		t==0: poll
//		t>0:  timeout after t
func (r *BufEx) SetReadTimeout(t time.Duration) time.Duration {
	prev := r.rdTimeout
	r.rdTimeout = t
	return prev
}

func (r *BufEx) Read(b []byte) (n int, err error) {
	if r.ch == nil {
		if r.rdTimeout < 0 || r.b.Buffered() > 0 {
			return r.b.Read(b)
		}
		ch := make(chan error, 1)
		r.ch = ch
		go func() {
			_, err := r.b.Peek(1)
			ch <- err
		}()
		runtime.Gosched()
	}

	if r.rdTimeout < 0 {
		err = <-r.ch // Block
	} else {
		select {
		case err = <-r.ch: // Poll
		default:
			if r.rdTimeout == 0 {
				return 0, gerrors.ErrTimeout
			}
			select {
			case err = <-r.ch: // Timeout
			case <-time.After(r.rdTimeout):
				return 0, gerrors.ErrTimeout
			}
		}
	}
	r.ch = nil
	if r.b.Buffered() > 0 {
		n, _ = r.b.Read(b)
	}
	return
}
