//go:build !windows && !linux && !darwin && !dragonfly && !freebsd && !netbsd && !openbsd

package gnet

func Control(network, address string, c syscall.RawConn) error {
	return nil
}
