package gsniff

import (
	"goutil/basic/gerrors"
	"goutil/container/gany"
	"goutil/net/gnet"
	"net"
	"strings"
)

// NewClient wraps net.Conn to UMUxConn before any Read/Write.
func NewClient(conn net.Conn, usage Usage) (net.Conn, error) {
	if strings.Contains(gany.Type(conn), "UMuxConn") {
		return nil, gerrors.Errorf("WrapConn: %s type not accepted, don't wraps UMuxConn multiple times", gany.Type(conn))
	}
	if conn == nil {
		return nil, gerrors.New("WrapConn: null connection")
	}
	written := 0
	for written < len(usage) {
		n, err := conn.Write(usage.Bytes()[written:])
		if err != nil {
			return nil, err
		}
		written += n
	}

	return conn, nil
}

func Dial(network, addr string, dialer gnet.Dialer, usage Usage) (net.Conn, error) {
	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	written := 0
	for written < len(usage) {
		n, err := conn.Write(usage.Bytes()[written:])
		if err != nil {
			return nil, err
		}
		written += n
	}
	return conn, nil
}
