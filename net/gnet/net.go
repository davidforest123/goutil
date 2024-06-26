package gnet

import (
	"net"
	"strings"

	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/net/gkcp"
)

// Dial multiple protocols.
func Dial(network, address string) (net.Conn, error) {
	switch strings.ToLower(network) {
	case "tcp":
		return net.Dial("tcp", address)
	case "udp":
		return net.Dial("udp", address)
	case "kcp":
		return gkcp.Dial(address)
	//case "quic":
	//	return gquic.Dial(address)
	default:
		return nil, gerrors.New("unsupported network %s", network)
	}
}

// ListenCop listens multiple connection-oriented protocols.
func ListenCop(network, addr string) (net.Listener, error) {
	switch strings.ToLower(network) {
	case "tcp", "tcp4", "tcp6":
		return net.Listen(network, addr)
	case "kcp":
		return gkcp.Listen(addr)
	//case "quic":
	//	return gquic.Listen(addr)
	default:
		return nil, gerrors.New("unsupported network %s", network)
	}
}

// ListenAny listens any supported protocols.
func ListenAny(network, addr string) (net.Listener, error) {
	switch strings.ToLower(network) {
	case "udp", "udp4", "udp6", "unixgram", "ip:1", "ip:icmp":
		return ListenPop(network, addr)
	default:
		return ListenCop(network, addr)
	}
}
