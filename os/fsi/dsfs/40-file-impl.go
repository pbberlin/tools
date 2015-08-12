package dsfs

import (
	"io"
	"os"
	"sync/atomic"

	"github.com/pbberlin/tools/os/fsi"
)

func (f *DsFile) Close() error {

	atomic.StoreInt64(&f.at, 0)
	f.Lock()
	f.closed = true

	err := f.Sync()
	if err != nil {
		return err
	}

	f.Unlock()

	return nil
}

// See fsi.File interface.
// Adapt (f *AeDir) Readdir also
func (f *DsFile) Readdir(n int) (fis []os.FileInfo, err error) {

	fis, err = f.fSys.ReadDir(f.Dir)

	wantAll := n <= 0

	if wantAll {
		return fis, nil
	}

	// Actually we would need memDirFetchPos
	// holding the latest retrieved file in
	// a forwardly-linked-list mimic.
	// Compare https://golang.org/src/os/file_windows.go
	// Instead: We either or return *all* available files
	// or empty slice plus io.EOF
	if f.memDirFetchPos == 0 {
		f.memDirFetchPos = len(fis)
		return fis, nil
	} else {
		f.memDirFetchPos = 0
		return []os.FileInfo{}, io.EOF
	}

}

// See fsi.File interface.
func (f *DsFile) Readdirnames(n int) (names []string, err error) {
	fis, err := f.Readdir(n)
	names = make([]string, 0, len(fis))
	for _, lp := range fis {
		names = append(names, lp.Name())
	}
	return names, err
}

func (f *DsFile) Read(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	if f.closed == true {
		return 0, fsi.ErrFileClosed
	}
	if len(b) > 0 && int(f.at) == len(f.Data) {
		return 0, io.EOF
	}
	if len(f.Data)-int(f.at) >= len(b) {
		n = len(b)
	} else {
		n = len(f.Data) - int(f.at)
	}
	copy(b, f.Data[f.at:f.at+int64(n)])
	atomic.AddInt64(&f.at, int64(n))
	return
}

func (f *DsFile) ReadAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Read(b)
}

func (f *DsFile) Seek(offset int64, whence int) (int64, error) {
	if f.closed == true {
		return 0, fsi.ErrFileClosed
	}
	switch whence {
	case 0:
		atomic.StoreInt64(&f.at, offset)
	case 1:
		atomic.AddInt64(&f.at, int64(offset))
	case 2:
		atomic.StoreInt64(&f.at, int64(len(f.Data))+offset)
	}
	return f.at, nil
}

func (f *DsFile) Stat() (os.FileInfo, error) {
	return os.FileInfo(*f), nil
}

func (f *DsFile) Sync() error {

	err := f.fSys.saveFileByPath(f, f.Dir+f.BName)
	if err != nil {
		return err
	}
	return nil
}

func (f *DsFile) Truncate(size int64) error {
	if f.closed == true {
		return fsi.ErrFileClosed
	}
	if size < 0 {
		return fsi.ErrOutOfRange
	}
	if size > int64(len(f.Data)) {
		diff := size - int64(len(f.Data))
		// f.Content = append(f.Content, bytes.Repeat([]byte{00}, int(diff))...)
		sb := make([]byte, int(diff))
		f.Data = append(f.Data, sb...)
	} else {
		f.Data = f.Data[0:size]
	}
	return nil
}

func (f *DsFile) Write(b []byte) (n int, err error) {
	n = len(b)
	cur := atomic.LoadInt64(&f.at)
	f.Lock()
	defer f.Unlock()
	diff := cur - int64(len(f.Data))
	var tail []byte
	if n+int(cur) < len(f.Data) {
		tail = f.Data[n+int(cur):]
	}
	if diff > 0 {
		sb := make([]byte, int(diff))
		f.Data = append(sb, b...)
		// f.Content = append(bytes.Repeat([]byte{00}, int(diff)), b...)
		f.Data = append(f.Data, tail...)
	} else {
		f.Data = append(f.Data[:cur], b...)
		f.Data = append(f.Data, tail...)
	}

	atomic.StoreInt64(&f.at, int64(len(f.Data)))
	return
}

func (f *DsFile) WriteAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Write(b)
}

func (f *DsFile) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}
