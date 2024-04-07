package main

import (
	"fmt"
	"github.com/davidforest123/goutil/basic/glog"
	"github.com/davidforest123/goutil/net/garp"
)

func main() {
	s := garp.NewScanner()
	devices, err := s.Scan()
	if err != nil {
		glog.Erro(err)
		return
	}
	for _, v := range devices {
		fmt.Println(v.Name, v.IP.String(), v.MAC)
	}
}
