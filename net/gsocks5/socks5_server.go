package gsocks5

import (
	"net"

	"github.com/davidforest123/goutil/basic/glog"
	"github.com/davidforest123/goutil/net/gnet"
	"github.com/davidforest123/goutil/net/gsocks5/socks5internal"
)

type (
	Server struct {
		srv *socks5internal.Server
	}
)

// NewServer create new socks5 server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
func NewServer() (*Server, error) {
	s := &Server{}

	conf := &socks5internal.Config{}
	err := error(nil)
	s.srv, err = socks5internal.New(conf)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// SetCustomDialer sets custom dialer for requests,
// users can customize DNS lookup behavior in 'dialer'.
// This operation is optional.
func (s *Server) SetCustomDialer(dialer gnet.DialWithCtxFunc) {
	s.srv.SetCustomDialer(dialer)
}

func (s *Server) SetCustomLogger(log glog.Interface) {
	s.srv.SetCustomLogger(log)
}

// ServeAddr listen TCP at `listenAddr` and serves the net.Listener as Socks5 proxy server.
// This function waits until an error is returned.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
// listenAddr example: "127.0.0.1:8000"
func (s *Server) ServeAddr(listenAddr string) error {
	// Listen a SOCKS5 server
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	// Serve it
	return s.srv.Serve(lis)
}

// ServeListener takes an external incoming net.Listener and serves it as Socks5 proxy server.
func (s *Server) ServeListener(lis net.Listener) error {
	return s.srv.Serve(lis)
}

// ServeConn takes an external incoming net.Conn and serves it as Socks5 proxy request connection.
func (s *Server) ServeConn(conn net.Conn) error {
	return s.srv.ServeConn(conn)
}
