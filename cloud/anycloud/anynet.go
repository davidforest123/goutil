package anycloud

import (
	"goutil/net/gmsg"
	"net"
)

type (
	AnyNet struct {
		ac *AnyCloud
	}
)

func (an *AnyNet) MQ(targetNodeID string) (*MQ, error) {
	streamConn, err := net.Dial(an.ac.network, an.ac.addr)
	if err != nil {
		return nil, err
	}
	if _, err = streamConn.Write([]byte(MfsUsageAnyNetRpc.String() + targetNodeID + "\n")); err != nil {
		return nil, err
	}
	msgConn, err := gmsg.ClientConn(streamConn, "")
	if err != nil {
		return nil, err
	}

	c := &MQ{anyNet: an, msgConn: msgConn}
	go c.rcvMsg()
	return c, nil
}

func (an *AnyNet) Gossip() (*Gossip, error) {
	return nil, nil
}
