package filesys

import ds "appengine/datastore"

func (fs FileSys) CreateDir(path string) {
}

func (fs FileSys) recurseMkDir(path string) {
	if path == "" {
		return
	}
	par, err := fs.GetFsoByQuery(path)
	_ = par
	if err == ds.ErrNoSuchEntity {
		par := parent(path)
		fs.recurseMkDir(par)
		fs.CreateDir(path)
	}
}
