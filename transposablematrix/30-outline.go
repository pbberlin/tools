package transposablematrix

import "sort"

type SectionBoundaryTrajectory int //
const (
	Convex SectionBoundaryTrajectory = iota
	Concave
)

const fictSegLen = 6 // fictitious segment length

var BottomlessPitError = epf("outline: bottomless pit")

func (a SectionBoundaryTrajectory) String() string {
	if a == Convex {
		return "convx"
	}
	return "cncav"
}

// PartialOutline() computes a slice of Points
// starting at NWW of filled slots, stopping at NEE.
// PartialOutline() adds fictitious sections.
// a.) to force viewport restrictions
// b.) to give FuseTwoSections a two-sect outline at bootstrap
func (m *TransposableMatrix) PartialOutline() ([]Point, error) {

	ret := m.FilledMinMax(false)

	p1 := ret[1] // NWW
	p2 := ret[6] // NEE

	xm1, ym1 := m.transposeBase2Mapped(p1.x, p1.y)
	xm2, ym2 := m.transposeBase2Mapped(p2.x, p2.y)

	// completely empty yet?
	if xm1 == xm2 && ym1 == ym2 {
		// pf("%v %v", p1, p2, xm1, ym1,xm2, ym2)
		line := m.enhanceLine(true, []Point{}, xm1, xm2)
		return line, nil
	}

	// m.SetLabel( xm1, ym1-1, Slot{Label:"--"})
	// m.SetLabel( xm2, ym2-1, Slot{Label:"--"})

	// init the loop
	y := ym1
	direction := 0 // 0 up, 1 right, 2 down
	prevDir := 0   // previous direction, init causing direction change at start
	line := []Point{Point{xm1, ym1}}
	cntr := 0

	breakExhaustive := true
	for x := xm1; x <= xm2+1; x++ {

		cntr++
		if cntr > m.nx+m.ny {
			err := epf("break outline imperfect - cntr%v - len%v %v \n", cntr, len(line), line)
			return nil, err
		}

		NE := !m.Empty(x, y-1)
		SE := !m.Empty(x, y)

		// three principal directions
		switch {
		case NE && SE:
			direction = 0 // up
		case !NE && SE:
			direction = 1 // right
		case !NE && !SE:
			direction = 2 // down
			// anomaly: bottomless pit - disk world abyss
			// if y >= 0 {
			if y >= 100 {
				return nil, BottomlessPitError
				direction = 1 //   => go right - this "fix" confused the algorithm
			}
		default: // NE && !SE
			direction = 1 // anomaly: scraping at the lower side => go right
		}
		// if x>4{
		// 	pf("%v,%v;%v| ",x,y,direction)
		// }

		// when direction was changed => plot
		if direction != prevDir {
			switch direction {
			case 0:
				line = appendDistinct(line, Point{x, y})
			case 1:
				if prevDir == 0 {
					line = appendDistinct(line, Point{x, y})
				} else if prevDir == 2 {
					line = appendDistinct(line, Point{x - 1, y})
				}
			case 2:
				line = appendDistinct(line, Point{x - 1, y})
			}

			// regular loop exit
			if len(line) > 0 {
				lastPoint := line[len(line)-1]
				if lastPoint.x == xm2 && lastPoint.y == ym2 {
					// pf("outline break pfct\n")
					breakExhaustive = false
					break
				}
			}
		}
		prevDir = direction

		// adjusting indexes; x gets regularly incremented
		switch direction {
		case 0: // up
			x--
			y--
		case 1: // right
		case 2: // down
			x--
			y++
		}

	}

	if breakExhaustive {
		err := epf("outline break exhaustive %v,%v\n", xm1, xm2)
		return nil, err
	}

	// Chop off double point at line-start
	if len(line) > 1 && line[0].x == line[1].x && line[0].y == line[1].y {
		line = line[1:]
	}

	line = m.enhanceLine(false, line, xm1, xm2)
	return line, nil

}

// ConcaveSections assumes mapped coordinates,
// centered 0;0 x growing positive, y growing negative.
// Consideration rests on the y-coords.
// ConcaveSections() ferrets out concave sections
// based on following consideration:
// A preceding downward slope, spanning several points
// An upward slope, of even the smallest extend.
// Each upward slope erases previous downslopes.
//
// ConcaveSections returns indexes of Points of "floor" sections
// Thus the neighboring sections can be pursued using the original outline
func (m *TransposableMatrix) ConcaveSections(outline []Point) [][]int {

	sections := [][]int{}

	if len(outline) < 2 { // too short to have concavity
		return sections
	}

	var P1, P2, P3 *Point
	idxP2, idxP3 := -1, -1

	for i := 1; i < len(outline); i++ { // starting at second point
		x := outline[i].x
		y := outline[i].y

		switch {
		case P1 == nil:
			P1 = &Point{x, y} // first
		case P1 != nil && P2 == nil && y <= P1.y:
			P1 = &Point{x, y} // north or east move - no downslope
		case P1 != nil && P2 == nil && y > P1.y:
			P2 = &Point{x, y} // initial south move
			idxP2 = i
		case P1 != nil && P2 != nil && y > P2.y:
			P1 = P2           // continuous south move - update
			P2 = &Point{x, y} //
			idxP2 = i
		case P1 != nil && P2 != nil && y == P2.y:
			P3 = &Point{x, y} // east move - with previous slope
			idxP3 = i
		case P1 != nil && P2 != nil && P3 != nil && y < P3.y:
			sections = append(sections, []int{idxP2, idxP3})
			P1, P2, P3 = nil, nil, nil
			idxP2, idxP3 = -1, -1
		}

	}
	//pf("%v \n", sections)

	return sections

}

