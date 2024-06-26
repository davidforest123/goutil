package route_with_api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type NetRoute struct {
	Metric      uint32 `json:"metric"`
	Destination string `json:"dest"`
	Gateway     string `json:"gateway"`
	Flags       string `json:"flags"`
	NetIf       string `json:"iface"`
}

func (nr NetRoute) ToPortableJSON() string {
	data, err := json.Marshal(nr)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (nr NetRoute) ToTableString() string {
	// NR doesn't have any header
	return fmt.Sprintf("%s\tvia %s\tdev %s\tflags %s\tmetric %d\n", nr.Destination, nr.Gateway, nr.NetIf, nr.Flags, nr.Metric)
}

type RouteFlag struct {
	U        bool `json:"up"`                       // 路由是活动的
	H        bool `json:"dest_is_single_host"`      // 目标是一个主机
	G        bool `json:"gateway"`                  // 路由指向网关
	S        bool `json:"static"`                   // bsd
	Cloned   bool `json:"clone_based_on_route"`     // bsd
	W        bool `json:"clone_auto_local"`         // bsd
	L        bool `json:"link_to_hw"`               // bsd
	Reinsta  bool `json:"reinstate_route"`          // linux
	D        bool `json:"dynamic_installed"`        // linux 由路由的后台程序动态地安装
	M        bool `json:"modified_from_routing_sw"` // linux 由路由的后台程序修改
	A        bool `json:"installed_by_addrconf"`    // linux
	Cached   bool `json:"cached"`                   // linux
	Rejected bool `json:"rejected"`                 // linux `!` 拒绝路由
}

func (rf RouteFlag) ToTableString() string {
	hasPrevious := false
	var sb strings.Builder
	v := reflect.ValueOf(rf)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Bool() {
			if hasPrevious {
				sb.WriteString(",")
			}
			sb.WriteString(v.Type().Field(i).Name)
			hasPrevious = true
		}
	}
	return sb.String()
}

func (rf RouteFlag) ToPortableJSON() string {
	data, err := json.Marshal(rf)
	if err != nil {
		panic(err)
	}
	return string(data)
}
