//go:build !linux && !darwin && !windows

package gtuntap

import "github.com/songgao/water"

func setConfigName(dst *water.Config, name string) {
}
