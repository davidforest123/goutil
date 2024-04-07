package gnet

import (
	"errors"
	"goutil/sys/gcmd"
	"runtime"
	"strings"
)

// SetIPForwarding - Sets IP forwarding if it's mac or linux
func SetIPForwarding() error {
	os := runtime.GOOS
	var err error
	switch os {
	case "linux":
		err = setIPForwardingUnix()
	case "freebsd":
		err = setIPForwardingFreeBSD()
	case "darwin":
		err = setIPForwardingMac()
	default:
		err = errors.New("this OS is not currently supported")
	}
	return err
}

// setIPForwardingUnix - sets the ipforwarding for linux
func setIPForwardingUnix() error {

	// ipv4
	out, err := gcmd.RunScript("sysctl net.ipv4.ip_forward")
	if err != nil {
		return err
	} else {
		s := strings.Fields(string(out))
		if s[2] != "1" {
			_, err = gcmd.RunScript("sysctl -w net.ipv4.ip_forward=1")
			if err != nil {
				return err
			}
		}
	}

	// ipv6
	out, err = gcmd.RunScript("sysctl net.ipv6.conf.all.forwarding")
	if err != nil {
		return err
	} else {
		s := strings.Fields(string(out))
		if s[2] != "1" {
			_, err = gcmd.RunScript("sysctl -w  net.ipv6.conf.all.forwarding=1")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// setIPForwardingFreeBSD - sets the ipforwarding for linux
func setIPForwardingFreeBSD() error {
	out, err := gcmd.RunScript("sysctl net.inet.ip.forwarding")
	if err != nil {
		return err
	} else {
		s := strings.Fields(string(out))
		if s[1] != "1" {
			_, err = gcmd.RunScript("sysctl -w net.inet.ip.forwarding=1")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetIPForwardingMac - sets ip forwarding for mac
func setIPForwardingMac() error {
	_, err := gcmd.RunScript("sysctl -w net.inet.ip.forwarding=1")
	return err
}
