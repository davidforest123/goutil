package gnet

/**
IPV4的保留地址汇总
0.0.0.0/8：用于广播信息到当前主机
10.0.0.0/8：用于专用网络中的本地通信
100.64.0.0/10：用于在电信级NAT环境中服务提供商与其用户通信
127.0.0.0/8：用于到本地主机的环回地址
169.254.0.0/16：用于单链路的两个主机之间的链路本地地址，而没有另外指定IP地址，例如通常从DHCP服务器所检索到的IP地址
172.16.0.0/12：用于专用网络中的本地通信
192.0.0.0/24：用于IANA的IPv4特殊用途地址表
192.0.2.0/24：分配为用于文档和示例中的“TEST-NET”（测试网），它不应该被公开使用
192.88.99.0/24：用于6to4任播中继，已废弃
192.168.0.0/16：用于专用网络中的本地通信
198.18.0.0/15：用于测试两个不同的子网的网间通信
198.51.100.0/24：分配为用于文档和示例中的“TEST-NET-2”（测试-网-2），它不应该被公开使用
203.0.113.0/24：分配为用于文档和示例中的“TEST-NET-3”（测试-网-3），它不应该被公开使用
224.0.0.0/4：用于多播
240.0.0.0/4：用于将来使用
255.255.255.255/32：用于受限广播地址
*/

import (
	"encoding/binary"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gconv"
	"net"
	"strings"
)

type (
	// IP address.
	IP net.IP

	// IPNet defines IP network, or IP range.
	// Notice:
	// Valid IPNet samples: 192.168.7.0/24
	// Invalid IPNet samples: 192.168.7.123/24
	IPNet net.IPNet
)

var (
	addrSpaceOfLoopBackIPv4 = [...]string{
		"127.0.0.0/8",
	}

	addrSpaceOfLoopBackIPv6 = [...]string{
		"::1/128",
	}

	addrSpaceOfAnyIPv6 = [...]string{
		"::/128",
	}

	addrSpacesOfLANIPv4 = [...]string{
		"10.0.0.0/8",     // 10.*.*.*
		"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
		"192.168.0.0/16", // 192.168.*.*
	}

	addrSpacesOfLANIPv6 = [...]string{
		"fd00::/8",
	}
)

func WrapIP(ip net.IP) IP {
	return IP(ip)
}

func ParseIP(s string) (IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, gerrors.New("Invalid IP address string '" + s + "'")
	}
	return IP(ip), nil
}

func (ip IP) Raw() net.IP {
	return net.IP(ip)
}

func (ip IP) String() string {
	return net.IP(ip).String()
}

func (ip IP) IsLoopBack() bool {
	for _, it := range addrSpaceOfLoopBackIPv4 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	for _, it := range addrSpaceOfLoopBackIPv6 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	return false
}

// To4 converts the IPv4 address ip to a 4-byte representation.
// IPv4 maybe takes 16-byte buffer, maybe not, that's why we need To4() function.
// If ip is not an IPv4 address, To4 returns nil.
func (ip IP) To4() IP {
	return WrapIP(ip.Raw().To4())
}

func (ip IP) IsV4() bool {
	return ip.Raw().To4() != nil
}

func (ip IP) IsV6() bool {
	return !ip.IsV4()
}

func (ip IP) IsPublic() bool {
	return !ip.IsPrivate()
}

func (ip IP) IsPrivate() bool {
	return ip.Raw().IsPrivate()
	/*for _, it := range addrSpacesOfLANIPv4 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	for _, it := range addrSpacesOfLANIPv6 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	return false*/
}

func (ip IP) IsAny() bool {
	if ip.String() == "0.0.0.0" {
		return true
	}

	for _, it := range addrSpaceOfLoopBackIPv6 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	return false
}

func IsIPString(s string) bool {
	_, err := ParseIP(s)
	return err == nil
}

func LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}

// GetAllPrivateIPs returns all my local IPs.
func GetAllPrivateIPs() ([]net.IP, error) {
	var result = make([]net.IP, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			if WrapIP(ip).IsPrivate() {
				result = append(result, ip)
			}
		}
	}

	return result, nil
}

func WrapIPNet(ipNet net.IPNet) IPNet {
	return IPNet(ipNet)
}

