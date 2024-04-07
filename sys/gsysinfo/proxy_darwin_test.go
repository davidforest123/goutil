package gsysinfo

import (
	"fmt"
	"testing"
)

func TestGetCurrentNetworkInterface(t *testing.T) {
	fmt.Println(GetCurrentNetworkInterface())
}

func TestSetGlobalSocks5ProxyOn(t *testing.T) {
	fmt.Println(GetGlobalSocks5Proxy())
	fmt.Println(SetGlobalSocks5ProxyOn("127.0.0.1:8000"))
	fmt.Println(GetGlobalSocks5Proxy())
	fmt.Println(SetGlobalSocks5ProxyOn("127.0.0.1:6000"))
	fmt.Println(GetGlobalSocks5Proxy())
}
