package gio

import (
	"fmt"
	"goutil/basic/gtest"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestReadFull(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:12345")
	gtest.Assert(t, err)
	go func() {
		for {
			lis.Accept()
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	gtest.Assert(t, err)

	buf := make([]byte, 10)
	timeout := 10 * time.Second
	n, err := ReadFull(conn, buf, &timeout)
	gtest.Assert(t, err)
	fmt.Println(n)
}

func TestRead(t *testing.T) {
	rd := strings.NewReader("hello")
	rdc := io.NopCloser(rd)
	readTimeout := time.Millisecond * 10
	tickerTimeout := time.NewTicker(readTimeout)
	chDone := make(chan struct{}, 1)

	p := make([]byte, 1024)
	n := 0
	err := error(nil)
	go func() {
		defer close(chDone)
		defer fmt.Println("exit")
		n, err = io.ReadFull(rdc, p)
		fmt.Println(n, err)
	}()

	select {
	case <-chDone:
		fmt.Println("read done")
	case <-tickerTimeout.C:
		rdc.Close()
		fmt.Println("timeout")
		return
	}
}

func TestReadLine(t *testing.T) {
	rd := strings.NewReader("good\nmorning\nhello\nworld")

	buf, err := ReadLine(rd, false, 1024)
	if err != nil {
		gtest.Assert(t, err)
	}
	if string(buf) != "good" {
		gtest.PrintlnExit(t, "ReadLine test error1")
	}

	buf, err = ReadLine(rd, false, 1024)
	if err != nil {
		gtest.Assert(t, err)
	}
	if string(buf) != "morning" {
		gtest.PrintlnExit(t, "ReadLine test error2")
	}

	buf, err = ReadLine(rd, false, 1024)
	if err != nil {
		gtest.Assert(t, err)
	}
	if string(buf) != "hello" {
		gtest.PrintlnExit(t, "ReadLine test error3")
	}

	buf, err = io.ReadAll(rd)
	if err != nil {
		gtest.Assert(t, err)
	}
	if string(buf) != "world" {
		gtest.PrintlnExit(t, "ReadLine test error4")
	}
}
