package transposablematrix

import "github.com/pbberlin/tools/util"

type CurveDesc int

const (
	cncave CurveDesc = iota
	// cncveCore  // basic indicator: do we have concave core somewhere
	stairW
	stairE
	convex
)

func (c CurveDesc) String() string {
	switch c {
	case cncave:
		return "cncave"
	// case cncveCore:
	// 	return "cncveCore"
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

	w, e            []int     // the lengths of the west/eastward continuations edges: directions y-x-y
	pm              []int     // permissiveness westwards/eastwards/northwards
	curveDesc       CurveDesc // strongly concave, weakly concave or convex surroundings
	concaveCore     bool      // any concave y-x-y ? Turned out, the heuristics applicability depends on that
	dirIdx, maxOffs int       // direction to grow into, amount to grow
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

	dir := "UNKNOWN"
	switch f.dirIdx {
	case -1:
		dir = "blckd"
	case 0:
		dir = "westw"
	case 1:
		dir = "eastw"
	}

	if f.dirIdx > -1 {
		pf("%v%v | ", dir, f.maxOffs)
	} else {
		pf("")
	}

	pf("\n")

}

// FuseTwoSections takes an outline with all concavest lowest narrowest sections,
// and one or two designated section indize out of clns.
// So far mostly *one* section is given, while the second is the
// computed right neighbor. Rightmost sect is complemented by the left neighbor.
// Those designated sections are fused.
// Param l    => outline
// Param clns => concavest, lowest, narrowest sections
// sct1       => index to clns for which to find a pair
// Returns: Fusion
func (ar *Reservoir) FuseTwoSections(l []Point, clns [][]int, sct1 []int) (Fusion, error) {

	fs := NewFusion()
	var err error

	// find the two sections - westward or eastward
	west, east := findYourNeighbor(clns, sct1)
	var sct2 []int
	if sct1[1] == 0 {
		// only fuse eastw possible
		if east < 0 {
			err = epf("NO EASTW NEIGHBOR 1 to %v", sct1)
			return fs, err
		}
		sct2 = clns[east]
	} else if sct1[1] == len(l)-1 {
		// only fuse westw possible
		if west < 0 {
			err = epf("NO WESTW NEIGHBOR 1 to %v", sct1)
			return fs, err
		}
		sct2 = sct1
		sct1 = clns[west]
	} else {
		if util.Abs(sct1[3]) > 0 &&
			util.Abs(sct1[3]) < util.Abs(sct1[6]) {
			// east flank higher than west flank?
			// => fuse westwards
			sct2 = sct1
			if west < 0 {
				err = epf("NO WESTW NEIGHBOR 2 to %v", sct1)
				return fs, err
			}
			sct1 = clns[west]
		} else {
			// fuse eastwards
			if east < 0 {
				err = epf("NO EASTW NEIGHBOR 2 to %v", sct1)
				return fs, err
			}
			sct2 = clns[east]
		}
	}
	// PrintSectOfCLNS(sct1)
	// PrintSectOfCLNS(sct2)

	fs.idxL = append(fs.idxL, sct1[0], sct1[1], sct2[0], sct2[1])
	fs.xyx = append(fs.xyx, sct1[2], sct1[6], sct2[2])
	fs.w = append(fs.w, sct1[3], sct1[4], sct1[5])
	fs.e = append(fs.e, sct2[6], sct2[7], sct2[8])

	sdgN := ar.SmallestDesirableHeight
	sdgWE := ar.SmallestDesirableWidth

	pW, pE, pN := Permissiveness(sdgN, sdgWE, fs.xyx, fs.w, fs.e)
	fs.pm = []int{pW, pE, pN}

	//
	// base point
	// always bottom left
	fs.base.x = l[sct1[0]].x             // x-coord taken from beginning
	baseY1 := l[sct1[0]].y               // y-coord taken either from beginning
	baseY2 := l[sct2[1]].y               //      or taken from end
	fs.base.y = util.Max(baseY1, baseY2) // take lowest y - meaning Max(), not Min()
	dyw := sct1[3]                       // correction for concave angles
	if dyw > 0 {
		fs.base.x++
	}

	switch {
	case fs.pm[0] < 0 && fs.pm[1] < 0: // eastw, westw blocked, concave
		fs.curveDesc = cncave
		fs.dirIdx, fs.maxOffs = -1, 0
	case fs.pm[0] > 0 && fs.pm[1] < 0: //  westwards
		fs.curveDesc = stairW
		fs.dirIdx, fs.maxOffs = 0, fs.pm[0]
	case fs.pm[0] < 0 && fs.pm[1] > 0: //  eastwards
		fs.curveDesc = stairE
		fs.dirIdx, fs.maxOffs = 1, fs.pm[1]
	case fs.pm[0] > 0 && fs.pm[1] > 0: // utterly convex
		fs.curveDesc = convex
		fs.dirIdx, fs.maxOffs = 1, fs.pm[1] // wanton choice: grow east
	}

	lowerX := 0
	if fs.xyx[1] < 0 {
		lowerX = 1
	}
	if lowerX == 0 && fs.w[0] > 0 {
		fs.concaveCore = true
	}
	if lowerX == 1 && fs.e[0] > 0 {
		fs.concaveCore = true
	}

	return fs, nil
}

// Permissiveness collects information on permissiveness to west, east and north
// Parameters sdgN, sdgWE => smallest desirable gap north, west+eastwards
func Permissiveness(sdgN, sdgWE int, xyx, w, e []int) (pW, pE, pN int) {

	yw, xw2, yw2 := w[0], w[1], w[2]
	ye, xe2, ye2 := e[0], e[1], e[2]

	pW = PermissivenessEastOrWest(sdgWE, yw, xw2, yw2)
	pE = PermissivenessEastOrWest(sdgWE, ye, xe2, ye2)

	pNW := PermissivenessNorth(sdgN, yw, xw2, yw2)
	pNE := PermissivenessNorth(sdgN, ye, xe2, ye2) + xyx[1]

	pN = pNW
	if pNE < pNW {
		pN = pNE
	}

	return
}

// PermissivenessNorth looks at possible neighboring
// narrow chimneys, and gives a measure to prevent those
func PermissivenessNorth(sdgN, y, x1, y1 int) (p int) {

	if y1 < 0 || y1 == 0 {
		p = 99
	} else if y1 > 0 {
		y1abs := y1 // y, y1 can be negative, x cannot
		if y1 < 0 {
			y1abs = -y1
		}
		switch {
		case x1 >= sdgN: // broad vertical corridor
			p = 50
		case x1 < sdgN && y1abs <= sdgN: // smaller vertical corridor, but not very high
			p = 2
		case x1 < sdgN && y1abs > sdgN: // narrow high vertical corridor, allow low permissiveness
			p = sdgN
		}

	}

	return
}

// PermissivenessEastOrWest investigates possible tubes
// and returns a measure to avoid those
func PermissivenessEastOrWest(sdgWE, y, x1, y1 int) (p int) {

	// pf("y:%v, x1:%v, y1:%v - ", y, x1, y1)
	if y > 0 {
		p = -10 // walled => growth impossible
	} else if y == 0 {
		p = 99 // nothing there => unlimited
	} else if y < 0 {
		yabs := y // y, y1 can be negative, x cannot
		if y < 0 {
			yabs = -y
		}
		switch {
		case y1 > 0: // y<0 and y1>0 indicate a concavity. Should not happen
			p = -11
		case y < 0 && yabs < sdgWE && x1 >= sdgWE: // long horizontal tube
			p = sdgWE
		default:
			// all other cases with y<0, y1<0
			p = 5
		}
	}

	return
}

// FuseAllSections pairs/combines all sections of an outline to Fusions
func (m *TransposableMatrix) FuseAllSections(ar *Reservoir,
	clns [][]int, l []Point) (fs []Fusion, err error) {

	// prevent the rightmost clns being processed twice
	// (from-left-to-right, from-right-to-left)
	spentFusedSects := map[int]bool{}

	for i := 0; i < len(clns); i++ {

		f := Fusion{}

		f, err = ar.FuseTwoSections(l, clns, clns[i])
		if err != nil {
			return []Fusion{}, err
		}

		if spentFusedSects[util.Intslice2Int(f.idxL)] {
			continue
		}
		spentFusedSects[util.Intslice2Int(f.idxL)] = true

		fs = append(fs, f)
		// f.Print()

	}
	// pf("%+v\n", spentFusedSects)
	return fs, nil
}
