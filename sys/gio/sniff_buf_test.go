package gio

import (
	"bytes"
	"github.com/davidforest123/goutil/basic/gtest"
	"testing"
)

func TestWrapConn(t *testing.T) {
	buf := bytes.NewReader([]byte("0123456789"))
	sbuf := NewSniffBuf(buf)
	nr := sbuf.NormalReader()
	rr := sbuf.RewindReader()

	buf1 := make([]byte, 1)

	n, err := rr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '0', "should be `0` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = rr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '1', "should be `1` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	rr.Rewind()
	n, err = rr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '0', "should be `0` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = nr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '0', "should be `0` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = nr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '1', "should be `1` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = nr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '2', "should be `2` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = nr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '3', "should be `3` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = nr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '4', "should be `4` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = rr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '5', "should be `5` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	n, err = rr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '6', "should be `6` but got `%d`", buf1[0])
	//fmt.Println(buf1)

	rr.Rewind()
	n, err = rr.Read(buf1)
	gtest.AssertTrue(t, n == 1, "should read 1 byte, but read %d", n)
	gtest.Assert(t, err)
	gtest.AssertTrue(t, buf1[0] == '5', "should be `5` but got `%d`", buf1[0])
	//fmt.Println(buf1)
}
