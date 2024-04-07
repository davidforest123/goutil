package gtime

import (
	"time"
)

var (
	TimeZoneAsiaShanghai *time.Location = nil
)

func init() {
	// `TimeZoneAsiaShanghai, err = time.LoadLocation("Asia/Shanghai")` is not a good method
	// Under Windows, time.LoadLocation may return an error - "unknown time zone Asia/Shanghai"
	// time.LoadLocation relies on the IANA Time Zone Database (tzdata for short),
	// which is usually included in the Linux system, but not in the Windows system.
	TimeZoneAsiaShanghai = time.FixedZone("CST", 8*3600)
}

func GetLocalTimezone() (int, error) {
	return 0, nil
}

func SetLocalTimezone(timezone int) error {
	return nil
}

func ParseTimezoneCode(tz string) (offset int, err error) {
	return 0, nil //timezone.GetOffset()
}
