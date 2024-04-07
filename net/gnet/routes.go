package gnet

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gstring"
	"github.com/davidforest123/goutil/sys/gcmd"
	"github.com/libp2p/go-netroute"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Notice
//
// 1, TAP is not supported for now.
//
// 2, macOS上，把TUN设备的IP（比如192.168.23.1）作为Gateway IP添加进路由规则时，TUN必须提前设置好然后再执行路由添加，
//    否则这条路由的Gateway尽管还是设置的那个IP（比如192.168.23.1），但是路由规则的关联的网卡（Netif）就变成默认物理网卡而不是TUN设备了
//
//

/**
3种路由类型

1、主机路由
主机路由是路由选择表中指向单个 IP 地址或主机名的路由记录。主机路由的 Flags 字段为 H。例如，在下面的示例中，
本地主机通过 IP 地址 192.168.1.1 的路由器到达 IP 地址为 10.0.0.10 的主机。
Destination    Gateway       Genmask        Flags     Metric    Ref    Use    Iface
-----------    -------       -------        -----     ------    ---    ---    -----
10.0.0.10     192.168.1.1    255.255.255.255   UH          0      0      0     eth0

2、网络路由
网络路由是代表主机可以到达的网络。网络路由的 Flags 字段为 N。例如，在下面的示例中，本地主机将发送到网络 192.19.12
的数据包转发到 IP 地址为 192.168.1.1 的路由器。
Destination    Gateway       Genmask      Flags    Metric    Ref     Use    Iface
-----------    -------       -------      -----    -----     ---     ---    -----
192.19.12     192.168.1.1    255.255.255.0   UN        0       0       0     eth0

3、默认路由
默认路由属于替补路由，优先级最低，当其他路由不可达时，换言之当主机不能在路由表中查找到目标主机的 IP 地址或网络路由时，
数据包就被发送到默认路由（默认网关）上。默认路由的 Flags 字段为 G。例如，在下面的示例中，默认路由是 IP 地址为 192.168.1.1 的路由器。
Destination    Gateway       Genmask    Flags     Metric    Ref    Use    Iface
-----------    -------       -------    -----     ------    ---    ---    -----
default       192.168.1.1    0.0.0.0       UG          0      0      0     eth0
*/

/**
路由标志位（Flags）
U			`up`                       // 路由是活动的
H			`dest_is_single_host`      // 目标是一个主机
G			`gateway`                  // 路由指向网关
S			`static`                   // bsd
Cloned		`clone_based_on_route`     // bsd
W 			`clone_auto_local`         // bsd
L			`link_to_hw`               // bsd
Reinsta		`reinstate_route`          // linux
D 			`dynamic_installed`        // linux 由路由的后台程序动态地安装
M			`modified_from_routing_sw` // linux 由路由的后台程序修改
N           `dest_is_network`		   // linux 目标是个网络，即目标IP和目标子网掩码代表的所有主机
A			`installed_by_addrconf`    // linux
Cached		`cached`                   // linux
Rejected	`rejected`                 // linux `!` 拒绝路由
*/

/**
获取某个目的IP的路由（已验证macOS）
bsd：route get $DEST-IP
*/

/**
Linux下的route 命令？亲测和macOS上不符合，比如macOS上没有gw和dev参数，等等
设置和查看路由表都可以用 route 命令，设置内核路由表的命令格式是：
# route  [add|del] [-net|-host] $TARGET [netmask $NM] [gw $GW] [[dev] $IFACE]
其中：
add : 		添加一条路由规则
del : 		删除一条路由规则
-net : 		目的地址是一个网络
-host : 	目的地址是一个主机
$TARGET : 	目的网络或主机
netmask : 	目的地址的网络掩码
gw : 		路由数据包通过的网关
dev : 		为路由指定的网络接口

参考： http://www.manongjc.com/detail/27-zqgfyhucblawzff.html   https://www.cyberciti.biz/faq/linux-ip-command-examples-usage-syntax/



网络路由（已验证macOS）
macOS Monterey上，$MAC-STYLE-DEST-IPNet不可以使用192.168.13.0/32这种形式，亲测无效，应该使用192.168.13这种格式，其他版本未知
bsd添加方法1： route add -net $MAC-STYLE-DEST-IPNet(比如192.168.13) $GATEWAY-IP
bsd添加方法2： route add -net $MAC-STYLE-DEST-IPNet(比如192.168.13) -interface $DEV-NAME(比如utun6)
bsd添加方法3： route add -net $DEST-IP(网段比如192.168.13.0，最后一位数字可以随便填写，比如1、2、3，只要子网掩码对就行) -netmask $NETMASK $GATEWAY-IP
bsd添加方法4（推荐）： route add -net $DEST-IP(网段比如192.168.13.0，最后一位数字可以随便填写，比如1、2、3，只要子网掩码对就行) -netmask $NETMASK -interface $DEV-NAME(比如utun6)
bsd删除方法1： route delete -net $MAC-STYLE-DEST-IPNet(比如192.168.13)
bsd删除方法2： route delete -net $MAC-STYLE-DEST-IPNet(比如192.168.13) -interface $DEV-NAME(比如utun6)
bsd删除方法3： route delete -net $DEST-IP(网段比如192.168.13.0，最后一位数字不能随便写，如果子网掩码是255.255.255.0就只能写0) -gateway $GATEWAY-IP
bsd删除方法4（推荐）： route delete -net $DEST-IP(网段比如192.168.13.0，最后一位数字不能随便写，如果子网掩码是255.255.255.0就只能写0) -interface $DEV-NAME(比如utun6)

默认路由（未验证）
bsd: route add default $GATEWAY-IP   or   route add -net 0.0.0.0 $GATEWAY-IP
bsd: route delete default $GATEWAY-IP
bsd example: route -n delete default -link 8` // delete default route line#8

主机路由（暂无需要）
*/

