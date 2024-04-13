package gfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gstring"
	"github.com/davidforest123/goutil/container/gternary"
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

func DetectPathOS(path string) (string, bool) {
	// unix absolute path
	if strings.HasPrefix(path, "/") {
		return "linux", true
	}

	// windows absolute path
	pathLower := strings.ToLower(path)
	if len(path) == 2 && (pathLower[0] >= 'a' && pathLower[0] <= 'z') && pathLower[1] == ':' {
		return "windows", true
	}
	if len(path) > 2 && (pathLower[0] >= 'a' && pathLower[0] <= 'z') && pathLower[1] == ':' && pathLower[2] == '\\' {
		return "windows", true
	}

	// unix relative path
	// NOTICE: "\\" could be windows relative path or linux relative path, so please don't use "\\" to detect os
	if strings.Contains(path, "/") {
		return "linux", true
	}

	// not sure
	return "", false
}

// IsAbs reports whether the path is absolute.
// This function differs from filepath.IsAbs(),
// this function can handle paths to different os than the runtime.GOOS,
// while filepath.IsAbs() can only handle paths to same os as the current runtime.GOOS.
func IsAbs(path string) bool {
	isAbsUnix := strings.HasPrefix(path, "/")
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
}

// first elem must be detectable
func Join(elem ...string) (string, error) {
	if len(elem) == 0 {
		return "", nil
	}

	firstPathOs, firstOk := DetectPathOS(elem[0])
	if !firstOk {
		return "", gerrors.New("can not detect os of first elem(%s)", elem[0])
	}
	for _, v := range elem[1:] {
		currPathOs, currOk := DetectPathOS(v)
		if currOk && currPathOs != firstPathOs {
			return "", gerrors.New("current elem(%s) os(%s) != first elem(%s) os(%s)", v, currPathOs, elem[0], firstPathOs)
		}
	}

	seperator := gternary.If(firstPathOs == "windows").String("\\", "/")

	splitList := []string{}
	for _, e := range elem {
		splitList = append(splitList, strings.Split(e, seperator)...)
	}
	splitList = gstring.RemoveEmpty(splitList)

	for i, v := range splitList {
		if v == "." {
			splitList[i] = ""
		}
		if v == ".." && (i-1) >= 0 {
			splitList[i] = ""
			splitList[i-1] = ""
		}
	}

	splitList = gstring.RemoveEmpty(splitList)
	if strings.HasPrefix(elem[0], seperator) {
		splitList = append([]string{""}, splitList...)
	}

	return strings.Join(splitList, seperator), nil
}

// srcDir must be detectable
// Notice: PathJoin("/Users/tony", "test.js") = "/Users/test.js"
// Notice: DirJoinFile("/Users/tony", "test.js") = "/Users/tony/test.js"
// Combine absolute dir path and relative file path to get a new absolute file path
func DirJoinFile(srcDir, targetPath string) (string, error) {
	if IsAbs(targetPath) {
		return "", gerrors.New("targetPath(%s) is an absolute path", targetPath)
	}
	return Join(srcDir, targetPath)
}

// srcDir must be detectable
// Notice: PathJoin("/Users/tony", "test.js") = "/Users/test.js"
// Notice: DirJoinDir("/Users/tony", "test.js") = "/Users/tony/test.js/"
// Combine absolute dir path and relative dir path to get a new absolute dir path
func DirJoinDir(srcDir, targetPath string) (string, error) {
	if IsAbs(targetPath) {
		return "", gerrors.New("targetPath(%s) is an absolute path", targetPath)
	}
	srcDir = AppendPathSeparatorIfNecessary(srcDir, AsDirOS)
	result, err := Join(srcDir, targetPath)
	if err != nil {
		return "", err
	}
	return AppendPathSeparatorIfNecessary(result, AsDirOS), nil
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
