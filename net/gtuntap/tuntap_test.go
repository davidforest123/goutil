package gtuntap

import (
	"fmt"
	"goutil/net/gnet"
	"net"
	"testing"
)

func TestIface_Read(t *testing.T) {
	deviceIPNet, err := gnet.ParseIPNet("192.168.23.1/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	iface := New("tun", "", 1300, deviceIPNet.Std())
	if err := iface.SudoSetup(); err != nil {
		fmt.Println(err)
		return
	}

	routeIPNet := net.IPNet{
		IP:   net.ParseIP("192.168.13.0"),
		Mask: gnet.MustParseIPMask("255.255.255.0").Std(),
	}
	err = gnet.SudoAddNetRoute(iface.Name(), routeIPNet, routeIPNet.IP.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("setup finished")

	packet := make([]byte, 2000)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Printf("Packet Received: % x\n", packet[:n]))
	}
}
