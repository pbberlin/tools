package osfs

import "runtime"

type osFileSys struct {
	replacePath bool
}

func New() *osFileSys {

	repl := false
	if runtime.GOOS == "windows" {
		repl = true
	}

	o := &osFileSys{replacePath: repl}
	return o
}
