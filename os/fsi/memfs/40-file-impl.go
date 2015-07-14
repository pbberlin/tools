package memfs

import (
	"bytes"
	"io"
	"os"
	"sync/atomic"

	"github.com/pbberlin/tools/os/fsi"
)

func (f *InMemoryFile) Open() error {
	atomic.StoreInt64(&f.at, 0)
	f.Lock()
	f.closed = false
	f.Unlock()
	return nil
}

func (f *InMemoryFile) Close() error {
	atomic.StoreInt64(&f.at, 0)
	f.Lock()
	f.closed = true
	f.Unlock()
	return nil
}

// To remain consistent with osfs, we can only return base name.
func (f *InMemoryFile) Name() string {
	_, bname := f.fs.pathInternalize(f.name)
	return bname
}

func (f *InMemoryFile) Stat() (os.FileInfo, error) {
	return &InMemoryFileInfo{f}, nil
}

func (f *InMemoryFile) Readdir(count int) (res []os.FileInfo, err error) {

	limit := len(f.memDir)

	if len(f.memDir) == 0 {
		return
	}

	if count > 0 {
		limit = count
	}

	if len(f.memDir) < limit {
		err = io.EOF
	}

	res = make([]os.FileInfo, len(f.memDir))

	i := 0
	for _, file := range f.memDir {
		res[i], _ = file.Stat()
		i++
	}
	return res, nil
}

func (f *InMemoryFile) Readdirnames(n int) (names []string, err error) {
	fi, err := f.Readdir(n)
	names = make([]string, len(fi))
	for i, f := range fi {
		names[i] = f.Name()
	}
	return names, err
}

func (f *InMemoryFile) Read(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	if f.closed == true {
		return 0, fsi.ErrFileClosed
	}
	if len(b) > 0 && int(f.at) == len(f.data) {
		return 0, io.EOF
	}
	if len(f.data)-int(f.at) >= len(b) {
		n = len(b)
	} else {
		n = len(f.data) - int(f.at)
	}
	copy(b, f.data[f.at:f.at+int64(n)])
	atomic.AddInt64(&f.at, int64(n))
	return
}

func (f *InMemoryFile) ReadAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Read(b)
}

func (f *InMemoryFile) Truncate(size int64) error {
	if f.closed == true {
		return fsi.ErrFileClosed
	}
	if size < 0 {
		return fsi.ErrOutOfRange
	}
	if size > int64(len(f.data)) {
		diff := size - int64(len(f.data))
		f.data = append(f.data, bytes.Repeat([]byte{00}, int(diff))...)
	} else {
		f.data = f.data[0:size]
	}
	return nil
}

func (f *InMemoryFile) Seek(offset int64, whence int) (int64, error) {
	if f.closed == true {
		return 0, fsi.ErrFileClosed
	}
	switch whence {
	case 0:
		atomic.StoreInt64(&f.at, offset)
	case 1:
		atomic.AddInt64(&f.at, int64(offset))
	case 2:
		atomic.StoreInt64(&f.at, int64(len(f.data))+offset)
	}
	return f.at, nil
}

func (f *InMemoryFile) Write(b []byte) (n int, err error) {
	n = len(b)
	cur := atomic.LoadInt64(&f.at)
	f.Lock()
	defer f.Unlock()
	diff := cur - int64(len(f.data))
	var tail []byte
	if n+int(cur) < len(f.data) {
		tail = f.data[n+int(cur):]
	}
	if diff > 0 {
		f.data = append(bytes.Repeat([]byte{00}, int(diff)), b...)
		f.data = append(f.data, tail...)
	} else {
		f.data = append(f.data[:cur], b...)
		f.data = append(f.data, tail...)
	}

	atomic.StoreInt64(&f.at, int64(len(f.data)))
	return
}

func (f *InMemoryFile) WriteAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Write(b)
}

func (f *InMemoryFile) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}

func (f *InMemoryFile) Info() *InMemoryFileInfo {
	return &InMemoryFileInfo{file: f}
}
