package gcompress

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
)

type (
	BufferCompress struct {
		algo Comp
	}
)

func NewBufferCompress(inBuf []byte, outBuf bytes.Buffer, algo Comp) (*BufferCompress, error) {
	return nil, nil
}

// UnCompressGzip un-compress using Gzip
func UnCompressGzip(buf []byte) ([]byte, error) {
	if buf == nil || len(buf) == 0 {
		return nil, errors.New("invalid input buf!")
	}
	rbuf := bytes.NewReader(buf)
	gzipReader, err := gzip.NewReader(rbuf)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(gzipReader)
	return res, err
}

// CompressGzip compress using Gzip
func CompressGzip(buf []byte) []byte {
	var res bytes.Buffer
	gzipWriter := gzip.NewWriter(&res)
	gzipWriter.Write(buf)
	gzipWriter.Close()
	return res.Bytes()
}
