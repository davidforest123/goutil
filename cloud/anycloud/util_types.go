package anycloud

import (
	"fmt"
	"goutil/net/gsniff"
	"goutil/net/gstun"
	"net"
)

type (
	AnyStoreServe func(conn net.Conn)

	StreamUsageMod string

	// StreamUsage
	// Stream usage format: TurnNode|ModName|ModTarget
	// examples:
	// |modsocks|tcp://google.com"
	// Any1BCJFVHjhcDkEdOCMKFnfdD|modsocks|tcp://google.com"
	StreamUsage struct {
		TurnNode  string
		ModName   StreamUsageMod
		ModTarget string
	}

	AnyCloudMgr struct {
		NodeID   string
		NodeName string
		AnyNet   struct {
			InternalAddr   []net.IP
			Nat            gstun.Nat
			ConnectedNodes map[string]string // map[NodeID]NodeName
		}
		AnyStore struct {
		}
		AnyServer struct {
		}
	}
)

var (
	// MfsUsageAnyNetRaw ...
	// Mfs connection usages.
	MfsUsageAnyNetRaw    = gsniff.MustMakeUsage("@antraw\n")
	MfsUsageAnyNetStun   = gsniff.MustMakeUsage("@antstn\n")
	MfsUsageAnyNetRpc    = gsniff.MustMakeUsage("@antrpc\n")
	MfsUsageAnyNetMQ     = gsniff.MustMakeUsage("@antmsq\n")
	MfsUsageAnyStoreRpc  = gsniff.MustMakeUsage("@aystrc:") // ':' means there is more
	MfsUsageAnyServerRpc = gsniff.MustMakeUsage("@aysvrc\n")

	// StreamUsageModVPN ...
	// Accepted|Built-in stream usage mod names.
	StreamUsageModVPN          = makeStreamUsageMod("modvpn")
	StreamUsageModSocks        = makeStreamUsageMod("modsocks")
	StreamUsageModTunnel       = makeStreamUsageMod("modtunnel")
	StreamUsageModMsg          = makeStreamUsageMod("modmsg")
	StreamUsageModAnyStoreRpc  = makeStreamUsageMod("modanystorerpc")
	StreamUsageModAnyServerRpc = makeStreamUsageMod("modanyserverrpc")
	StreamUsageModMQ           = makeStreamUsageMod("modmq") // anyNet: ping, connect-info-update, bootstraps-query, anyStore: stat, anyServer: stat

	AllStreamUsageMods []string
)

func MakeStreamUsage(turnNode string, modName StreamUsageMod, modTarget string) *StreamUsage {
	res := StreamUsage{
		TurnNode:  turnNode,
		ModName:   modName,
		ModTarget: modTarget,
	}
	return &res
}

func makeStreamUsageMod(mod string) StreamUsageMod {
	AllStreamUsageMods = append(AllStreamUsageMods, mod)
	return StreamUsageMod(mod)
}

func (cs *StreamUsage) String() string {
	return fmt.Sprintf("%s|%s|%s", cs.TurnNode, cs.ModName, cs.ModTarget)
}
