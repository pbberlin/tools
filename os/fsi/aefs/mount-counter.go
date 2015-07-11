package aefs

var cntr = 0

func MountPointLast() string {
	ret := cntr - 1
	if ret < 0 {
		ret = 0
	}
	return spf("mnt%02v", ret)
}

func MountPointNext() string {
	ret := spf("mnt%02v", cntr)
	cntr++
	return ret
}

func MountPointReset() string {
	cntr = 0
	return spf("mnt%02v", cntr)
}

func MountPointDecr() string {
	cntr--
	return spf("mnt%02v", cntr)
}
