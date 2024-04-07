package gproxy

import (
	"fmt"
	"github.com/davidforest123/goutil/net/gprobe"
	"time"
)

func IsVisitable(url string) bool {
	fmt.Println(gprobe.TcpingOnline("www.youtube.com", 443, time.Millisecond*500))
	return false
}
