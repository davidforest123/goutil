package anycloud

import (
	"goutil/net/gssh"
)

type (
	AnyServer struct {
	}
)

func (s *AnyServer) GlobalTopology() {

}

func (s *AnyServer) Ssh(targetNodeID, username, password, privateKeyFile, passphrase string) (*gssh.Client, error) {
	return nil, nil
}
