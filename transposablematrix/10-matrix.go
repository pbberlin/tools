package transposablematrix

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pbberlin/tools/util"
)

/*
usage:

import "github.com/pbberlin/tools/transposablematrix"

func main() {
	transposablematrix.Main1()
	<-transposablematrix.TermBoxDone
}
*/

func init() {
	rand.Seed(time.Now().UnixNano())
	primes = util.PrecomputePrimes(2560)
}

const (
	IncreaseSlack = false

	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1) // max int domain
	MinInt  = -MaxInt - 1       // min int domain

	multiSortFactor   = 100
	multiSortFactorSq = multiSortFactor * multiSortFactor
)

var (
	primes []bool
)

var (
	lineHoriz = rune(0x2500) //  |
	lineVertc = rune(0x2502) //  -

	fullBlock = rune(0x2588) //  █
	space     = rune(0x0020) // ' '

	leftSquareBracket  = rune(0x005B) // '['
	rightSquareBracket = rune(0x005D) // ']'

	// unused:
	arropUp   = rune(0x005E) // '^'
	arropDown = rune(0x0076) // 'v'

	//                   left top    , right top   , bottom left , bottom right
	runeCorners = []rune{0x250C + 0*4, 0x250C + 1*4, 0x250C + 2*4, 0x250C + 3*4}

	//               vert branch right, v branch lft, horiz br dwn, h. branc top
	runeFittings = []rune{0x251C + 0*8, 0x251C + 1*8, 0x251C + 2*8, 0x251C + 3*8}
)

// from noon opposing clock dir
var spriteCorners1 = []string{
	string(runeCorners[0]) + string(lineHoriz), // ┌-
	string(runeCorners[2]) + string(lineHoriz), // └-
	string(lineHoriz) + string(runeCorners[3]), // -┘
	string(lineHoriz) + string(runeCorners[1]), // -┐
}

var spriteUturnsHoriz = []string{
	string(leftSquareBracket) + string(lineHoriz),  // [-
	string(lineHoriz) + string(rightSquareBracket), // -]
}

// unused
var spriteUturnsVert = []string{
	string(runeCorners[2]) + string(runeCorners[3]), // └┘
	string(runeCorners[0]) + string(runeCorners[1]), // ┌┐
}

