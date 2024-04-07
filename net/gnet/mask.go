package gnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"goutil/basic/gerrors"
	"goutil/container/gconv"
	"net"
	"strings"
)

type (
	// IPMask only exists in IPv4, there is no netmask in IPv6.
	IPMask net.IPMask
)

func WrapIPMask(ipMask net.IPMask) IPMask {
	return IPMask(ipMask)
}

// ParseIPMask casts string like `255.255.255.0` to IPv4 net mask, IPv6 doesn't need net mask.
// Subnet Mask represents the Network by defining the Leading Bits
// as 1’s while the Hosts with Trailing Bits as 0’s.
// 子网掩码必须以1为前导，以0为尾巴，不可以交叉存在
// Invalid mask sample: 255.255.253.0(11111111.11111111.11111101.00000000).
func ParseIPMask(s string) (IPMask, error) {
	defErr := gerrors.New("invalid IPMask %s", s)
	m := net.IPMask(net.ParseIP(s).To4())
	if m == nil {
		return nil, defErr
	}

	// verify if IPMask is valid
	decStr := gconv.BytesToBinString(m, "")
	if decStr == "0.0.0.0" || decStr == "255.255.255.255" {
		return IPMask(m), nil
	}
	binStr := gconv.BytesToBinString(m, "")
	ss := strings.Split(binStr, "10")
	if len(ss) != 2 {
		return nil, defErr
	}
	for _, v := range ss[0] { // Leading Bits must be 1’s
		if v != '1' {
			return nil, defErr
		}
	}
	for _, v := range ss[1] { // Trailing Bits must be 0’s
		if v != '0' {
			return nil, defErr
		}
	}

	return IPMask(m), nil
}

func MustParseIPMask(s string) IPMask {
	ipMask, err := ParseIPMask(s)
	if err != nil {
		panic(err)
	}
	return ipMask
}

func (im IPMask) Std() net.IPMask {
	return net.IPMask(im)
}

func (im IPMask) DecString() string {
	return fmt.Sprintf("%d.%d.%d.%d", im[0], im[1], im[2], im[3])
}

// ToSimpleNotation Converts IP mask to 16 bit unsigned integer.
func (im IPMask) ToSimpleNotation(mask net.IPMask) (uint16, error) {
	var i uint16
	buf := bytes.NewReader(mask)
	err := binary.Read(buf, binary.BigEndian, &i)
	return i, err
}
