package gstun

import (
	"github.com/ccding/go-stun/stun"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/net/gnet"
)

// https://www.cnblogs.com/xinbigworld/p/16297220.html 什么情况下可以进行NAT？

type (
	NatType      string
	BehaviorType string

	Nat struct {
		ExternalAddr  []gnet.IPPort
		NatType       NatType
		MappingType   BehaviorType
		FilteringType BehaviorType
	}
)

var (
	NatTypeUnknown            = NatType("NatUnknown")
	NatTypeBlocked            = NatType("NatBlocked")
	NatTypeNotBehind          = NatType("NatNotBehind")
	NatTypeSymmetric          = NatType("NatSymmetric")
	NatTypeFullCone           = NatType("NatFullCone")
	NatTypeAddrRestrictedCone = NatType("NatAddrRestrictedCone")
	NatTypePortRestrictedCone = NatType("NatPortRestrictedCone")

	BehaviorTypeUnknown     = BehaviorType("BehaviorUnknown")
	BehaviorTypeEndpoint    = BehaviorType("BehaviorEndpointIndependent")
	BehaviorTypeAddr        = BehaviorType("BehaviorAddressDependent")
	BehaviorTypeAddrAndPort = BehaviorType("BehaviorAddressAndPortDependent")

	// copied from https://github.com/pradt2/always-online-stun/blob/master/valid_hosts.txt
	publicStunServers = []string{
		"stun.l.google.com:19302",
		"stun.l.google.com:19305",
		"stun1.l.google.com:19302",
		"stun1.l.google.com:19305",
		"stun2.l.google.com:19302",
		"stun2.l.google.com:19305",
		"stun3.l.google.com:19302",
		"stun3.l.google.com:19305",
		"stun4.l.google.com:19302",
		"stun4.l.google.com:19305",
		"stun.tel.lu:3478",
		"stun.voip.blackberry.com:3478",
		"stun.nextcloud.com:443",
		"stun.vo.lu:3478",
	}
)

func convertBehavior(b stun.BehaviorType) (BehaviorType, error) {
	switch b {
	case stun.BehaviorTypeUnknown:
		return BehaviorTypeUnknown, nil
	case stun.BehaviorTypeAddr:
		return BehaviorTypeAddr, nil
	case stun.BehaviorTypeAddrAndPort:
		return BehaviorTypeAddrAndPort, nil
	case stun.BehaviorTypeEndpoint:
		return BehaviorTypeEndpoint, nil
	default:
		return BehaviorTypeUnknown, gerrors.New("unknown behavior %d", b)
	}
}

func Discover(stunServer string) (*Nat, error) {
	// Creates a STUN client. NewClientWithConnection can also be used if you want to handle the UDP listener by yourself.
	client := stun.NewClient()
	client.SetServerAddr(stunServer)
	client.SetVerbose(false)
	client.SetVVerbose(false)

	result := &Nat{}

	// Test behavior.
	natBehavior, err := client.BehaviorTest()
	if err != nil && err.Error() != "Not behind a NAT" {
		return nil, err
	}
	if natBehavior != nil {
		result.MappingType, err = convertBehavior(natBehavior.MappingType)
		if err != nil {
			return nil, err
		}
		result.FilteringType, err = convertBehavior(natBehavior.FilteringType)
		if err != nil {
			return nil, err
		}
	}

	// Discover the NAT.
	nat, host, err := client.Discover()
	if err != nil {
		return nil, err
	}
	ipPort, err := gnet.ParseIPPort(host.String())
	if err != nil {
		return nil, err
	}
	result.ExternalAddr = append(result.ExternalAddr, *ipPort)
	switch nat {
	case stun.NATUnknown:
		result.NatType = NatTypeUnknown
	case stun.NATNone:
		result.NatType = NatTypeNotBehind
	case stun.NATBlocked:
		result.NatType = NatTypeBlocked
	case stun.NATFull:
		result.NatType = NatTypeFullCone
	case stun.NATSymmetric:
		result.NatType = NatTypeSymmetric
	case stun.NATRestricted:
		result.NatType = NatTypeAddrRestrictedCone
	case stun.NATPortRestricted:
		result.NatType = NatTypePortRestrictedCone
	case stun.SymmetricUDPFirewall:
		result.NatType = NatTypeSymmetric
	}

	return result, nil
}
