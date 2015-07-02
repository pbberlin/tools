package gaefs

import (
	"bytes"
	"io/ioutil"
	"os"

	"appengine/datastore"

	"github.com/pbberlin/tools/logif"
	"golang.org/x/tools/godoc/vfs"

	pth "path"
)

func OS(mount string) FileSys {
	panic(`
		Sadly, google app engine file system requires a
	 	http.Request based context object.
	 	Use NewFs(string, appengine.Context) instead of OS.
	`)
}

func (fs *FileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f := File{}
	f.BName = pth.Base(filename)
	f.Dir = pth.Dir(filename)
	f.Content = data

	err := fs.SaveFile(&f, filename)
	return err
}

func ReadFile(fs FileSys, path string) ([]byte, error) {

	file, err := fs.GetFile(path)
	if err != nil {
		return []byte{}, err
	}
	return file.Content, err

	//
	// ====================
	// pointless detour via read-seak-closer:
	rsc, err := fs.OpenVFS(path)
	if err != nil {
		return []byte{}, err
	}

	defer rsc.Close()
	b, err := ioutil.ReadAll(rsc)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

// satisfy vfs.Opener
// Conflicts with Afero Open method
func (fs FileSys) OpenVFS(path string) (vfs.ReadSeekCloser, error) {

	var b []byte

	file, err := fs.GetFile(path)
	if err != nil {
		return NopCloser(bytes.NewReader(b)), err
	}

	b = file.Content
	return NopCloser(bytes.NewReader(b)), nil
}

func (fs FileSys) Lstat(path string) (os.FileInfo, error) {
	panic(spf("Links not implemented for %v", fs))
	var fi os.FileInfo
	return fi, nil
}

func (fs FileSys) Stat(path string) (os.FileInfo, error) {
	f, err := fs.GetFile(path)
	if err != nil {
		dir, err := fs.GetDirByPath(path)
		if err != nil {
			return nil, err
		}
		return os.FileInfo(dir), nil
	} else {
		return os.FileInfo(f), nil
	}
}

// satisfying vfs interface
// and ioutil.ReadDir
// similar to GetFiles, but returning only dirs
// Todo: Sort dirs by name
func (fs *FileSys) ReadDir(path string) ([]os.FileInfo, error) {

	path = cleanseLeadingSlash(path)

	var dirs []Directory
	var fis []os.FileInfo

	dir, err := fs.GetDirByPath(path)
	if err == datastore.ErrNoSuchEntity {
		return fis, err
	} else if err != nil {
		logif.E(err)
		return fis, err
	}

	q := datastore.NewQuery(tdir).Ancestor(dir.Key)
	keys, err := q.GetAll(fs.Ctx(), &dirs)
	_ = keys
	if err != nil {
		fs.Ctx().Errorf("Error getching dir children of %v => %v", dir.Key, err)
		return fis, err
	}

	for _, v := range dirs {
		fi := os.FileInfo(v)
		fis = append(fis, fi)
	}

	return fis, err

}

func (fs FileSys) String() string {
	return "gaefs"
}
