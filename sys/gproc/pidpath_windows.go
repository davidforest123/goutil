//go:build windows

package gproc

import "github.com/davidforest123/goutil/basic/gerrors"

// TODO
func GetExePathFromPid(pid int) (path string, err error) {
	return "", gerrors.New("windows unsupported for now")
}
