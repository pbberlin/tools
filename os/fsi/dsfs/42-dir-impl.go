package dsfs

import (
	"io"
	"os"
)

// These methods on AeDir are mostly useless,
// except for Readdir(), Readdirnames() and Stat().
// But we have to implement them in order to
// convert AeDir into fsi.File.
//

func (f *AeDir) Close() error {
	return nil
}

func (f *AeDir) Read(b []byte) (n int, err error) {
	return
}
func (f *AeDir) ReadAt(b []byte, off int64) (n int, err error) {
	return f.Read(b)
}

// Adapt (f *AeFile) Readdir also
func (f *AeDir) Readdir(n int) (fis []os.FileInfo, err error) {

	wantAll := n <= 0
	fis, err = f.fSys.ReadDir(f.Dir + f.BName)
	if wantAll {
		return fis, nil
	}

	if f.memDirFetchPos == 0 {
		f.memDirFetchPos = len(fis)
		return fis, nil
	} else {
		f.memDirFetchPos = 0
		return []os.FileInfo{}, io.EOF
	}
}

func (f *AeDir) Readdirnames(n int) (names []string, err error) {
	fis, err := f.Readdir(n)
	names = make([]string, 0, len(fis))
	for _, lp := range fis {
		names = append(names, lp.Name())
	}
	return names, err
}

func (f *AeDir) Seek(offset int64, whence int) (int64, error) {
	return int64(0), nil
}

func (f *AeDir) Stat() (os.FileInfo, error) {
	return os.FileInfo(*f), nil
}

func (f *AeDir) Truncate(size int64) error {
	return nil
}

func (f *AeDir) Write(b []byte) (n int, err error) {
	return
}

func (f *AeDir) WriteAt(b []byte, off int64) (n int, err error) {
	return f.Write(b)
}

func (f *AeDir) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}
