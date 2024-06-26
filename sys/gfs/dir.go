package gfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/container/gstring"
	"github.com/davidforest123/goutil/container/gternary"
)

type (
	ProcessPathAs string
)

const (
	AsDirOS     = ProcessPathAs("AsDirOS")
	AsRuntimeOS = ProcessPathAs("AsRuntimeOS")
)

// base name: final element of the path, not the entire path.
// full name: the entire path.

func DirExits(dir string) bool {
	pi, err := GetPathInfo(dir)
	if err != nil {
		return false
	}
	return pi.IsFolder && pi.Exist
}

func MakeDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func RemoveDir(dir string) error {
	return os.RemoveAll(dir)
}

// TODO: test required
func RemoveAllEmptyDirs(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			files, errRD := os.ReadDir(path)
			if errRD != nil {
				return errRD
			}
			if len(files) == 0 {
				errRM := os.Remove(path)
				if errRM != nil {
					return errRM
				}
			}
		}
		return nil
	})
	return err
}

func DirEmpty(dir string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(files) == 0, nil
}

// get total size of directory
// os.Stat(dir string) read index size of dir, not total size
func DirSize(dir string) (int64, map[string]int64, error) {
	var size int64
	everyFileSize := map[string]int64{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
			everyFileSize[path] = info.Size()
		}
		return err
	})
	return size, everyFileSize, err
}

// Remove all content under dir, but keep dir folder
func CleanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func AppendPathSeparatorIfNecessary(dir string, as ProcessPathAs) string {
	switch as {
	case AsDirOS:
		pathOS, ok := DetectPathOS(dir)
		if !ok {
			panic(gerrors.New("AppendPathSeparatorIfNecessary=>DetectPathOS(%s) failed", dir))
		}
		seperator := gternary.If(pathOS == "windows").String("\\", "/")
		if strings.HasSuffix(dir, seperator) {
			return dir
		} else {
			return dir + seperator
		}
	case AsRuntimeOS:
		if dir != "" && dir != "." && dir != ".." && dir[len(dir)-1] != os.PathSeparator {
			return dir + string(os.PathSeparator)
		}
		return dir
	default:
		panic(gerrors.New("AppendPathSeparatorIfNecessary accept error param(%s)", as))
	}
}

// WalkDirTopLevel lists the top level file/dir at dir, it will NOT scan any subdirectory in `dir`.
func WalkDirTopLevel(dir string) (dirs []string, files []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range entries {
		fullPath, err := DirJoinFile(dir, v.Name())
		if err != nil {
			return nil, nil, err
		}
		if v.IsDir() {
			dirs = append(dirs, fullPath)
		} else {
			files = append(files, fullPath)
		}
	}
	return dirs, files, nil
}

// WalkDir walks the file tree rooted at dir, it will automatically scan each subdirectory in `dir`.
func WalkDir(dir string) (dirs []string, files []string, err error) {
	if err := filepath.WalkDir(dir,
		func(path string, de os.DirEntry, err error) error {
			if de == nil {
				return err
			}
			if de.IsDir() {
				dirs = append(dirs, path)
			} else {
				files = append(files, path)
			}
			return nil
		}); err != nil {
		return nil, nil, err
	}
	return dirs, files, nil
}

func ListDirContains(dir, contains string) (dirs []string, files []string, err error) {
	if err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if strings.Contains(path, contains) {
				if f.IsDir() {
					dirs = append(dirs, path)
				} else {
					files = append(files, path)
				}
			}
			return nil
		}); err != nil {
		return nil, nil, err
	}
	return dirs, files, nil
}

func ListDirEndWith(dir, tail string) (dirs []string, files []string, err error) {
	if err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if gstring.EndsWith(path, tail) {
				if f.IsDir() {
					dirs = append(dirs, path)
				} else {
					files = append(files, path)
				}
			}
			return nil
		}); err != nil {
		return nil, nil, err
	}
	return dirs, files, nil
}

func DirSlash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

type totalInfo struct {
	mu              sync.RWMutex
	copiedSizeOfDir int64
}

func copyDirEx(src string, dst string, dirSizeCallback DirCopiedSizeCallback, totalInfo *totalInfo) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if totalInfo == nil {
		return gerrors.Errorf("nil totalInfo")
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDirEx(srcPath, dstPath, dirSizeCallback, totalInfo)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			fileSizeCallback := func(currFileCopiedSize int64) {
				dirSizeCallback(srcPath, dstPath, currFileCopiedSize, totalInfo.copiedSizeOfDir+currFileCopiedSize)
			}
			written, err := CopyFileEx(srcPath, dstPath, fileSizeCallback)
			if totalInfo != nil {
				totalInfo.mu.Lock()
				totalInfo.copiedSizeOfDir += written
				totalInfo.mu.Unlock()
			}
			dirSizeCallback(srcPath, dstPath, written, totalInfo.copiedSizeOfDir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type DirCopiedSizeCallback func(currSrcFile, currDstFile string, copiedSizeOfCurrFile, copiedSizeOfDir int64)

func CopyDirEx(src string, dst string, dirSizeCallback DirCopiedSizeCallback) (err error) {
	totalInfo := &totalInfo{}
	return copyDirEx(src, dst, dirSizeCallback, totalInfo)
}

/*
func CopyDirEx(src string, dst string, dirSizeCallback DirCopiedSizeCallback, alreadyDirCopiedSize *int64) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	dirCopiedSize := int64(0)
	if alreadyDirCopiedSize != nil {
		dirCopiedSize = *alreadyDirCopiedSize
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDirEx(srcPath, dstPath, dirSizeCallback, &dirCopiedSize)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			fileSizeCallback := func(currFileCopiedSize int64) {
				dirSizeCallback(srcPath, dstPath, currFileCopiedSize, dirCopiedSize+currFileCopiedSize)
			}
			written, err := CopyFileEx(srcPath, dstPath, fileSizeCallback)
			dirCopiedSize += written
			dirSizeCallback(srcPath, dstPath, written, dirCopiedSize)
			if err != nil {
				return err
			}
		}
	}

	return nil
}*/