func WrapIPNetPtr(ipNet *net.IPNet) *IPNet {
	return (*IPNet)(ipNet)
}

// ParseIPNet casts string like `192.168.1.3/24` into IPNet.
func ParseIPNet(s string) (IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return IPNet{}, err
	}
	// ipNet returned by net.ParseCIDR has been removed Host IP info, just Net info left in it.
	// For example, net.Parse(`192.168.1.3/24`) returns ip `192.168.1.3` and ipNet `192.168.1.0/24`.
	// So I'll add the host IP back in.
	ipNet.IP = ip
	return WrapIPNet(*ipNet), nil
}

func (in IPNet) Std() net.IPNet {
	return net.IPNet(in)
}

// IsHostType returns if a IPNet is a host ip address.
// "192.168.3.3/24": true
// "192.168.3.0/24": false
// Notice: IPNet ending with `/32` is both a Host type and a Net type.
func (in IPNet) IsHostType() bool {
	return in.String() != in.OnlyKeepNet().String()
}

// IsNetType returns if a IPNet is a network.
// "192.168.3.0/24": true
// "192.168.3.3/24": false
// Notice: IPNet ending with `/32` is both a Host type and a Net type.
func (in IPNet) IsNetType() bool {
	return in.String() == in.OnlyKeepNet().String()
}

// OnlyKeepNet returns pure network format IPNet, drop single host info.
// IPNet: `192.168.1.2/24`
// Return: `192.168.1.0/24`
// FIXME 需要查资料验证结果是否正确
func (in IPNet) OnlyKeepNet() IPNet {
	// implement 1
	// ParseCIDR returns net.IP and net.IPNet, in the net.IPNet, it keeps only network info
	/*
		_, ipNet, err := net.ParseCIDR(in.String())
		if err != nil {
			panic(err)
		}
		return WrapIPNet(*ipNet)*/

	// implement 2
	std := in.Std()
	// IPv4 net.IP does not necessarily occupy 4 bytes, it may occupy 16 bytes,
	// call To4() can ensure that it occupies 4 bytes
	newIP := std.IP.To4()
	for i, op := range in.Mask {
		newIP[i] = newIP[i] & op
	}
	return IPNet{
		IP:   newIP,
		Mask: in.Mask,
	}
}

func (in IPNet) Contains(ip net.IP) bool {
	std := in.Std()
	return std.Contains(ip)
}

// StdString returns standard format IPNet.
// IPNet: `192.168.1.2/24`
// Return: `192.168.1.2/24`
func (in IPNet) String() string {
	return in.CidrStyleString()
}

// CidrStyleString returns CIDR format IPNet.
// IPNet: `192.168.1.2/24`
// Return: `192.168.1.2/24`
func (in IPNet) CidrStyleString() string {
	std := in.Std()
	return std.String()
}

// MacStyleString returns mac style network format IPNet.
// IPNet: `192.168.1.2/24`
// Return: `192.168.1`
// FIXME 需要查资料验证结果是否正确
func (in IPNet) MacStyleString() string {
	// IPv4 net.IP does not necessarily occupy 4 bytes, it may occupy 16 bytes,
	// call To4() can ensure that it occupies 4 bytes
	ipv4 := in.Std().IP.To4()
	var result []byte
	for i, op := range in.Mask {
		b := ipv4[i] & op
		if int(b) != 0 {
			result = append(result, b)
		}
	}
	return gconv.BytesToDecString(result, ".", 0)
}

// Verify checks if IPNet is a valid IP network or not.
// Valid IPNet samples: 192.168.7.0/24
// Invalid IPNet samples: 192.168.7.12/24
func (in IPNet) Verify() error {
	_, parsedIpNet, err := net.ParseCIDR(in.String())
	if err != nil {
		return err
	}
	if parsedIpNet.String() != in.String() {
		return gerrors.New("IPNet(%s) is not a valid IP network", in.String())
	}
	return nil
}

// ListAll returns all IPs of current IP network.
// FIXME: IPv6 not supported for now.
func (in IPNet) ListAll() []net.IP {
	var result []net.IP
	num := binary.BigEndian.Uint32(in.IP)
	mask := binary.BigEndian.Uint32(in.Mask)
	num &= mask
	for mask < 0xffffffff {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], num)
		result = append(result, buf[:])
		mask += 1
		num += 1
	}
	return result
}
