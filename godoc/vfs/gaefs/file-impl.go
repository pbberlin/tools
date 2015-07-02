package gaefs

import (
	"io"
	"sync/atomic"
)

func (f *File) Open() error {
	atomic.StoreInt64(&f.at, 0)
	f.Lock()
	f.closed = false
	f.Unlock()
	return nil
}

func (f *File) Close() error {
	atomic.StoreInt64(&f.at, 0)
	f.Lock()
	f.closed = true
	f.Unlock()
	return nil
}

// func (f *File) Stat() (os.FileInfo, error) {
// 	return &InMemoryFileInfo{f}, nil
// }

// func (f *File) Readdir(count int) (res []os.FileInfo, err error) {
// }

// func (f *File) Readdirnames(n int) (names []string, err error) {
// }

func (f *File) Read(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	if f.closed == true {
		return 0, ErrFileClosed
	}
	if len(b) > 0 && int(f.at) == len(f.Content) {
		return 0, io.EOF
	}
	if len(f.Content)-int(f.at) >= len(b) {
		n = len(b)
	} else {
		n = len(f.Content) - int(f.at)
	}
	copy(b, f.Content[f.at:f.at+int64(n)])
	atomic.AddInt64(&f.at, int64(n))
	return
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Read(b)
}

func (f *File) Truncate(size int64) error {
	if f.closed == true {
		return ErrFileClosed
	}
	if size < 0 {
		return ErrOutOfRange
	}
	if size > int64(len(f.Content)) {
		diff := size - int64(len(f.Content))
		// f.Content = append(f.Content, bytes.Repeat([]byte{00}, int(diff))...)
		sb := make([]byte, int(diff))
		f.Content = append(f.Content, sb...)
	} else {
		f.Content = f.Content[0:size]
	}
	return nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.closed == true {
		return 0, ErrFileClosed
	}
	switch whence {
	case 0:
		atomic.StoreInt64(&f.at, offset)
	case 1:
		atomic.AddInt64(&f.at, int64(offset))
	case 2:
		atomic.StoreInt64(&f.at, int64(len(f.Content))+offset)
	}
	return f.at, nil
}

func (f *File) Write(b []byte) (n int, err error) {
	n = len(b)
	cur := atomic.LoadInt64(&f.at)
	f.Lock()
	defer f.Unlock()
	diff := cur - int64(len(f.Content))
	var tail []byte
	if n+int(cur) < len(f.Content) {
		tail = f.Content[n+int(cur):]
	}
	if diff > 0 {
		sb := make([]byte, int(diff))
		f.Content = append(sb, b...)
		// f.Content = append(bytes.Repeat([]byte{00}, int(diff)), b...)
		f.Content = append(f.Content, tail...)
	} else {
		f.Content = append(f.Content[:cur], b...)
		f.Content = append(f.Content, tail...)
	}

	atomic.StoreInt64(&f.at, int64(len(f.Content)))
	return
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Write(b)
}

func (f *File) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}
