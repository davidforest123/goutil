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
		/*log         glog.Interface
		dialer      gnet.DialWithCtxFunc
		dnsRequire  gnet.LookupIPRequiredFunc
		dnsResolver gnet.LookupIPWithCtxFunc*/
	}
)

// NewServer create new socks5 server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
func NewServer() (*Server, error) {
	s := &Server{}

	conf := &socks5internal.Config{}
	/*if s.dialer != nil {
		conf.Dial = s.dialer
	}
	if s.dnsRequire != nil {
		conf.ResolveRequire = s.dnsRequire
	}
	if s.dnsResolver != nil {
		conf.Resolver = s.dnsResolver
	}
	if s.log != nil {
		conf.Log = s.log
	} else {
		conf.Log = glog.DefaultLogger
	}*/
	err := error(nil)
	s.srv, err = socks5internal.New(conf)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// SetCustomDialer sets custom dialer for requests.
// This operation is optional.
func (s *Server) SetCustomDialer(dialer gnet.DialWithCtxFunc) {
	//s.dialer = dialer
	s.srv.SetCustomDialer(dialer)
}

// SetCustomDNSResolver sets custom DNS resolver for requests.
// This operation is optional.
// dnsRequire: check whether dns-lookup required, if false, dnsResolver will not be called,
// and domain (but not net.IP) will be sent to custom dialer, and user want to do dns-lookup in custom dialer callback.
func (s *Server) SetCustomDNSResolver(dnsRequire gnet.LookupIPRequiredFunc, dnsResolver gnet.LookupIPWithCtxFunc) {
	//s.dnsRequire = dnsRequire
	//s.dnsResolver =
	s.srv.SetCustomDNSResolver(dnsRequire, dnsResolver)
}

func (s *Server) SetCustomLogger(log glog.Interface) {
	//s.log = log
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