// TODO handle ipv6 in future

var (
	errUnsupportedOS = gerrors.New("unsupported OS %s", runtime.GOOS)
)

type (
	RouteRule struct {
		Type    string // net/host/default
		IPVer   int    // 4 or 6
		Dest    string // 目标网段或者主机
		Gateway string // 网关地址，”*” 表示目标是本主机所属的网络，不需要路由
		Flags   string
		Iface   string // 该路由表项对应的输出接口
		Expire  *int
	}
)

// GetRoute returns route info of some specified destination IP, like Interface and gateway.
// On macOS, this function is the functional equivalent of the `route get $DEST-IP` command.
func GetRoute(dst net.IP) (iface *net.Interface, gateway net.IP, err error) {
	r, err := netroute.New()
	if err != nil {
		return nil, nil, err
	}
	iface, gw, _, err := r.Route(dst)
	return iface, gw, err
}

func GetRoutes() ([]RouteRule, error) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "netstat -rn"
	case "linux":
		cmd = "ip route"
	case "freebsd":
		cmd = "route -n"
	case "windows":
		cmd = "route print"
	default:
		return nil, errUnsupportedOS
	}
	out, err := gcmd.RunScript(cmd)
	if err != nil {
		return nil, err
	}
	lines := gstring.SplitToLines(string(out))
	lines = gstring.RemoveEntirelySpaces(lines)

	var result []RouteRule

	ipVer := 4
	for _, line := range lines {
		if line == "Internet:" {
			ipVer = 4
			continue
		}
		if line == "Internet6:" {
			ipVer = 6
			continue
		}
		sarr := strings.Fields(line)
		if len(sarr) < 4 {
			continue
		}
		if sarr[0] == "Destination" {
			continue
		}

		newRule := RouteRule{
			IPVer:   ipVer,
			Dest:    sarr[0],
			Gateway: sarr[1],
			Flags:   sarr[2],
			Iface:   sarr[3],
			Expire:  nil,
		}
		if len(sarr) == 5 {
			expire, err := strconv.Atoi(sarr[4])
			if err == nil {
				newRule.Expire = &expire
			}
		}

		result = append(result, newRule)
	}

	return result, nil
}

// SudoAddNetRoute adds `net` type route rule to local operating system
// networksetup -setadditionalroutes Ethernet 192.168.1.0 255.255.255.0 10.0.0.2 persistent
// address: 不知道干嘛用的
//
// Notice:
// 1、如果程序进程退出，包括意外退出，通常TUN设备会被系统清除，TUN设备相关的路由也会被系统清除。macOS下亲测如此
func SudoAddNetRoute(iface string, destIPNet net.IPNet, address string) error {
	switch runtime.GOOS {
	case "darwin": // verified
		out, err := gcmd.RunScript("sudo route get " + destIPNet.IP.String())
		if err != nil {
			return err
		}
		if !(strings.Contains(string(out), iface)) {
			_, err = gcmd.RunScript(fmt.Sprintf("sudo route add -net %s -netmask %s -interface %s", destIPNet.IP.String(), WrapIPMask(destIPNet.Mask).DecString(), iface))
		}
		return err
	case "linux": // FIXME 待验证
		if gcmd.CommandExists("ip") {
			out, err := gcmd.RunScript(fmt.Sprintf("ip route get %s", destIPNet.String()))
			if err != nil || !strings.Contains(string(out), iface) {
				_, err = gcmd.RunScript(fmt.Sprintf("ip route add %s dev %s", WrapIPNet(destIPNet).CidrStyleString(), iface))
			}
		} else if gcmd.CommandExists("route") {
			out, err := gcmd.RunScript(fmt.Sprintf("route get %s", destIPNet.String()))
			if err != nil || !strings.Contains(string(out), iface) {
				_, err = gcmd.RunScript(fmt.Sprintf("route add -net %s netmask %s dev %s", destIPNet.IP.String(), WrapIPMask(destIPNet.Mask).DecString(), iface))
			}
		} else {
			return gerrors.New("neither `ip` nor `route` command found")
		}
		return nil
	case "freebsd": // FIXME 待验证
		_, err := gcmd.RunScript("route add -net " + destIPNet.String() + " -interface " + iface)
		return err
	case "windows": // FIXME 待验证
		_, err := gcmd.RunScript("route ADD " + destIPNet.String() + " " + address)
		if err != nil {
			return err
		}
		time.Sleep(time.Second >> 2)
		_, err = gcmd.RunScript("route CHANGE " + destIPNet.IP.String() + " MASK " + destIPNet.Mask.String() + " " + address)
		return err
	default:
		return errUnsupportedOS
	}
}

