package osfs

type OsFileSys struct{}

func New() *OsFileSys {
	o := &OsFileSys{}
	return o
}
