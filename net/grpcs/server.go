package grpcs

// Notice:
// If a function of the server-side registered Receiver returns an error,
// the result of the output parameter will not be transmitted to client when the function is called by RPC.

// Registered Receiver member function sample:
// (r *Recv) Method(in InputParam, out *OutputParam) error

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync/atomic"
)

type (
	Server struct {
		rpcSvr *rpc.Server
		//netSvr        net.Listener
		paramChecker  ParamChecker
		onRequestUser OnRequest
		rpcType       RpcType
	}

	Svr Server

	InitChecker func() ParamChecker
	OnRequest   func(in Request, out *Reply) error
)

var (
	jsonRpcRegisterDone = int32(0)
)

// NewServer creates new rpc server.
// before every 'onReq' call, rpc server will check input and out put param with 'checker'.
func NewServer(rpcType RpcType, checker ParamChecker, onReq OnRequest) (*Server, error) {
	s := &Server{
		rpcSvr:        rpc.NewServer(),
		rpcType:       rpcType,
		paramChecker:  checker,
		onRequestUser: onReq,
	}

	// Register
	if s.rpcType == RpcTypeJSON {
		// In json rpc, types can be registered only once,
		// no matter how many servers are started.
		if atomic.LoadInt32(&jsonRpcRegisterDone) == 0 {
			atomic.StoreInt32(&jsonRpcRegisterDone, 1)
			if err := rpc.Register((*Svr)(s)); err != nil {
				return nil, err
			}
		}
	} else if s.rpcType == RpcTypeGOB {
		if err := s.rpcSvr.Register((*Svr)(s)); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Svr) OnRequestInternal(in Request, out *Reply) error {
	/*if err := s.paramChecker.VerifyIn(in.Func, in); err != nil {
		return err
	}*/
	if err := s.onRequestUser(in, out); err != nil {
		return err
	}
	if err := s.paramChecker.VerifyOut(in.Func, out, false); err != nil {
		return err
	}
	return nil
}

func (s *Server) Serve(conn net.Conn) {
	if s.rpcType == RpcTypeJSON {
		jsonrpc.ServeConn(conn)
	} else if s.rpcType == RpcTypeGOB {
		s.rpcSvr.ServeConn(conn)
	}
}

func (s *Server) Close() error {
	return nil
	//return s.netSvr.Close()
}
