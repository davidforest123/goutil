package gsysinfo

import (
	"goutil/basic/gerrors"
	//"github.com/klauspost/cpuid" // x86/x64 is supported only for now
	"github.com/shirou/gopsutil/v3/cpu"
	"runtime"
	"time"
)

// Get unique serial number of CPU
func GetSerialNumber() (string, error) {
	return "Unsupported for now", nil
}

// 获取所有CPU的使用百分比，以数组返回
func GetAllUsedPercent(duration time.Duration) ([]float64, error) {
	return cpu.Percent(duration, true)
}

// 获取所有CPU的使用百分比，组合成总百分比后返回
func GetCombinedUsedPercent(duration time.Duration) (float64, error) {
	p, err := cpu.Percent(duration, false)
	if err != nil {
		return 0, err
	}
	return p[0], err
}

func GetCpuCount() int {
	return runtime.NumCPU()
}

func isAllDigit(name string) bool {
	for _, c := range name {
		if c < '0' || c >= '9' {
			return false
		}
	}
	return true
}

// 控制CPU使用率，动态调整sleep时间
type DyncSleep struct {
	cpuUsage      float64 // 允许的CPU百分比
	lastSleepTime time.Duration
}

func NewDyncSleep(cpuUsage float64) (*DyncSleep, error) {
	if cpuUsage <= 0 || cpuUsage >= 100 {
		return nil, gerrors.Errorf("Invalid cpuUsage %f", cpuUsage)
	}
	return &DyncSleep{cpuUsage: cpuUsage, lastSleepTime: time.Millisecond}, nil
}

func (s *DyncSleep) Sleep() {
	used, err := GetCombinedUsedPercent(time.Second)
	if err != nil {
		time.Sleep(s.lastSleepTime)
	} else {
		if used > s.cpuUsage {
			s.lastSleepTime += time.Millisecond
		}
		if used < s.cpuUsage {
			s.lastSleepTime -= time.Millisecond

		}
		if s.lastSleepTime <= 0 {
			s.lastSleepTime = time.Millisecond
		}
		time.Sleep(s.lastSleepTime)
	}
}
