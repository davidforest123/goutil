package gcompress

// FIXME: this package doesn't work well for now, generated .zip file can't be read by other zip apps.

import (
	"github.com/davidforest123/goutil/basic/gerrors"
	yekaZip "github.com/davidforest123/goutil/compress/gcompress/yekazip"
	"github.com/davidforest123/goutil/sys/gfs"
	"github.com/davidforest123/goutil/sys/gio"
	"io"
	"os"
	"sync"
)

type (
	CompFile struct {
		filename string // zip filename
		algo     Comp
		password string
		param    *CompParam
		mu       sync.RWMutex
		f        *os.File        // zip file
		fw       *yekaZip.Writer // zip file used to write, if user called fw.Close(), f will be closed too, so don't call fw.Close() before CompFile.Close()
		ingTasks map[string]bool // map[itemFilenameInZip]readTrueWriteFalse
	}

	compReadCloser struct {
		compFile *CompFile
		inRc     io.ReadCloser
	}

	compWriteFlushCloser struct {
		compFile *CompFile
		inW      io.Writer
	}
)

func (cf *CompFile) Walk() (dirs []string, files []string, err error) {
	rLockOK := cf.mu.TryRLock()
	if !rLockOK {
		return nil, nil, gerrors.New("Can't walk zip file [%s] for now, there's write task in progress at the moment.", cf.filename)
	}
	defer cf.mu.RUnlock()

	switch cf.algo {
	case CompZip:
		rd, err := yekaZip.OpenReaderWithFile(cf.f)
		if err != nil {
			return nil, nil, err
		}
		for _, v := range rd.File {
			if v.FileInfo().IsDir() {
				dirs = append(dirs, v.Name)
			} else {
				files = append(files, v.Name)
			}
		}
		err = rd.Close()
		if err != nil {
			return nil, nil, err
		}
		return dirs, files, nil
	default:
		return nil, nil, gerrors.New("Walk unsupported compress algorithm %s", cf.algo)
	}
}

func (cf *CompFile) newReadCloser(filenameInsideZip string) (io.ReadCloser, error) {
	rLockOK := cf.mu.TryRLock()
	if !rLockOK {
		return nil, gerrors.New("Can't read zip file [%s][%s] for now, there's write task in progress at the moment.", cf.filename, filenameInsideZip)
	}

	rd, err := yekaZip.OpenReaderWithFile(cf.f)
	if err != nil {
		cf.mu.RUnlock()
		return nil, err
	}
	err = error(nil)
	inRc := io.ReadCloser(nil)
	for _, v := range rd.File {
		if v.Name == filenameInsideZip {
			inRc, err = v.Open()
			if err != nil {
				cf.mu.RUnlock()
				return nil, err
			}
		}
	}
	if inRc == nil {
		cf.mu.RUnlock()
		return nil, gerrors.New("file[%s] not found in zip file[%s]", filenameInsideZip, cf.filename)
	}
	if err = rd.Close(); err != nil {
		cf.mu.RUnlock()
		return nil, err
	}

	return &compReadCloser{
		compFile: cf,
		inRc:     inRc,
	}, nil
}

func (rc *compReadCloser) Read(p []byte) (int, error) {
	return rc.inRc.Read(p)
}

func (rc *compReadCloser) Close() error {
	rc.compFile.mu.RUnlock()
	return rc.inRc.Close()
}

func (cf *CompFile) ReadFile(filenameInsideZip string) (io.ReadCloser, error) {
	switch cf.algo {
	case CompZip:
		return cf.newReadCloser(filenameInsideZip)
	default:
		return nil, gerrors.New("Walk unsupported compress algorithm %s", cf.algo)
	}
}

