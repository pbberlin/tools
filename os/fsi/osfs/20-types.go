package osfs

type osFileSys struct{}

func New() *osFileSys {
	o := &osFileSys{}
	return o
}
