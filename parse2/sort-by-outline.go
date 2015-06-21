package parse2

import (
	"strings"

	"github.com/pbberlin/tools/util"
)

type sortByOutline []string

func (s sortByOutline) Len() int {
	return len(s)
}

func (s sortByOutline) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByOutline) Less(i, j int) bool {

	var sortByLevelFirst bool = true

	if sortByLevelFirst {
		lvl1 := strings.Count(s[i], ".")
		lvl2 := strings.Count(s[j], ".")
		if lvl1 < lvl2 {
			return true
		}
		if lvl1 > lvl2 {
			return false
		}
	}

	// A pure number comparison
	// 1.1, 1.2, 2.1, 2.1.1.
	st1 := strings.Split(s[i], ".")
	st2 := strings.Split(s[j], ".")

	for idx, v1 := range st1 {

		if idx > len(st2)-1 {
			// i.e. 2.37.2 > 2.
			return false
		}
		v2 := st2[idx]

		if util.Stoi(v1) < util.Stoi(v2) {
			return true
		}

		if util.Stoi(v1) > util.Stoi(v2) {
			return false
		}

	}

	// i.e. 2 < 2.26.1.1
	return true
}

//
// px prints debug info for the outline sorting
func px(st1, st2 []string, idx int, op string) {

	var ps1, ps2 []byte

	for i := 0; i <= idx; i++ {
		if i > 0 {
			ps1 = append(ps1, '.')
		}
		v1 := "-"
		if i < len(st1)-1 {
			v1 = st1[i]
		}
		ps1 = append(ps1, v1...)

		//
		if i > 0 {
			ps2 = append(ps2, '.')
		}
		v2 := "-"
		if i < len(st2)-1 {
			v2 = st2[i]
		}
		ps2 = append(ps2, v2...)

		if v1 == "-" && v2 == "-" {
			break
		}
	}

	pf("%10v %s %10v    \n", string(ps1), op, string(ps2))

}
