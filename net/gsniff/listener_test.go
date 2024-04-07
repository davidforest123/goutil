package gsniff

import (
	"fmt"
	"google.golang.org/grpc"
	"goutil/basic/gerrors"
	"goutil/basic/gtest"
	"net"
	"testing"
)

func TestTrySniff(t *testing.T) {
	l, err := net.Listen("tcp", "0.0.0.0:34567")
	gtest.Assert(t, err)

	go func() {
		_, err := grpc.Dial("127.0.0.1:34567", grpc.WithInsecure())
		gtest.Assert(t, gerrors.Wrap(err, "rpc.Dial"))
		//rpcConn.Call("a", grpcs.NewRequest(), nil)
	}()

	for {
		rawConn, err := l.Accept()
		gtest.Assert(t, gerrors.Wrap(err, "Accept"))
		sniffConn, err := trySniff(rawConn, AllMatchers...)
		gtest.Assert(t, gerrors.Wrap(err, "TrySniff"))
		fmt.Println(sniffConn.SniffedInfo())
	}

	l.Close()
}

func TestWrapConn(t *testing.T) {
	network := "tcp"
	addr := "localhost:23456"
	usage1, err := MakeUsage("@rpcapi!")
	gtest.Assert(t, err)

	l, err := net.Listen(network, addr)
	gtest.Assert(t, err)

	go func() {
		c, err := net.Dial(network, addr)
		gtest.Assert(t, err)
		_, err = NewClient(c, usage1)
		gtest.Assert(t, err)
		fmt.Println("Dial usage: ", usage1)
	}()

	rawConn, err := l.Accept()
	gtest.Assert(t, err)
	conn, err := trySniff(rawConn, AllMatchers...)
	gtest.Assert(t, err)
	fmt.Println("Listen usage: ", conn.SniffedInfo().Usage)
}
