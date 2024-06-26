package socks5internal

import (
	"bufio"
	"fmt"
	"net"

	"github.com/davidforest123/goutil/basic/glog"
	"github.com/davidforest123/goutil/net/gnet"
	"golang.org/x/net/context"
)

// Config is used to setup and configure a Server
type Config struct {
	// AuthMethods can be provided to implement custom authentication
	// By default, "auth-less" mode is enabled.
	// For password-based auth use UserPassAuthenticator.
	AuthMethods []Authenticator

	// If provided, username/password authentication is enabled,
	// by appending a UserPassAuthenticator to AuthMethods. If not provided,
	// and AUthMethods is nil, then "auth-less" mode is enabled.
	Credentials CredentialStore

	// Rules is provided to enable custom logic around permitting
	// various commands. If not provided, PermitAll is used.
	Rules RuleSet

	// Rewriter can be used to transparently rewrite addresses.
	// This is invoked before the RuleSet is invoked.
	// Defaults to NoRewrite.
	Rewriter AddressRewriter

	// BindIP is used for bind or udp associate
	BindIP net.IP

	// Log can be used to provide a custom log target.
	Log glog.Interface

	// active proxy go routines count
	aliveProxy int64

	// Optional function for dialing out
	// Users can customize DNS lookup behavior in 'CustomDial'.
	CustomDial func(ctx context.Context, network, addr string) (net.Conn, error)
}

// Server is reponsible for accepting connections and handling
// the details of the SOCKS5 protocol
type Server struct {
	config      *Config
	authMethods map[uint8]Authenticator
}

// New creates a new Server and potentially returns an error
func New(conf *Config) (*Server, error) {
	// Ensure we have at least one authentication method enabled
	if len(conf.AuthMethods) == 0 {
		if conf.Credentials != nil {
			conf.AuthMethods = []Authenticator{&UserPassAuthenticator{conf.Credentials}}
		} else {
			conf.AuthMethods = []Authenticator{&NoAuthAuthenticator{}}
		}
	}

	// Ensure we have a rule set
	if conf.Rules == nil {
		conf.Rules = PermitAll()
	}

	// Ensure we have a log target
	if conf.Log == nil {
		conf.Log = glog.DefaultLogger
	}

	server := &Server{
		config: conf,
	}

	server.authMethods = make(map[uint8]Authenticator)

	for _, a := range conf.AuthMethods {
		server.authMethods[a.GetCode()] = a
	}

	return server, nil
}

// SetCustomDialer sets custom dialer for requests,
// users can customize DNS lookup behavior in 'dialer'.
// This operation is optional.
func (s *Server) SetCustomDialer(dialer gnet.DialWithCtxFunc) {
	s.config.CustomDial = dialer
}

// SetCustomLogger sets custom logger
func (s *Server) SetCustomLogger(log glog.Interface) {
	s.config.Log = log
}

// ListenAndServe is used to create a listener and serve on it
func (s *Server) ListenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

// Serve is used to serve connections from a listener
func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err := s.ServeConn(conn); err != nil {
				s.config.Log.Erro(err)
			}
		}()
	}
	return nil
}

// ServeConn is used to serve a single connection.
func (s *Server) ServeConn(conn net.Conn) error {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	// Read the version byte
	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		s.config.Log.Errof("[ERR] socks: Failed to get version byte: %v", err)
		return err
	}

	// Ensure we are compatible
	if version[0] != Version {
		err := fmt.Errorf("Unsupported SOCKS version: %v", version)
		s.config.Log.Errof("[ERR] socks: %v", err)
		return err
	}

	// Authenticate the connection
	authContext, err := s.authenticate(conn, bufConn)
	if err != nil {
		err = fmt.Errorf("Failed to authenticate: %v", err)
		s.config.Log.Errof("[ERR] socks: %v", err)
		return err
	}

	request, err := NewRequest(bufConn)
	if err != nil {
		if err == unrecognizedAddrType {
			if err := sendReply(conn, addrTypeNotSupported, nil); err != nil {
				err = fmt.Errorf("Failed to send reply: %v", err)
				s.config.Log.Errof("[ERR] socks: %v", err)
				return err
			}
		}
		err = fmt.Errorf("Failed to read destination address: %v", err)
		s.config.Log.Errof("[ERR] socks: %v", err)
		return err
	}
	request.AuthContext = authContext
	if client, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: client.Port}
	}

	// Process the client request
	if err := s.handleRequest(request, conn); err != nil {
		s.config.Log.Infof("waiting for jumpbox(from %s, dest %s) to be available...close it", request.RemoteAddr.String(), request.DestAddr.FQDN)
		conn.Close()
		return err
	}

	return nil
}
