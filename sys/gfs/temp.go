package gfs

import (
	"github.com/google/uuid"
	"goutil/basic/gerrors"
	"runtime"
)

// Generate a new temp filename for cache
func NewTempFilename() (string, error) {
	if runtime.GOOS == "windows" {
		return "", gerrors.New("Unsupport windows for now")
	} else {
		return "/tmp/" + uuid.New().String() + ".temp", nil
	}
}