// ConcavestLowestNarrowestSections() combines Parialoutline() and ConcaveSections()
// It takes the perspective outline,
// and extracts all horizontal sections of it.
// The sections are sorted by being concave, low, short.
// Hereby we rely on ConcaveSections() to put concave first.
// Finally, the lowest smallest sections
// are returned as indizes to points in the original outline.
func (m *TransposableMatrix) ConcavestLowestNarrowestSections() ([][]int, []Point, error) {

	l, err := m.PartialOutline()

	mappedByY := make(map[int][]int, 10) // careful about double keys

	cs := m.ConcaveSections(l)
	for i := 0; i < len(cs); i++ {
		pIdx1 := cs[i][0]
		pIdx2 := cs[i][1]
		P1 := l[pIdx1]
		P2 := l[pIdx2]
		mappedByY[sortFactor(true, P1, P2, mappedByY)] = []int{pIdx1, pIdx2}
	}

	csIndex := map[int]int{}
	for _, v := range cs {
		csIndex[v[0]] = v[1]
	}

	// pf("concave_sects    : %v \n", mappedByY)
	// pf("concave_sects idx: %v \n", csIndex)

	for k := 1; k < len(l); k++ {
		P2 := l[k]
		P1 := l[k-1]
		if P1.y == P2.y { // => horiz segments
			// pf("%c %v,%v;%v | ", 'a'+(k-1), w.x, v.x, v.y)
			alreadyConcave := (csIndex[k-1] == k)
			if !alreadyConcave {
				accommodatedSortFact := sortFactor(false, P1, P2, mappedByY)
				mappedByY[accommodatedSortFact] = []int{k - 1, k}
				// pf("added   %2v %2v - %3v,%3v,%3v  with            %v \n", k-1,k, P1.y, P1.x, P2.x , accommodatedSF)
			} else {
				// pf("skipped %2v %2v - %3v,%3v,%3v  already concave %v \n", k-1,k,P1.y, P1.x, P2.x , alreadyConcave)
			}
		}
	}

	sortedKeys := make([]int, 0, len(mappedByY))
	for k, _ := range mappedByY {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)

	ret := make([][]int, len(mappedByY))
	for i, sk := range sortedKeys {
		v := mappedByY[sk]
		ret[i] = v
	}

	ret = decorateClns(ret, l)

	return ret, l, err
}

// decorateClns() takes a slice of concavest, lowest narrowest sections (clns)
// and appends (decorates) their eastern and western dy,dx lengths.
// The return slice has format
// [idx section][ idxP1 idxP2 DX | westDY westDX  2westDY | eastDY eastDX 2eastDY ].
// By convention ...DX is always positive.
// ...DY > 0 indicates upward DY, DY < 0 indicates downward precipice
func decorateClns(clns [][]int, l []Point) [][]int {
	r := clns
	for i := 0; i < len(clns); i++ {
		pIdx1 := clns[i][0]
		pIdx2 := clns[i][1]

		var P1, P11, P12, P13, P2, P21, P22, P23 Point
		var dxbase, dyw, dxw, dye, dxe int
		var dyw2, dye2 int

		P1 = l[pIdx1]
		if pIdx1-1 >= 0 {
			P11 = l[pIdx1-1]
			dyw = (P1.y - P11.y)
		}
		if pIdx1-2 >= 0 {
			P12 = l[pIdx1-2]
			dxw = P11.x - P12.x
		}
		if pIdx1-3 >= 0 {
			P13 = l[pIdx1-3]
			dyw2 = (P12.y - P13.y)
		}

		P2 = l[pIdx2]
		if pIdx2+1 < len(l) {
			P21 = l[pIdx2+1]
			dye = (P2.y - P21.y)
		}
		if pIdx2+2 < len(l) {
			P22 = l[pIdx2+2]
			dxe = P22.x - P21.x
		}
		if pIdx2+3 < len(l) {
			P23 = l[pIdx2+3]
			dye2 = (P22.y - P23.y)
		}

		dxbase = P2.x - P1.x
		// The outline goes like this:
		//
		//                  PXP
		//                  XXX
		// XXP            PXPXPXX
		// XXX            XXXXXXX
		// XXPXXXXXXXXXXXXPXXXXXX
		//
		// Therefore we have to make corrections:
		if dyw <= 0 && dye <= 0 {
			dxbase++ //  convex sourrounding
		}
		if dyw > 0 && dye > 0 {
			dxbase-- // concave sourrounding
		}

		// decr incr for the neighboring x-sections
		// but dyw inverted - 'cause relative to *this* section
		if dyw > 0 && dyw2 <= 0 { // dyw2 == 0 meaning end of world => convex surrounding
			dxw++
		}
		if dyw < 0 && dyw2 > 0 {
			dxw--
		}

		if dye > 0 && dye2 <= 0 {
			dxe++
		}
		if dye < 0 && dye2 > 0 {
			dxe--
		}

		// pf("%v %v %v . %v %v %v \n", P1, P11, P12, P2, P21, P22)

		r[i] = append(r[i], dxbase,
			dyw, dxw, dyw2,
			dye, dxe, dye2)
	}
	return r
}

