package gaefs

import (
	"path"
	"strings"

	"appengine/datastore"
)

func cleanseLeadingSlash(p string) string {
	p = path.Join(p)
	for {
		if strings.HasPrefix(p, sep) {
			p = p[1:]
		} else {
			break
		}
	}
	return p
}

func dirFromKey(key *datastore.Key) string {

	dir := key.String()
	dir = cleanseLeadingSlash(dir)

	dir = strings.Replace(dir, tdirsep, "", -1) // removing all fsd,
	dirs := strings.Split(dir, sep)

	// logif.Pf("AAAAAAAAAA %#v", dirs)
	// logif.Pf("  ")

	if len(dirs) > 1 {
		dirs = dirs[1:]
		dir = path.Join(dirs...)
		dir = dir + sep
		return dir
	} else {
		return ""
	}

}
