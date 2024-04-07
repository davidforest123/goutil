package anycloud

import (
	"goutil/net/grpcs"
)

type (
	AnyCloud struct {
		network   string
		addr      string
		accessKey string
	}
)

func New(network, addr, accessKey string) (*AnyCloud, error) {
	return &AnyCloud{
		network:   network,
		addr:      addr,
		accessKey: accessKey,
	}, nil
}

// AnyNet connects to AnyNet system.
func (ac *AnyCloud) AnyNet() (*AnyNet, error) {
	return &AnyNet{ac: ac}, nil
}

// AnyStore connects to AnyStore node.
func (ac *AnyCloud) AnyStore(nodeID string) (*AnyStore, error) {
	rpc, err := grpcs.Dial(grpcs.RpcTypeJSON, ac.network, ac.addr, NewRPCChecker(),
		[]byte(MfsUsageAnyStoreRpc.String()+nodeID+"\n"))
	if err != nil {
		return nil, err
	}

	c := &AnyStore{ac: ac, rpc: rpc}
	return c, nil
}

// AnyServer connects to AnyServer system.
func (ac *AnyCloud) AnyServer() (*AnyServer, error) {
	return &AnyServer{}, nil
}
