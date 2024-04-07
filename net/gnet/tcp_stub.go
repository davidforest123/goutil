//go:build !linux
// +build !linux

package gnet

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	"net"
)

func EnableBBR(conn *net.TCPConn) error {
	return gerrors.ErrNotSupport
}

func GetBBR(conn *net.TCPConn) (*BBRInfo, error) {
	return nil, gerrors.ErrNotSupport
}
