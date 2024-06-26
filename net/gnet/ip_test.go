package gnet

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gtest"
	"net"
	"strings"
	"testing"
)

/*func TestCheckIPString(t *testing.T) {
	ip4lanlist := []string{"10.78.1.2", "192.168.1.12"}
	for _, v := range ip4lanlist {
		if CheckIPString(v) != IPv4_LAN {
			t.Error(fmt.Sprintf("%s test failed", v))
		}
	}

	ip4loopbacklist := []string{"127.0.0.1"}
	for _, v := range ip4loopbacklist {
		if CheckIPString(v) != IPv4_LOOPBACK {
			t.Error(fmt.Sprintf("%s test failed", v))
		}
	}

	ip4anylist := []string{"0.0.0.0"}
	for _, v := range ip4anylist {
		if CheckIPString(v) != IPv4_ANY {
			t.Error(fmt.Sprintf("%s test failed", v))
		}
	}
}*/

func TestIPNet_ListAll(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("192.168.3.139/24")
	if err != nil {
		t.Error(err)
		return
	}
	expect := "192.168.3.0,192.168.3.1,192.168.3.2,192.168.3.3,192.168.3.4,192.168.3.5,192.168.3.6,192.168.3.7," +
		"192.168.3.8,192.168.3.9,192.168.3.10,192.168.3.11,192.168.3.12,192.168.3.13,192.168.3.14,192.168.3.15," +
		"192.168.3.16,192.168.3.17,192.168.3.18,192.168.3.19,192.168.3.20,192.168.3.21,192.168.3.22,192.168.3.23," +
		"192.168.3.24,192.168.3.25,192.168.3.26,192.168.3.27,192.168.3.28,192.168.3.29,192.168.3.30,192.168.3.31," +
		"192.168.3.32,192.168.3.33,192.168.3.34,192.168.3.35,192.168.3.36,192.168.3.37,192.168.3.38,192.168.3.39," +
		"192.168.3.40,192.168.3.41,192.168.3.42,192.168.3.43,192.168.3.44,192.168.3.45,192.168.3.46,192.168.3.47," +
		"192.168.3.48,192.168.3.49,192.168.3.50,192.168.3.51,192.168.3.52,192.168.3.53,192.168.3.54,192.168.3.55," +
		"192.168.3.56,192.168.3.57,192.168.3.58,192.168.3.59,192.168.3.60,192.168.3.61,192.168.3.62,192.168.3.63," +
		"192.168.3.64,192.168.3.65,192.168.3.66,192.168.3.67,192.168.3.68,192.168.3.69,192.168.3.70,192.168.3.71," +
		"192.168.3.72,192.168.3.73,192.168.3.74,192.168.3.75,192.168.3.76,192.168.3.77,192.168.3.78,192.168.3.79," +
		"192.168.3.80,192.168.3.81,192.168.3.82,192.168.3.83,192.168.3.84,192.168.3.85,192.168.3.86,192.168.3.87," +
		"192.168.3.88,192.168.3.89,192.168.3.90,192.168.3.91,192.168.3.92,192.168.3.93,192.168.3.94,192.168.3.95," +
		"192.168.3.96,192.168.3.97,192.168.3.98,192.168.3.99,192.168.3.100,192.168.3.101,192.168.3.102,192.168.3.103," +
		"192.168.3.104,192.168.3.105,192.168.3.106,192.168.3.107,192.168.3.108,192.168.3.109,192.168.3.110,192.168.3.111," +
		"192.168.3.112,192.168.3.113,192.168.3.114,192.168.3.115,192.168.3.116,192.168.3.117,192.168.3.118,192.168.3.119," +
		"192.168.3.120,192.168.3.121,192.168.3.122,192.168.3.123,192.168.3.124,192.168.3.125,192.168.3.126,192.168.3.127," +
		"192.168.3.128,192.168.3.129,192.168.3.130,192.168.3.131,192.168.3.132,192.168.3.133,192.168.3.134,192.168.3.135," +
		"192.168.3.136,192.168.3.137,192.168.3.138,192.168.3.139,192.168.3.140,192.168.3.141,192.168.3.142,192.168.3.143," +
		"192.168.3.144,192.168.3.145,192.168.3.146,192.168.3.147,192.168.3.148,192.168.3.149,192.168.3.150,192.168.3.151," +
		"192.168.3.152,192.168.3.153,192.168.3.154,192.168.3.155,192.168.3.156,192.168.3.157,192.168.3.158,192.168.3.159," +
		"192.168.3.160,192.168.3.161,192.168.3.162,192.168.3.163,192.168.3.164,192.168.3.165,192.168.3.166,192.168.3.167," +
		"192.168.3.168,192.168.3.169,192.168.3.170,192.168.3.171,192.168.3.172,192.168.3.173,192.168.3.174,192.168.3.175," +
		"192.168.3.176,192.168.3.177,192.168.3.178,192.168.3.179,192.168.3.180,192.168.3.181,192.168.3.182,192.168.3.183," +
		"192.168.3.184,192.168.3.185,192.168.3.186,192.168.3.187,192.168.3.188,192.168.3.189,192.168.3.190,192.168.3.191," +
		"192.168.3.192,192.168.3.193,192.168.3.194,192.168.3.195,192.168.3.196,192.168.3.197,192.168.3.198,192.168.3.199," +
		"192.168.3.200,192.168.3.201,192.168.3.202,192.168.3.203,192.168.3.204,192.168.3.205,192.168.3.206,192.168.3.207," +
		"192.168.3.208,192.168.3.209,192.168.3.210,192.168.3.211,192.168.3.212,192.168.3.213,192.168.3.214,192.168.3.215," +
		"192.168.3.216,192.168.3.217,192.168.3.218,192.168.3.219,192.168.3.220,192.168.3.221,192.168.3.222,192.168.3.223," +
		"192.168.3.224,192.168.3.225,192.168.3.226,192.168.3.227,192.168.3.228,192.168.3.229,192.168.3.230,192.168.3.231," +
		"192.168.3.232,192.168.3.233,192.168.3.234,192.168.3.235,192.168.3.236,192.168.3.237,192.168.3.238,192.168.3.239," +
		"192.168.3.240,192.168.3.241,192.168.3.242,192.168.3.243,192.168.3.244,192.168.3.245,192.168.3.246,192.168.3.247," +
		"192.168.3.248,192.168.3.249,192.168.3.250,192.168.3.251,192.168.3.252,192.168.3.253,192.168.3.254"
	ipList := WrapIPNetPtr(ipNet).ListAll()
	var ipArr []string
	for _, v := range ipList {
		ipArr = append(ipArr, v.String())
	}
	if strings.Join(ipArr, ",") != expect {
		t.Error("IPNet.ListAll error")
	}
}