// DeleteNetRoute removes `net` type route rule from local operating system
// FIXME 是否需要sudo？
func DeleteNetRoute(iface string, addr net.IPNet, address string) error {
	switch runtime.GOOS {
	case "darwin": // verified
		_, err := gcmd.RunScript(fmt.Sprintf("route delete -net %s -interface %s", addr.IP.String(), iface))
		return err
	case "linux": // FIXME 待验证
		if gcmd.CommandExists("ip") {
			if _, err := gcmd.RunScript(fmt.Sprintf("ip route del %s dev %s", addr.String(), iface)); err != nil {
				return err
			}
		} else if gcmd.CommandExists("route") {
			if _, err := gcmd.RunScript(fmt.Sprintf("route del %s dev %s", addr.String(), iface)); err != nil {
				return err
			}
		} else {
			return gerrors.New("neither `ip` nor `route` command found")
		}
		return nil
	case "freebsd": // FIXME 待验证
		_, err := gcmd.RunScript("route delete -net " + addr.String() + " -interface " + iface)
		return err
	case "windows": // FIXME 待验证
		_, err := gcmd.RunScript("route DELETE " + addr.IP.String() + " MASK " + addr.Mask.String() + " " + address)
		return err
	default:
		return errUnsupportedOS
	}
}

/*
// AddCidr - sets the CIDR route, used on join and restarts
func AddCidr(iface, address string, addr *net.IPNet) error {
	switch runtime.GOOS {
	case "darwin":
		var cmd string
		if WrapIP(addr.IP).IsV4() {
			cmd = "route -q -n add -net " + addr.String() + " " + address
		} else if WrapIP(addr.IP).IsV6() {
			cmd = "route -A inet6 -q -n add -net " + addr.String() + " " + address
		} else {
			return gerrors.New("could not parse address: " + addr.String())
		}
		_, err := gcmd.RunScript(cmd)
		return err
	case "linux":
		var cmd string
		if WrapIP(addr.IP).IsV4() {
			cmd = "ip -4 route add " + addr.String() + " dev " + iface
		} else if WrapIP(addr.IP).IsV6() {
			cmd = "ip -6 route add " + addr.String() + " dev " + iface
		} else {
			return gerrors.New("could not parse address: " + addr.String())
		}
		_, err := gcmd.RunScript(cmd)
		return err
	case "freebsd":
		var cmd string
		if WrapIP(addr.IP).IsV4() {
			cmd = "route add -net " + addr.String() + " -interface " + iface
		} else if WrapIP(addr.IP).IsV6() {
			cmd = "route add -net -inet6 " + addr.String() + " -interface " + iface
		} else {
			return gerrors.New("could not parse address: " + addr.String())
		}
		_, err := gcmd.RunScript(cmd)
		return err
	case "windows":
		_, err := gcmd.RunScript("route ADD " + addr.String() + " " + address)
		if err != nil {
			return err
		}
		time.Sleep(time.Second >> 2)
		_, err = gcmd.RunScript("route CHANGE " + addr.IP.String() + " MASK " + addr.Mask.String() + " " + address)
		return err
	default:
		return errUnsupportedOS
	}
}

// DeleteCidr - removes a static cidr route
func DeleteCidr(iface string, addr *net.IPNet, address string) error {
	switch runtime.GOOS {
	case "darwin":
		_, err := gcmd.RunScript("route -q -n delete " + addr.String() + " -interface " + iface)
		return err
	case "linux":
		_, err := gcmd.RunScript("ip route delete " + addr.String() + " dev " + iface)
		return err
	case "freebsd":
		_, err := gcmd.RunScript("route delete -net " + addr.String() + " -interface " + iface)
		return err
	case "windows":
		_, err := gcmd.RunScript("route DELETE " + addr.String())
		return err
	default:
		return errUnsupportedOS
	}
}
*/