// appendDistinct prevents double points at line start.
// appendDistinct removes double points
// from "one slot" u-turns
// up=>right=>down:
//
//            X
//   XXXXXXXXXXXXXXXXX
//
// However this *broke* the rule that horizontal segments start with even slice indize!
// Therefore disabled.
// The double point at line start is now removed at the end of PartialOutline()
func appendDistinct(l []Point, p Point) []Point {

	// DISABLED
	if false {
		if len(l) > 0 {
			lastPoint := l[len(l)-1]
			if lastPoint.x == p.x && lastPoint.y == p.y {
				// pf("doubler %v %v|", p.x, p.y)
				return l
			}
		}
	}

	l = append(l, p)
	return l
}

// findYourNeighbor() gives the clns index of the leftmost and rightmost section.
// Remember, that clns itself is sorted concave,lowest,narrowest
// -1 indicates: NO neighbor in that direction.
func findYourNeighbor(clns [][]int, sct []int) (int, int) {

	outlSelf := sct[0]       // index to point in outline
	outlPred := outlSelf - 2 // index to predecessor point
	outlSucc := outlSelf + 2 // index to successor

	if outlPred < 0 {
		outlPred = -1
	}

	// clns has one entry for *two* pointss
	if outlSucc/2 > len(clns)-1 {
		outlSucc = -1
	}
	// pf("self%v pred%v succ%v \n", outlSelf, outlPred, outlSucc)

	Pred, Succ := -1, -1 // indize to clns sections

	for i := 0; i < len(clns); i++ {
		// pf("%v", clns[i][0])
		if clns[i][0] == outlPred {
			Pred = i
		}
		if clns[i][0] == outlSucc {
			Succ = i
		}
	}
	// pf("\n")
	return Pred, Succ

}

// sortFactor() gives a weight to horizontal sections of an outline
// First come concave sections, then lower sections.
// Equally low sections are ordered by the width.
// Finally, there might be equal sections, needing differentiation
// Ultrafinally: sometimes we dont want the differentiation mentioned above,
// but merely check existenc, thus accommodateNext tells us,
// whether we wanna grope for empty space
func sortFactor(isConcave bool, P1, P2 Point, m map[int][]int) int {

	accommodateNext := true

	const fac = 1000 // summands may overwhelm the factor - but slightly imperfect sort is no debacle
	rv := fac * fac
	if isConcave {
		rv -= (fac * fac) // concave sections very first
	}

	rv += -(fac * P1.y) // lowest first, -2 => 2000, -1 => 1000, 4 => -4000
	if P1.y == 0 {
		rv += -(0.8 * fac) // make sure, null ordinate gets sorted high enough
	}

	rv += (P2.x - P1.x) // smaller first; 2,3,5  - assuming mapped coordinates, x always inreasing to right

	if accommodateNext {
		for {
			if _, ok := m[rv]; !ok {
				break
			}
			rv++ // key already exists => increase until finding free slot
		}
	}
	// pf("%5v %5v %5v %2v => %7v \n", accommodateNext, isConcave, fac*P1.y, P2.x-P1.x, rv)
	// if !accommodateNext {
	// 	pf("%5v %5v %5v => %7v \n", P1.y, P2.x, P1.x, rv)
	// }
	return rv
}

func PrintCLNS(clns [][]int) {
	pf("clns: ")
	newL := false
	cntr := 0
	const nCols = 6
	for i := 0; i < len(clns); i++ {
		PrintSectOfCLNS(clns[i])
		newL = false
		if cntr%nCols == nCols-1 {
			pf("\n")
			newL = true
		}
		cntr++
	}
	if !newL {
		pf("\n")
	}
}

// PrintSectOfCLNS prints ONE, not all
func PrintSectOfCLNS(clnsX []int) {
	pf("%v;%v|", clnsX[0], clnsX[1])
	// pf("xyx% 2v % 2v % 2v|", clnsX[2], clnsX[3], clnsX[6])
	pf("%v;%v;%v|", clnsX[2], clnsX[3], clnsX[6])
	// pf("w%v;%v|", clnsX[4], clnsX[5])
	// pf("e%v;%v| ", clnsX[7], clnsX[8])
	pf("%v%v|", clnsX[4], clnsX[5])
	pf("%v%v", clnsX[7], clnsX[8])
	pf("  ")
}

func PrintOutline(l []Point) {
	for i := 0; i < len(l); i++ {
		p := l[i]
		pf("%v:%v,%v|", i, p.x, p.y)
	}
	pf("\n")
}