var oldCorners = []string{"┌-", "└-", "-┘", "-┐"}
var otherSprites = [][]string{
	{"←", "↑", "→", "↓"},
	{"┤", "┘", "┴", "└", "├", "┌", "┬", "┐"},
	{"|", "/", "-", `\`},
	{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"},
	{"■", "□", "▪", "▫"},
}

// Type Point only when naming becomes too tedious.
// We mostly use explicit x,y variables.
type Point struct {
	x, y int
}

// Slot is conveyor to the underlying Amorph data
type Slot struct {
	AmX    *Amorph
	Label  string // to render meta-information, w + w/o amorph
	Ax, Ay int    // amorph coordinates, left top is at 0;0
}

type TransposableMatrix struct {
	m      [][]Slot // two dimensional plotting matrix
	nx, ny int      // it's lengths: len(m[]), len(m[][])

	cx, cy int // center coords for perspectivized views

	// The following points are maintained during Set()
	// They grow with Set(), but do *not* shrink.
	// We'd have to buffer them before CastStitch() and restore them after DeleteAmorph()

	// Total min and max points:
	xmif, ymif int // min point (left-top)     of rect of *filled* slots
	xmaf, ymaf int // max point (bottom-right) of rect of *non empty* slots

	//
	// Outmost points:
	// One coordinate at the minimum/maximum line of any filled Slot
	// The other coordinates at the endpoints of this line of filled Slots.
	// Could be losely described as the "concavest points"
	//
	// Illustration1:
	//
	//              NNW N NNE
	//           NWW         NEE
	//         	W               E
	//           SWW         SEE
	//              SSW S SSE
	//
	// Illustration2:
	//
	//            NNW<--N-------->NNE
	//
	//
	//     NWW
	//     ^            ^          NEE
	//     |            |            ^
	//     W        <---|--->        E
	//     |            |            v
	//     |            v          SEE
	//     |
	//     v
	//     SWW
	//      SSW<--------S--->SSE
	//
	//
	// Illustration3:
	//                    cx
	//                     |
	//      P0        P1XXX|XXX...
	//                XXXXX|XXX...
	//           XXXXXXXXXX|XXX...
	//           XXXXXXXXXX|XXX...
	//      P2XXXXXXXXXXXXX|XXX...
	//      XXXXXXXXXXXXXXX|XXX...
	//         XXXXXXXXXXXX|XXX...
	//-cy------------------|---...
	//         XXXXXXXXXXXX|XXX...
	//         XXXXXXXXXXXX|XXX...
	//
	// P0 := {xmif,ymif}   => total min point
	// P1  => outmostN[0]  ; P2 =>  outmostW[1]

	NNW Point
	NWW Point
	SWW Point
	SSW Point
	SSE Point
	SEE Point
	NEE Point
	NNE Point

	// Same points redundantly as slice
	// holding NNW, NWW ... counter clockwise
	sOutmost []*Point

	// adding viewport restrictions
	//
	// vly1--------------------------------------
	//     |                                    |
	//     |           NNW<--N-------->NNE      |
	//     |                                    |
	//     |                                    |
	//     |    NWW                             |
	//     |    ^                       NEE     |
	//     |    |                         ^     |
	//     |    W                         E     |
	//     |    |                         v     |
	//     |    |                       SEE     |
	//     |    |                               |
	//     |    v                               |
	//     |    SWW                             |
	//     |     SSW<--------S--->SSE           |
	//     |                                    |
	// vly2--------------------------------------
	//    vlx1                                vlx2

	vPRx1, vPRx2 int
	vPRy1, vPRy2 int
	enforceVPR   bool // enforce view limits

	persp int
}

func NewTransposableMatrix(argNx, argNy int) *TransposableMatrix {
	m := TransposableMatrix{}

	m.nx = argNx
	m.ny = argNy

	m.m = make([][]Slot, m.nx)
	for x := 0; x < m.nx; x++ {
		m.m[x] = make([]Slot, m.ny)
	}

	m.cx = m.nx / 2 // centering in the middle
	m.cy = m.ny / 2

	m.xmif = m.nx - 1 // init total min/max
	m.ymif = m.ny - 1

	m.xmaf = 0
	m.ymaf = 0

	m.sOutmost = make([]*Point, 8) // init outmost
	m.sOutmost[0] = &m.NNW         // thus we have verbose naming - and shift-rotatable access - without any new memory allocation during Set()
	m.sOutmost[1] = &m.NWW
	m.sOutmost[2] = &m.SWW
	m.sOutmost[3] = &m.SSW
	m.sOutmost[4] = &m.SSE
	m.sOutmost[5] = &m.SEE
	m.sOutmost[6] = &m.NEE
	m.sOutmost[7] = &m.NNE

	return &m
}

// vPRM.. - viewport restriction mapped coord
// vPRB.. - viewport restriction base   coord
func (m *TransposableMatrix) ViewportRestriction(vPRMx1, vPRMy1, vPRMx2, vPRMy2 int) {

	vPRBx1, vPRBy1 := m.transposeMapped2Base(vPRMx1, vPRMy1)
	vPRBx2, vPRBy2 := m.transposeMapped2Base(vPRMx2, vPRMy2)

	m.vPRx1 = vPRBx1
	m.vPRy1 = vPRBy1
	m.vPRx2 = vPRBx2
	m.vPRy2 = vPRBy2

	m.enforceVPR = true
}

// Move the center for mapped coordinates Get()/Set()
func (m *TransposableMatrix) SetCenter(x, y int) {
	newX, newY := m.nx/2+x, m.ny/2+y
	m.cx = newX
	m.cy = newY
}

func (m *TransposableMatrix) Center() (int, int) {
	return m.cx, m.cy
}

func (m *TransposableMatrix) SetPerspective(p int) {
	if !(p >= 0 && p <= 3) {
		panic(fmt.Sprintf("allowed: 0 - North - 1 West - 2 South - 3 East; is %v", p))
	}
	m.persp = p
}

func (m *TransposableMatrix) Perspective() int {
	return m.persp
}

func (m *TransposableMatrix) Slots() [][]Slot {
	return m.m
}

// TotalMinMax returns the outmost filled coordinates
func (m *TransposableMatrix) TotalMinMax(printIt bool) (int, int, int, int) {

	if printIt {
		m.totalMinMaxPrint()
	}

	xm0, ym0 := m.transposeBase2Mapped(m.xmif, m.ymif)
	xm1, ym1 := m.transposeBase2Mapped(m.xmaf, m.ymaf)
	switch m.persp {
	case 0:
		return xm0, ym0, xm1, ym1
	case 1:
		return xm1, ym0, xm0, ym1
	case 2:
		return xm1, ym1, xm0, ym0
	case 3:
		return xm0, ym1, xm1, ym0
	default:
		panic("allowed: 0 - North - 1 West - 2 South - 3 East")
	}

}

// totalMinMaxPrint() *fills* current total min/max lines intersect
func (m *TransposableMatrix) totalMinMaxPrint() {
	persp := 4 - m.persp
	m.m[m.xmif][m.ymif].Label = spriteCorners1[(persp+0)%4]
	m.m[m.xmif][m.ymaf].Label = spriteCorners1[(persp+1)%4]
	m.m[m.xmaf][m.ymaf].Label = spriteCorners1[(persp+2)%4]
	m.m[m.xmaf][m.ymif].Label = spriteCorners1[(persp+3)%4]
}

// FilledMinMax returns 8 outmost points, perspectivized
// FilledMinMax() shift-rotates the slice of outmost points
// so that from we start with the NNW concavest - in every given perspecitve
// The parameter printIt makes the points visible
//  - setting their slots to a sprite
func (m *TransposableMatrix) FilledMinMax(printIt bool) []*Point {

	ret := make([]*Point, 8)
	copy(ret, m.sOutmost)

	f := func(steps int) {
		for i := 0; i < steps; i++ {
			ret = append(ret[1:], ret[0])
		}
	}

	switch m.persp {
	case 0:
		// nothing
	case 1:
		f(2)
	case 2:
		f(4)
	case 3:
		f(6)
	default:
		panic("allowed: 0 North - 1 West - 2 South - 3 East")
	}

	if printIt {
		alreadyXY := map[int]bool{}
		for i := 0; i < len(ret); i++ {
			xm, ym := m.transposeBase2Mapped(ret[i].x, ret[i].y)
			m.SetLabel(xm, ym, Slot{Label: spriteCorners1[i/2]})
			if alreadyXY[multiSortFactor*xm+ym] {
				m.SetLabel(xm, ym, Slot{Label: spriteUturnsHoriz[i/4]})
			}
			alreadyXY[multiSortFactor*xm+ym] = true
		}
	}

	return ret
}

// updateLimits() is called upon each Set()
// updating total and outmost coordinates
func (m *TransposableMatrix) updateLimits(xb, yb int) {

	// total limits of filled slots
	if xb < m.xmif {
		m.xmif = xb
	}
	if xb >= m.xmaf {
		m.xmaf = xb
	}
	if yb < m.ymif {
		m.ymif = yb
	}
	if yb >= m.ymaf {
		m.ymaf = yb
	}

	//
	// the outmost slots - along horizontal perspectives
	if xb == m.xmif &&
		(yb < m.NWW.y || xb != m.NWW.x) { // second comparison: new total min line - xb < NWW.x no worki - since NWW.x inits to 0
		m.NWW.x, m.NWW.y = xb, yb
	}

	if xb == m.xmaf &&
		(yb < m.NEE.y || xb != m.NEE.x) {
		m.NEE.x, m.NEE.y = xb, yb
	}

	if xb == m.xmif &&
		(yb > m.SWW.y || xb != m.SWW.x) {
		m.SWW.x, m.SWW.y = xb, yb
	}

	if xb == m.xmaf &&
		(yb > m.SEE.y || xb != m.SEE.x) {
		m.SEE.x, m.SEE.y = xb, yb
	}

	//
	// the outmost slots - along vertical perspectives
	if yb == m.ymif &&
		(xb < m.NNW.x || yb != m.NNW.y) {
		m.NNW.x, m.NNW.y = xb, yb
	}

	if yb == m.ymif &&
		(xb > m.NNE.x || yb != m.NNE.y) {
		m.NNE.x, m.NNE.y = xb, yb
	}

	if yb == m.ymaf &&
		(xb < m.SSW.x || yb != m.SSW.y) {
		m.SSW.x, m.SSW.y = xb, yb
	}

	if yb == m.ymaf &&
		(xb > m.SSE.x || yb != m.SSE.y) {
		m.SSE.x, m.SSE.y = xb, yb
	}

}

// Get() gives us the mapped slot - observing perspective and center
// --x-       ---
// ----   /\  x--
// ----   _|  ---
//            ---
func (m *TransposableMatrix) Get(x, y int) Slot {
	xm, ym := m.transposeMapped2Base(x, y)
	if xm == -1 || ym == -1 {
		return Slot{Label: "~~"}
	}
	return m.m[xm][ym]
}

// Empty reports the state of a certain coord.
// If the label is filled, but there is no amorph
// empty returns true
func (m *TransposableMatrix) Empty(x, y int) bool {
	sl := m.Get(x, y)
	if sl.AmX == nil {
		return true
	}
	return false
}

// SetLabel() does not touch the *Slot.Amorph, only changes label
func (m *TransposableMatrix) SetLabel(x, y int, slt Slot) {
	before := m.Get(x, y)
	before.Label = slt.Label
	m.Set(x, y, before)
}

// Set() changes the mapped slot
// We could trigger a panic on override
// to prevent erroneous usages.
// Override would then require a preceding Delete()
func (m *TransposableMatrix) Set(x, y int, slt Slot) {

	xb, yb := m.transposeMapped2Base(x, y)
	// pf("%v %v  => %v %v\n", x, y, xm, ym)
	if xb == -1 || yb == -1 { // cant set
		return
	}
	m.m[xb][yb] = slt

	if slt.AmX != nil {
		m.updateLimits(xb, yb)
	}

}

func (m *TransposableMatrix) Delete(x, y int) {
	panic("delete is impossible, since it cannot restore xmif-xmaf NNW-NNE...")
	xb, yb := m.transposeMapped2Base(x, y)
	if xb == -1 || yb == -1 {
		return
	}
	m.m[xb][yb] = Slot{} // empty slot
}

// transposeMapped2Base() translates
// mapped, "in perspective" coords (xm,ym)
// to *real* slice coords, called base coords xb,yb
// -1 means out of range (too little or too large)
func (m *TransposableMatrix) transposeMapped2Base(xm, ym int) (xb, yb int) {

	lx := len(m.m) - 1
	ly := len(m.m[0]) - 1

	xb, yb = -1, -1

	switch m.persp {
	case 0:
		xm += m.cx
		ym += m.cy
		if ym > ly || xm > lx || ym < 0 || xm < 0 {
			xb, yb = -1, -1
			break
		}
		// return m.m[xm][ym]
		xb, yb = xm, ym
	case 1:
		xm -= m.cy
		ym += m.cx
		if ym > lx || -xm > ly || ym < 0 || -xm < 0 {
			xb, yb = -1, -1
			break
		}
		// return m.m[ym][-xm]
		xb, yb = ym, -xm
	case 2:
		xm -= m.cx
		ym -= m.cy
		if -xm > lx || -ym > ly || -xm < 0 || -ym < 0 {
			xb, yb = -1, -1
			break
		}
		// return m.m[-xm][-ym]
		xb, yb = -xm, -ym
	case 3:
		xm += m.cy
		ym -= m.cx
		if -ym > lx || xm > ly || -ym < 0 || xm < 0 {
			xb, yb = -1, -1
			break
		}
		// return m.m[-ym][xm]
		xb, yb = -ym, xm
	default:
		panic("allowed: 0 - North - 1 West - 2 South - 3 East")
	}

	return

}

// transposeBase2Mapped() is the Inverse() of transposeMapped2Base()
func (m *TransposableMatrix) transposeBase2Mapped(xb, yb int) (xm, ym int) {

	xm, ym = -1, -1
	switch m.persp {
	case 0:
		xb -= m.cx
		yb -= m.cy
		xm, ym = xb, yb
	case 1:
		xm, ym = -yb, xb
		xm += m.cy
		ym -= m.cx
	case 2:
		xm, ym = -xb, -yb
		xm += m.cx
		ym += m.cy
	case 3:
		xm, ym = yb, -xb
		xm -= m.cy
		ym += m.cx
	default:
		panic(fmt.Sprintf("allowed: 0 - North - 1 West - 2 South - 3 East; is %v", m.persp))
	}

	return

}

// DrawLine only writes to the slot label.
// The pointer to the amorph remains unchanged
// Empty() still reports true for the slot
func (m *TransposableMatrix) DrawLine(line []Point, suffix string, pivotPointsOnly bool) {
	for i := 0; i < len(line); i++ {
		s := spf("%v%v", i%10, suffix)
		if suffix == "" {
			s = ""
		}
		if pivotPointsOnly || i == len(line)-1 {
			m.SetLabel(line[i].x, line[i].y, Slot{Label: s})
		} else {
			x := util.Min(line[i+1].x, line[i].x)
			y := util.Min(line[i+1].y, line[i].y)
			dx := util.Abs(line[i+1].x - line[i].x)
			dy := util.Abs(line[i+1].y - line[i].y)
			// pf("sect :%v %v %v %v \n", x, x+dx, y, y+dy)
			for j := x; j <= x+dx; j++ {
				for k := y; k <= y+dy; k++ {
					m.SetLabel(j, k, Slot{Label: s})
				}
			}
		}
	}
}

func (m *TransposableMatrix) SetAmorphSnapLeftwise(base Point, am Amorph) {
	baseTransformed := Point{base.x - am.Cols, base.y - am.Rows}
	m.CastStitch(baseTransformed, am)
}
func (m *TransposableMatrix) SetAmorphSnapRightwise(base Point, am Amorph) {
	// don't know what xm0 + 1 is needed for
	baseTransformed := Point{base.x, base.y - am.Rows}
	m.CastStitch(baseTransformed, am)
}

// CastStitch() paints the elements of an amorph into the matrix.
// The amorph area can be *kerfed* southwest and southeast like this:
//    |XXXXXXXXXXXXXXX|
//    |XXXXXXXXXXXXXXX|
//    |XXXXXXXXXXXP   |
//    |XXXXXXXXXXX    |
//    |   PXXXXXXX    |
//    |    XXXXXXX    |
//    |    XXXXXXX    |
// The kerfs have positive coords, subtracted from the edges.
// The amorph is consumed, casted/molded into a permanent shape
// and finally it is stitched to the nearby amorphs
func (m *TransposableMatrix) CastStitch(base Point, am Amorph) {

	Spent[am.IdxArticle] = true

	var rightKerf, leftKerf Point

	if len(am.Edge) == 3 {
		if am.Edge[1] > 0 {
			rightKerf.x = am.Edge[2]
			rightKerf.y = am.Edge[1]
		} else if am.Edge[1] < 0 {
			leftKerf.x = am.Edge[0]
			leftKerf.y = -am.Edge[1]
		}
	}

	maxSlots := am.Cols*am.Rows - am.Slack // 3*3 -2 => 7
	cntrSlots := 0
	cntrInfo := 0 // which property to display - width, height...

	// debug alterations:
	if IncreaseSlack {
		maxSlots += 111
	}

	// leftKerf.x, leftKerf.y = 1, 2
	// rightKerf.x, rightKerf.y = 0, 0

lbl0:
	for i3 := 0; i3 < am.Cols; i3++ { //  columnwise

		limRow := am.Rows
		if i3 < leftKerf.x {
			limRow = am.Rows - leftKerf.y
		}
		if i3 > am.Cols-rightKerf.x-1 {
			limRow = am.Rows - rightKerf.y
		}
		for i4 := 0; i4 < limRow; i4++ { // down the rows - columnwise
			x := base.x + i3
			y := base.y + i4

			outr1 := false
			if i4 == 0 || i4 == am.Rows-1 || i3 == 0 || i3 == am.Cols-1 {
				outr1 = true
			}
			outr2 := false
			if i4 == 1 || i4 == am.Rows-2 || i3 == 1 || i3 == am.Cols-2 {
				outr2 = true
			}

			slt := Slot{}
			slt.AmX = &am
			xib, yib := am.transpose(i3, i4, m.persp)
			slt.Ax = xib
			slt.Ay = yib
			if outr1 {
				if i3%2 == 0 && i3 != am.Cols-2 ||
					am.Cols == 2 && i3 == 0 {
					// slt.Label = spf("%1v%1v", am.IdxArticle%10, yib%10)
					slt.Label = spf("%-2v", am.IdxArticle)
				} else {
					slt.Label = "  "
				}
				if i3 == am.Cols-1 {
					slt.Label = spf("%2v", am.IdxArticle)
				}
			} else if outr2 {
				switch cntrInfo % 3 {
				case 0:
					slt.Label = spf("c%v", am.Cols%10)
				case 1:
					slt.Label = spf("r%v", am.Rows%10)
				case 2:
					slt.Label = spf("  ")
				}
				cntrInfo++
			} else {
				slt.Label = "  "
			}

			if am.Padded > 0 {
				slt.Label = spf("p%v", am.Padded)
			}

			// slt.Label = "xx"
			// slt.Label = spf("%1v%1v", xib, yib)
			m.Set(x, y, slt)

			cntrSlots++
			if cntrSlots >= maxSlots {
				break lbl0
			}

		}
	}

}

// PrintMatrixBuffer prints those slots,
// that are not empty
// It prints the base coords
func PrintMatrixBuffer(m TransposableMatrix) {

	sl := m.Slots()
	for k1, ssl := range sl {
		hasCont := false
		for k2, slot := range ssl {
			if slot.AmX != nil {
				pf("%v-%v: %v%v |", k1, k2, slot.AmX.Cols, slot.AmX.Rows)
				hasCont = true
			}
		}
		if hasCont == true {
			pf("\n")
		}
	}

}