func TestParseIPNet(t *testing.T) {
	ipNet, err := ParseIPNet("192.168.1.3/24")
	gtest.Assert(t, err)
	fmt.Println(ipNet.String())
}

func TestIPNet_IsHostType(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("172.31.16.1/20").Expect(true)
	cl.New().Input("172.31.16.0/20").Expect(false)
	cl.New().Input("192.168.1.1/24").Expect(true)
	cl.New().Input("192.168.1.0/24").Expect(false)

	for _, v := range cl.Get() {
		ipNetStr := v.Inputs[0].(string)
		expect := v.Expects[0].(bool)

		ipNet, err := ParseIPNet(ipNetStr)
		gtest.Assert(t, err)
		result := ipNet.IsHostType()
		if result != expect {
			gtest.PrintlnExit(t, "IPNet(%s).IsHostType() expect %v but got %v", ipNetStr, expect, result)
		}
	}
}

func TestIPNet_Contains(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("172.31.16.0/20").Input("172.31.16.1").Expect(true)

	for _, v := range cl.Get() {
		ipNetStr := v.Inputs[0].(string)
		ipStr := v.Inputs[1].(string)
		expect := v.Expects[0].(bool)

		ipNet, err := ParseIPNet(ipNetStr)
		gtest.Assert(t, err)
		ip, err := ParseIP(ipStr)
		gtest.Assert(t, err)
		result := ipNet.Contains(ip.Raw())
		if result != expect {
			gtest.PrintlnExit(t, "IPNet(%s).Contains(%s) expect %v but got %v", ipNetStr, ipStr, expect, result)
		}
	}
}

func TestIPNet_String(t *testing.T) {
	ipNet := WrapIPNet(net.IPNet{
		IP:   net.ParseIP("192.168.1.2"),
		Mask: MustParseIPMask("255.255.255.0").Std(),
	})
	fmt.Println(ipNet.String())
}

func TestIPNet_CidrStyleNetString(t *testing.T) {
	ipNet := WrapIPNet(net.IPNet{
		IP:   net.ParseIP("192.168.1.2"),
		Mask: MustParseIPMask("255.255.255.0").Std(),
	})
	fmt.Println(ipNet.OnlyKeepNet().String())

	// FIXME 需要查资料验证结果是否正确
	ipNet = WrapIPNet(net.IPNet{
		IP:   net.ParseIP("192.168.23.2"),
		Mask: MustParseIPMask("255.255.254.0").Std(),
	})
	fmt.Println(ipNet.OnlyKeepNet().String())

	// FIXME 需要查资料验证结果是否正确
	ipNet = WrapIPNet(net.IPNet{
		IP:   net.ParseIP("192.168.23.2"),
		Mask: MustParseIPMask("255.255.252.0").Std(),
	})
	fmt.Println(ipNet.OnlyKeepNet().String())
}

func TestIPNet_MacStyleString(t *testing.T) {
	// FIXME 需要查资料验证结果是否正确
	ipNet := WrapIPNet(net.IPNet{
		IP:   net.ParseIP("192.168.23.2"),
		Mask: MustParseIPMask("255.255.252.0").Std(),
	})
	fmt.Println(ipNet.MacStyleString())
}