func (cf *CompFile) newWriteFlushCloser(filenameInsideZip string, encrypt Encrypt, level Level) (gio.WriteFlushCloser, error) {
	yekaMtd := yekaZip.LevelStore
	switch level {
	case LevelStore:
		yekaMtd = yekaZip.LevelStore
	case LevelDeflate:
		yekaMtd = yekaZip.LevelDeflate
	default:
		return nil, gerrors.New("Unsupported level %s", level)
	}

	yekaEnc := yekaZip.AES256Encryption
	switch encrypt {
	case EncryptAES128:
		yekaEnc = yekaZip.AES128Encryption
	case EncryptAES192:
		yekaEnc = yekaZip.AES192Encryption
	case EncryptAES256:
		yekaEnc = yekaZip.AES256Encryption
	default:
		return nil, gerrors.New("Unsupported encrypt %s", encrypt)
	}

	lockOK := cf.mu.TryLock()
	if !lockOK {
		return nil, gerrors.New("Can't write zip file [%s][%s] for now, there's read/write task in progress at the moment.", cf.filename, filenameInsideZip)
	}

	// check if filenameInsideZip exist
	/*_, err := cf.f.Seek(0, 0)
	if err != nil {
		cf.mu.Unlock()
		return nil, err
	}
	status, err := cf.f.Stat()
	if err != nil {
		cf.mu.Unlock()
		return nil, err
	}
	if status.Size() > 0 {
		rd, err := yekaZip.OpenReaderWithFile(cf.f)
		if err != nil {
			cf.mu.Unlock()
			return nil, err
		}
		for _, v := range rd.File {
			if v.Name == filenameInsideZip {
				cf.mu.Unlock()
				return nil, gerrors.New("file[%s] already exist in zip file [%s]", filenameInsideZip, cf.filename)
			}
		}
		err = rd.Close()
		if err != nil {
			cf.mu.Unlock()
			return nil, err
		}
	}*/

	// add new file into zip file
	inW, err := cf.fw.Encrypt(filenameInsideZip, cf.password, yekaEnc, yekaMtd)
	if err != nil {
		cf.mu.Unlock()
		return nil, err
	}
	return &compWriteFlushCloser{compFile: cf, inW: inW}, nil
}

func (wfc *compWriteFlushCloser) Write(p []byte) (int, error) {
	return wfc.inW.Write(p)
}

func (wfc *compWriteFlushCloser) Flush() error {
	return wfc.compFile.fw.Flush()
}

func (wfc *compWriteFlushCloser) Close() error {
	wfc.compFile.mu.Unlock()
	// NOTICE: don't call wfc.compFile.fw.Close(), otherwise, CompFile.f will be closed too.
	return nil
}

// AddFile adds new file into zip file.
func (cf *CompFile) AddFile(filenameInsideZip string, encrypt Encrypt, level Level) (gio.WriteFlushCloser, error) {
	return cf.newWriteFlushCloser(filenameInsideZip, encrypt, level)
}

// Close implements io.ReadWriteCloser.
func (cf *CompFile) Close() error {
	lockOK := cf.mu.TryLock()
	if !lockOK {
		return gerrors.New("Can't close zip file [%s] for now, there's read/write task in progress at the moment.", cf.filename)
	}
	defer cf.mu.Unlock()

	return cf.f.Close()
}

// NewCompFile open/create compressed file.
func NewCompFile(filename string, compAlgo Comp, password string, param *CompParam) (*CompFile, error) {
	rst := new(CompFile)
	rst.filename = filename
	rst.algo = compAlgo
	rst.password = password
	if gfs.DirExits(filename) {
		return nil, gerrors.New("%s is a directory but not a .zip file", filename)
	}
	err := error(nil)
	if gfs.FileExits(filename) {
		rst.f, err = os.Open(filename) // open zip file
	} else {
		rst.f, err = os.Create(filename) // create zip file
	}
	if err != nil {
		return nil, err
	}
	rst.fw = yekaZip.NewWriter(rst.f)
	if param != nil {
		*rst.param = *param
		if err := param.Verify(compAlgo); err != nil {
			return nil, err
		}
	}
	switch compAlgo {
	case CompNone:
		return nil, gerrors.New("can't create CompFile for 'none' algo")
	case CompZip:
	default:
		return nil, gerrors.New("NewCompFS unsupported compress algorithm %s", compAlgo)
	}

	return rst, nil
}
