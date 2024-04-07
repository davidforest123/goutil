package gio

import (
	"bytes"
	"io"
)

/**
What is SniffBuf?
Technically, io.Reader is treated like a stream, you cannot read it multiple times,
Read() will delete what already been read, then data will be lost in this stream.
That's why I create SniffBuf, it provides two different Readers, a NormalReader and a RewindReader.
NormalReader behaves like the standard io.Reader, it is non-rewindable.
But RewindReader is rewindable, and data read through its Read interface is cached in a temporary buffer,
and can be re-read through Read interface after Rewind.

How to use SniffBuf?
Note that data that has already been read through the NormalReader is not rewindable and cannot be read
through the RewindReader, so you should first read the data through the RewindReader to do the job you want,
and then read the data through the NormalReader as you would normally do with standard io.Reader.
*/

type (
	SniffBuf struct {
		src    io.Reader     // data source
		buffer *bytes.Buffer // buffered bytes when sniffing
		off    int           // offset
	}
	NormalReader SniffBuf // normal reader which isn't rewindable
	RewindReader SniffBuf // special reader which is rewindable
)

// NewSniffBuf creates a new SniffBuff
func NewSniffBuf(src io.Reader) *SniffBuf {
	if src == nil {
		panic("src is nil!")
	}
	return &SniffBuf{
		src:    src,
		off:    0,
		buffer: bytes.NewBuffer(nil),
	}
}

// NormalReader exports a normal io.Reader
func (sc *SniffBuf) NormalReader() *NormalReader {
	return (*NormalReader)(sc)
}

// Normal(not rewindable) io.Read()
func (sc *NormalReader) Read(p []byte) (n int, err error) {
	if sc.buffer.Len() > 0 {
		n, err = sc.buffer.Read(p)
		sc.off -= n
		if sc.off < 0 {
			sc.off = 0
		}
		return n, err
	} else {
		return sc.src.Read(p)
	}
}

// RewindReader exports a RewindReader
func (sc *SniffBuf) RewindReader() *RewindReader {
	return (*RewindReader)(sc)
}

// Rewindable io.Read()
func (sc *RewindReader) Read(p []byte) (n int, err error) {
	if sc.off < sc.buffer.Len() {
		n = copy(p, sc.buffer.Bytes()[sc.off:])
		sc.off += n
		return n, err
	} else {
		n, err = io.TeeReader(sc.src, sc.buffer).Read(p)
		sc.off += n
		return n, err
	}
}

// Rewind resets offset to 0.
func (sc *RewindReader) Rewind() {
	sc.off = 0
}
