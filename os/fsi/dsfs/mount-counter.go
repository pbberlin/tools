package dsfs

var cntr = 0

func MountPointLast() string {
	ret := cntr - 1
	if ret < 0 {
		ret = 0
	}
	return spf("mnt%02v", ret)
}

func MountPointIncr() string {
	ret := spf("mnt%02v", cntr)
	cntr++

	// deliberate - if confusing; preventing init -1
	return ret
}

func MountPointReset() string {
	cntr = 0
	return spf("mnt%02v", cntr)
}

func MountPointDecr() string {
	ret := cntr - 1
	if ret < 0 {
		ret = 0
	}
	cntr = ret
	return spf("mnt%02v", cntr)
}
