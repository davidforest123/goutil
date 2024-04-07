package gprobe

import (
	"goutil/basic/gerrors"
	"goutil/container/gconv"
	"goutil/net/gnet"
	"net"
	"time"
)

func TCPing(host string, port int, timeout time.Duration) (opened bool, duration time.Duration, err error) {
	if !gnet.IsValidPort(port) {
		return false, 0, gerrors.Errorf("Invalid port " + gconv.NumToString(port))
	}

	ip := ""
	if gnet.IsIPString(host) {
		ip = host
	} else {
		ipArr, err := gnet.LookupIP(host)
		if err != nil {
			return false, 0, err
		}
		if len(ipArr) == 0 {
			return false, 0, gerrors.New("ip addresses length is zero")
		}
		ip = ipArr[0].String()
	}

	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", ip+":"+gconv.NumToString(port), timeout)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return false, time.Now().Sub(startTime), nil // Maybe Closed
	}
	return true, time.Now().Sub(startTime), nil // Opened
}
