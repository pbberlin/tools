package aefs_sr

import (
	"io"
	"os"
	"sync/atomic"
)

func (f *AeFile) Close() error {
	atomic.StoreInt64(&f.at, 0)
	f.Lock()
	f.closed = true
	f.Unlock()
	return nil
}

func (f *AeFile) Readdir(count int) (res []os.FileInfo, err error) {
	return f.Fs.ReadDir(f.Dir)
}

func (f *AeFile) Readdirnames(n int) (names []string, err error) {
	fis, err := f.Readdir(n)
	names = make([]string, 0, len(fis))
	for _, lp := range fis {
		names = append(names, lp.Name())
	}
	return names, err
}

func (f *AeFile) Read(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	if f.closed == true {
		return 0, ErrFileClosed
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

func (f *AeFile) ReadAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Read(b)
}

func (f *AeFile) Seek(offset int64, whence int) (int64, error) {
	if f.closed == true {
		return 0, ErrFileClosed
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

func (f *AeFile) Stat() (os.FileInfo, error) {
	return os.FileInfo(*f), nil
}

func (f *AeFile) Sync() error {

	err := f.Fs.SaveFile(f, f.Dir)
	if err != nil {
		return err
	}
	return nil
}

func (f *AeFile) Truncate(size int64) error {
	if f.closed == true {
		return ErrFileClosed
	}
	if size < 0 {
		return ErrOutOfRange
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

func (f *AeFile) Write(b []byte) (n int, err error) {
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

func (f *AeFile) WriteAt(b []byte, off int64) (n int, err error) {
	atomic.StoreInt64(&f.at, off)
	return f.Write(b)
}

func (f *AeFile) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}
