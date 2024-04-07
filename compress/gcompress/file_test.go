package gcompress

import (
	"bytes"
	"github.com/yeka/zip"
	"goutil/basic/gerrors"
	"goutil/basic/gtest"
	"io"
	"log"
	"os"
	"testing"
)

func TestNewCompFile(t *testing.T) {
	cf, err := NewCompFile("/Users/tonystark/Downloads/test.zip", CompZip, "123456789", nil)
	gtest.Assert(t, err)

	w1, err := cf.AddFile("1.txt", EncryptAES256, LevelDeflate)
	gtest.Assert(t, gerrors.Wrap(err, "AddFile 1"))
	_, err = w1.Write([]byte("hello 1"))
	gtest.Assert(t, gerrors.Wrap(err, "Write 1"))
	err = w1.Flush()
	gtest.Assert(t, gerrors.Wrap(err, "Flush 1"))
	err = w1.Close()
	gtest.Assert(t, gerrors.Wrap(err, "Close 1"))

	/*w2, err := cf.AddFile("2.txt", EncryptAES128, LevelDeflate)
	gtest.Assert(t, gerrors.Wrap(err, "AddFile 2"))
	_, err = w2.Write([]byte("hello 2"))
	gtest.Assert(t, gerrors.Wrap(err, "Write 2"))
	err = w2.Flush()
	gtest.Assert(t, gerrors.Wrap(err, "Flush 2"))
	err = w2.Close()
	gtest.Assert(t, gerrors.Wrap(err, "Close 2"))*/

	err = cf.Close()
	gtest.Assert(t, gerrors.Wrap(err, "Close"))
}

func TestCompressGzip(t *testing.T) {
	contents := []byte("Hello World")
	fzip, err := os.Create(`/Users/tonystark/Downloads/alex-test.zip`)
	if err != nil {
		log.Fatalln(err)
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	w, err := zipw.Encrypt(`test.txt`, `golang`, zip.AES256Encryption)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, bytes.NewReader(contents))
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
}
