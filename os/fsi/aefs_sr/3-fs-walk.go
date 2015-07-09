package aefs_sr

/*

// SkipDir is an "error", which a walk-function can
// return, in order to signal, that walk should not traverse into this dir.
var SkipDir = errors.New("skip this directory")

// type WalkFunc func(path string, info os.FileInfo, err error) error

// walk recursively descends path, calling walkFn.
func (fs *AeFileSys) walk(path string, info os.FileInfo, walkFn fsi.WalkFunc) error {

	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	fis, err := fs.ReadDir(path)
	// logif.Pf("%11v => %+v", path, fis)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, fi := range fis {
		filename := pth.Join(path, fi.Name())
		fileInfo, err := fs.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != SkipDir {
				return err
			}
		} else {
			err = fs.walk(filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root.
//
// It requires only the FileSystem interface, and is therefore implementation indepdenent.
//
// It is similar to filepath.Walk(root string, walkFunc)
// In contrast to filepath.Walk, it does not enumerate files, only dirs.
// File enumeration can be added in the WalkFunc.
//
// Directories are walked in order of Readdirnames()
//
// Errors that arise visiting directories can be filtered by walkFn.
//
// Walk does not follow symbolic links.
func (fs *AeFileSys) Walk(root string, walkFn fsi.WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		// logif.Pf("walk start error %v", err)
		return walkFn(root, nil, err)
	}
	return fs.walk(root, info, walkFn)
}

// example walk
func ExampleWalk() {

	// first the per-node func:
	exWalkFunc := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, "my secret directory") {
			return SkipDir
		}

		if err == ErrTooLarge {
			return err // calling off the walk
		}

		tp := "file"
		if f.IsDir() {
			tp = "dir "
		}
		fmt.Printf("Visited: %s %s \n", tp, path)
		return nil
	}

	mnt := "mnt00"
	fs := NewAeFs(mnt) // add appengine context
	fsiX := fsi.FileSystem(fs)

	err := fsiX.Walk(mnt, exWalkFunc)
	fmt.Printf("Walk() returned %v\n", err)
}

*/
