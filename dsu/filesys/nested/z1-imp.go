package nested

import ds "appengine/datastore"

func (fs Nested) CreateDir(path string) {
}

func (fs Nested) recurseMkDir(path string) {
	if path == "" {
		return
	}
	par, err := fs.dirByPath(path)
	_ = par
	if err == ds.ErrNoSuchEntity {
		par := parent(path)
		fs.recurseMkDir(par)
		fs.CreateDir(path)
	}
}
