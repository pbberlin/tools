package transposablematrix

import "github.com/pbberlin/tools/util"

type CurveDesc int

const (
	cncave CurveDesc = iota
	stairW           // open towards west
	stairE           // open towards east
	convex
)

func (c CurveDesc) String() string {
	switch c {
	case cncave:
		return "cncave"
	case stairW:
		return "stairW"
	case stairE:
		return "stairE"
	case convex:
		return "convex"
	}
	return ""
}

type Fusion struct {
	idxL []int // indexes to the outline. Turned out unused, but kept.
	base Point // point to the leftmost, lowest position of the fusion - for later stitching

	// the lengths of the fusion edges, directions x-y-x
	// x - western horizontal section, y - shared vertical, x - eastern horiz. sect
	xyx []int

	w, e      []int     // the lengths of the west/eastward continuations edges: directions y-x-y
	pm        []int     // permissiveness westwards/eastwards/northwards
	curveDesc CurveDesc // strongly concave, weakly concave or convex surroundings
}

func NewFusion() Fusion {
	ret := Fusion{}
	ret.idxL = make([]int, 0, 4)
	ret.xyx = make([]int, 0, 3)
	ret.w = make([]int, 0, 3)
	ret.e = make([]int, 0, 3)
	ret.pm = make([]int, 0, 3)
	return ret
}

func (f Fusion) Print() {

	if f.idxL == nil {
		return
	}

	pf("fusn:%2v %2v %2v %2v | ", f.idxL[0], f.idxL[1], f.idxL[2], f.idxL[3])
	t1 := spf("%v;%v", f.base.x, f.base.y)
	pf("baseP%6s | ", t1)
	pf("xyx % 3v % 3v % 3v | ", f.xyx[0], f.xyx[1], f.xyx[2])
	pf("w % 2v % 2v % 2v | ", f.w[0], f.w[1], f.w[2])
	pf("e % 2v % 2v % 2v | ", f.e[0], f.e[1], f.e[2])

	pmd := make([]string, len(f.pm))
	for i := 0; i < len(f.pm); i++ {
		if f.pm[i] < 0 {
			pmd[i] = "bl"
		} else {
			pmd[i] = spf("%v", f.pm[i])
		}
	}

	pf("%s | overgr w%v e%v n%v | ", f.curveDesc, pmd[0], pmd[1], pmd[2])

	pf("\n")

}

// FillHeightFloor is the recommended minimal height
// of an amorph to fill in.
//
func (f Fusion) FillHeightFloor() int {

	if f.xyx == nil || len(f.xyx) < 2 {
		return cFusionHeightDefault
	}

	if f.xyx[1] > 0 {
		// lower western flank => max of western or middle flank
		if f.w == nil || len(f.w) < 1 {
			return cFusionHeightDefault
		}
		return util.Max(f.xyx[1], f.w[0])
	} else {
		// lower eastern flank => max of eastern or middle flank
		if f.e == nil || len(f.e) < 1 {
			return cFusionHeightDefault
		}
		return util.Max(-f.xyx[1], f.e[0])
	}

}
