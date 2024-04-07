package gio

import (
	"bytes"
	"errors"
	"github.com/davidforest123/goutil/basic/gerrors"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type (
	SetDeadlineCallback func(t time.Time) error
	CopiedSizeCallback  func(size int64)
	ErrNotify           func(err error)

	TimeoutR struct {
		rd       io.ReadCloser
		deadline time.Time
	}

	TimeoutW struct {
		wd       io.WriteCloser
		deadline time.Time
	}

	ReaderTimeout interface {
		io.ReadCloser
		SetReadDeadline(t time.Time) error
	}

	TimeoutReadCloser interface {
		io.ReadCloser
		SetReadDeadline(t time.Time) error
	}

	TimeoutWriteCloser interface {
		io.WriteCloser
		SetWriteDeadline(t time.Time) error
	}

	TimeoutReadWriteCloser interface {
		io.ReadWriteCloser
		SetReadDeadline(t time.Time) error
		SetWriteDeadline(t time.Time) error
		SetDeadline(t time.Time) error
	}
)

func TwoWaysCopy(s1, s2 io.ReadWriteCloser, errNotify ErrNotify) {
	// Memory optimized io.Copy function specified for this library
	const bufSize = 4096
	genericCopy := func(dst io.Writer, src io.Reader) (written int64, err error) {
		// If the reader has a WriteTo method, use it to do the copy.
		// Avoids an allocation and a copy.
		if wt, ok := src.(io.WriterTo); ok {
			return wt.WriteTo(dst)
		}
		// Similarly, if the writer has a ReadFrom method, use it to do the copy.
		if rt, ok := dst.(io.ReaderFrom); ok {
			return rt.ReadFrom(src)
		}

		// fallback to standard io.CopyBuffer
		buf := make([]byte, bufSize)
		return io.CopyBuffer(dst, src, buf)
	}

	// start tunnel & wait for tunnel termination
	streamCopy := func(dst io.Writer, src io.ReadCloser, chClose chan struct{}) {
		if _, err := genericCopy(dst, src); err != nil {
			if err != nil {
				errNotify(err)
			}
		}

		// Notice: you'd better don't close them, sometimes this could make something wrong,
		// for example closed a valid TUN/TAP device.
		// So let the user judge and decide whether to close them.
		// s1.Close()
		// s2.Close()
		close(chClose)
	}

	chClose21 := make(chan struct{}, 1)
	chClose12 := make(chan struct{}, 1)
	go streamCopy(s2, s1, chClose21)
	go streamCopy(s1, s2, chClose12)

	// continue if any copy routine exits
	select {
	case <-chClose21:
	case <-chClose12:
	}
}

// Note: among the two, 'dst' and 'src', there must be at least one net.Conn
// based on standard library io.copyBuffer
func copyBufferTimeout(dst io.Writer, src io.Reader, buf []byte, noDataTimeout time.Duration) (written int64, err error) {
	var errInvalidWrite = errors.New("invalid write result")
	dstConn, dstOk := dst.(net.Conn)
	srcConn, srcOk := src.(net.Conn)
	if !dstOk && !srcOk {
		return 0, gerrors.New("copyBufferTimeout only support at least one net.Conn")
	}
	readDeadline := time.Now().Add(noDataTimeout)
	writeDeadline := time.Now().Add(noDataTimeout)

	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {

		// set deadline for src
		if srcOk {
			if srdErr := srcConn.SetReadDeadline(readDeadline); srdErr != nil {
				return written, srdErr
			}
		}

		nr, er := src.Read(buf)

		// reset readDeadline
		if nr > 0 {
			readDeadline = time.Now().Add(noDataTimeout)
		}

		if nr > 0 {

			// set deadline for dst
			if dstOk {
				if swdErr := dstConn.SetWriteDeadline(writeDeadline); swdErr != nil {
					return written, swdErr
				}
			}

			nw, ew := dst.Write(buf[0:nr])

			// reset writeDeadline
			if nw > 0 {
				writeDeadline = time.Now().Add(noDataTimeout)
			}

			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// CopyTimeout is a Copy function with no data timeout parameter.
// Note: among the two, 'dst' and 'src', there must be at least one net.Conn
// More: https://groups.google.com/g/golang-nuts/c/byfD0YtIl0I
func CopyTimeout(dst io.Writer, src io.Reader, noDataTimeout time.Duration) (written int64, err error) {
	return copyBufferTimeout(dst, src, nil, noDataTimeout)
}

// A pool for stream copying
var xmitBuf sync.Pool

func init() {
	xmitBuf.New = func() interface{} {
		return make([]byte, 32768)
	}
}

// https://github.com/efarrer/iothrottler/
// https://github.com/jwkohnen/bwio/
// limit copy speed / duration, read copied size
func CopyEx(dst io.Writer, dstWriteCb SetDeadlineCallback, src io.Reader, srcReadCb SetDeadlineCallback, timeout time.Duration, sizeCallback CopiedSizeCallback) (written int64, err error) {
	buf := make([]byte, 32*1024)
	var nr, nw int
	var er, ew error
	lastNotify := time.Time{}

	/*
		// If the reader has a WriteTo method, use it to do the copy.
		// Avoids an allocation and a copy.
		if wt, ok := src.(WriterTo); ok {
			return wt.WriteTo(dst)
		}
		// Similarly, if the writer has a ReadFrom method, use it to do the copy.
		if rt, ok := dst.(ReaderFrom); ok {
			return rt.ReadFrom(src)
		}
		if buf == nil {
			buf = make([]byte, 32*1024)
		}
	*/

	for {
		if timeout > 0 {
			srcReadCb(time.Now().Add(timeout))
		}
		nr, er = src.Read(buf)
		if nr > 0 {
			if timeout > 0 {
				dstWriteCb(time.Now().Add(timeout))
			}
			nw, ew = dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)

				// notify callback
				if sizeCallback != nil {
					if time.Now().Sub(lastNotify) > time.Second {
						sizeCallback(written)
					}
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func ReaderToBytes(rd io.Reader) ([]byte, error) {
	b, err := io.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReaderToString(rd io.Reader) (string, error) {
	b, err := ReaderToBytes(rd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ReadCloserToBytes(rd io.ReadCloser) ([]byte, error) {
	b, err := io.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReadCloserToInterface(rd io.ReadCloser) (interface{}, error) {
	b, err := ReadCloserToBytes(rd)
	if err != nil {
		return nil, err
	}
	return interface{}(b), nil
}

func ReadCloserToString(rd io.ReadCloser) (string, error) {
	b, err := ReadCloserToBytes(rd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func BytesToReadCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}

func StringToReadCloser(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

// Read is based on standard io library.
func Read(r io.Reader, buf []byte, timeout *time.Duration) (n int, err error) {
	if timeout == nil {
		return r.Read(buf)
	}

	chDie := make(chan struct{}, 1)
	go func() {
		defer close(chDie)

		// FIXME: this will continue after timeout!
		n, err = r.Read(buf)
	}()

	ticker := time.NewTicker(*timeout)
	select {
	case <-ticker.C:
		return n, gerrors.ErrTimeout
	case <-chDie:
	}
	return n, err
}

// ReadFrom is based on standard io library.
func ReadFrom(dst *bytes.Buffer, src io.Reader, timeout *time.Duration) (n int64, err error) {
	if timeout == nil {
		return dst.ReadFrom(src)
	}

	chDie := make(chan struct{}, 1)
	go func() {
		defer close(chDie)

		// FIXME: this will continue after timeout!
		n, err = dst.ReadFrom(src)
	}()

	ticker := time.NewTicker(*timeout)
	select {
	case <-ticker.C:
		return n, gerrors.ErrTimeout
	case <-chDie:
	}
	return n, err
}

// ReadHead try best to Read until n==0 or error is returned.
func ReadHead(r io.Reader, timeout *time.Duration) (*bytes.Buffer, error) {
	if timeout == nil {
		b, err := io.ReadAll(r)
		return bytes.NewBuffer(b), err
	}

	buf := bytes.NewBuffer(nil)
	n := int64(0)
	readErr := error(nil)
	for {
		n, readErr = ReadFrom(buf, r, timeout)
		if n == 0 || readErr != nil {
			break
		}
	}
	return buf, readErr
}

// ReadFull reads from r into buf, until exactly len(buf) bytes read, or error returned,
// including gerrors.ErrTimeout, io.EOF, io.ErrUnexpectedEOF, and any other type of error.
func ReadFull(r io.Reader, buf []byte, timeout *time.Duration) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf), timeout)
}

// ReadAtLeast is based on standard io library.
func ReadAtLeast(r io.Reader, buf []byte, min int, timeout *time.Duration) (n int, err error) {
	if timeout == nil {
		return io.ReadAtLeast(r, buf, min)
	}

	chDie := make(chan struct{}, 1)
	go func() {
		defer close(chDie)

		// FIXME: this will continue after timeout!
		n, err = io.ReadAtLeast(r, buf, min)
	}()

	ticker := time.NewTicker(*timeout)
	select {
	case <-ticker.C:
		return n, gerrors.ErrTimeout
	case <-chDie:
	}

	return n, err
}

// ReadLine reads from r into buf, until '\n' read, or maxSize read, or error returned,
// including io.EOF, io.ErrUnexpectedEOF, and any other type of error.
// Note: this function is slow, DO NOT call this function too often!
func ReadLine(r io.Reader, keepN bool, maxSize int) ([]byte, error) {
	return ReadUntil(r, '\n', keepN, maxSize)
}

// ReadUntil reads until the first occurrence of delim in the input, or maxSize read, or error returned,
// including io.EOF, io.ErrUnexpectedEOF, and any other type of error.
// Note: this function is slow, DO NOT call this function too often!
func ReadUntil(r io.Reader, delim byte, keepDelim bool, maxSize int) ([]byte, error) {
	var result []byte
	b1 := make([]byte, 1)
	for {
		_, err := r.Read(b1)
		if err != nil {
			return nil, err
		}
		if b1[0] != delim || (b1[0] == delim && keepDelim) {
			result = append(result, b1...)
		}

		if b1[0] == delim {
			break
		}
		if maxSize > 0 && len(result) >= maxSize {
			break
		}
	}

	return result, nil
}
