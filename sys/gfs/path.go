package gfs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/davidforest123/goutil/basic/gerrors"
)

type PathInfo struct {
	Exist        bool
	IsFolder     bool
	ModifiedTime time.Time
}

const (
	InvalidFilenameCharsWindows = "\"\\:/*?<>|“”"
)

func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, gerrors.Wrap(err, fmt.Sprintf("Stat(%s)", path))
	}
	if fi.IsDir() {
		return 0, gerrors.Errorf("path(%s) is directory", path)
	}
	return fi.Size(), nil
}

func GetPathInfo(path string) (*PathInfo, error) {
	var pi PathInfo
	fi, err := os.Stat(path)
	if err == nil {
		pi.Exist = true
		pi.IsFolder = fi.IsDir()
		pi.ModifiedTime = fi.ModTime()
		return &pi, nil
	} else if err != nil && os.IsNotExist(err) {
		pi.Exist = false
		return &pi, nil
	} else {
		return &pi, err
	}
}

func FileExits(filename string) bool {
	pi, err := GetPathInfo(filename)
	if err != nil {
		return false
	}
	return !pi.IsFolder && pi.Exist
}

/*
// IsAbs reports whether the path is absolute.
func IsAbs(path string) bool {
	isAbsUnix := len(path) > 0 && path[0] == '/'
	pathLower := strings.ToLower(path)
	isAbsWindows := false
	if len(path) < 2 {
		isAbsWindows = false
	} else if len(path) == 2 {
		isAbsWindows = (pathLower[0] >= 'a' && pathLower[0] <= 'z') && pathLower[1] == ':'
	} else {
		isAbsWindows = (pathLower[0] >= 'a' && pathLower[0] <= 'z') && pathLower[1] == ':' && pathLower[2] == '\\'
	}
	return isAbsUnix || isAbsWindows
}*/

// Notice: PathJoin("/Users/tony", "test.js") = "/Users/test.js"
// Notice: DirJoinFile("/Users/tony", "test.js") = "/Users/tony/test.js"
// Combine absolute dir path and relative file path to get a new absolute file path
func DirJoinFile(srcDir, targetPath string) string {
	if path.IsAbs(targetPath) {
		return targetPath
	}
	srcDir = AppendPathSeparatorIfNecessary(srcDir)
	return path.Join(srcDir, targetPath)
}

// Notice: PathJoin("/Users/tony", "test.js") = "/Users/test.js"
// Notice: DirJoinDir("/Users/tony", "test.js") = "/Users/tony/test.js/"
// Combine absolute dir path and relative dir path to get a new absolute dir path
func DirJoinDir(srcDir, targetPath string) string {
	if path.IsAbs(targetPath) {
		return targetPath
	}
	srcDir = AppendPathSeparatorIfNecessary(srcDir)
	return AppendPathSeparatorIfNecessary(path.Join(srcDir, targetPath))
}

// PathBase returns last element of `path`.
// "/root/home/abc.txt" -> "abc.txt"
// note: this function doesn't work if file name contains '/', like "mydir/a/b.txt" and real file name is "a/b.txt"
func PathBase(path string) string {
	return filepath.Base(path)
}

// Replace illegal chars for short filename / dir name, not multi-level directory
func RefactShortPathName(path string) string {
	var illegalChars = "/\\:*\"<>|"
	for _, c := range illegalChars {
		path = strings.Replace(path, string(c), "-", -1)
	}
	return path
}
