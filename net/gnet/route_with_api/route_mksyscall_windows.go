//go:build windows

package route_with_api

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output gen_iphlpapi_windows.go prototype_iphlpapi_windows.go
