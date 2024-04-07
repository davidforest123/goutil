package grpcs

import (
	"github.com/davidforest123/goutil/net/gnet"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// RPC over SSL/TLS:
// https://gist.github.com/artyom/6897140

type (
	Client struct {
		rpcCli  *rpc.Client
		checker ParamChecker
	}

	AsyncCallback func(s string, b []byte)
)

// Dial connects to RPC server.
// featureBytes: send it before doing anything of RPC protocol, this buf is used for remote server to distinguish
// between different client connections for different purposes on the same port.
func Dial(rpcType RpcType, network, address string, checker ParamChecker, featureBytes []byte) (*Client, error) {
	netCli, err := gnet.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if len(featureBytes) > 0 {
		if _, err = netCli.Write(featureBytes); err != nil {
			return nil, err
		}
	}
	res := &Client{
		checker: checker,
	}
	if rpcType == RpcTypeJSON { // json rpc
		res.rpcCli = jsonrpc.NewClient(netCli)
		if err != nil {
			return nil, err
		}
	} else if rpcType == RpcTypeGOB { // gob rpc
		res.rpcCli = rpc.NewClient(netCli)
	}
	return res, nil
}

func (c *Client) Call(name string, args Request, reply *Reply) error {
	if err := c.checker.VerifyIn(name, args); err != nil {
		return err
	}

	// nil Reply means doesn't need output, but rpcCli.Call requires valid output structure.
	//isOriginReplyNil := reply == nil
	if reply == nil {
		tmp := NewReply()
		reply = &tmp
	}

	args.Func = name
	if err := c.rpcCli.Call("Svr.OnRequestInternal", args, reply); err != nil {
		return err
	}

	return nil
	/*
		if isOriginReplyNil {
			return nil
		}
		return c.checker.VerifyOut(name, reply, true)*/
}

func (c *Client) Close() error {
	return c.rpcCli.Close()
}
