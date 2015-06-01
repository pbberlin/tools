package transposablematrix

import "github.com/pbberlin/tools/util"

// enhanceLine extends the first segment horizontally westwards.
// It then adds a fictional vertical.
// Finally we add a horizontal margin.
// Eastwards we draw a pit down to vprFloor and then go vertically upwards.
// Param stillEmpty draws a bottom between the two vertical axis.
//
//
// ---┐                 ┌---
//    |                 |
//    |                 |
//    └---------┐       |
//          XXXX|       |
//          XXXX|       |
//              |       |
//              └-------┘
//
func (m *TransposableMatrix) enhanceLine(stillEmpty bool, line []Point,
	xm1, xm2 int) []Point {

	if !m.enforceVPR {
		return line
	}

	// mapping
	vprmx1, vprmy1 := m.transposeBase2Mapped(m.vPRx1, m.vPRy1)
	vprmx2, vprmy2 := m.transposeBase2Mapped(m.vPRx2, m.vPRy2)

	var vpr1, vpr2, vprFloor int // effective limits, depending on perspective
	if m.persp%2 == 1 {
		vpr1, vpr2 = vprmx1, vprmx2
		vprFloor = util.Max(vprmy1, vprmy2)
	} else {
		vpr1, vpr2 = vprmy1, vprmy2
		vprFloor = util.Max(vprmx1, vprmx2)
	}

	if stillEmpty {
		line = []Point{
			Point{vpr1, vprFloor},
			Point{vpr2, vprFloor},
		}
		pf("empty outline => artificial floor %v\n", line)
		xm1 = vpr1 + 1 // so that forthcoming comparisons apply
		xm2 = vpr2 - 1
	}

	// westward extension
	if vpr1 < xm1 {

		if vprFloor == line[0].y {
			line[0].x = vpr1 // simply prolong
		} else {
			// We could step up/down
			// but then we get gaps.
			// Thus we only step up/down eastwards
			line[0].x = vpr1
		}

		// fictitious vertical upward
		line = append([]Point{Point{}}, line...)
		line[0].x = line[1].x
		line[0].y = line[1].y - fictSegLen

		// fictitious westward horiz segment
		line = append([]Point{Point{}}, line...)
		line[0].x = line[1].x - fictSegLen
		line[0].y = line[1].y

	}

	// eastward extension
	if vpr2 > xm2 {

		if xm2 == vpr2-1 {

			// prolong eastwards one more slot
			lastIdx := len(line) - 1
			line[lastIdx].x = vpr2

		} else if vprFloor == line[0].y {

			// prolong eastwards along x-axis
			lastIdx := len(line) - 1
			line[lastIdx].x = vpr2

		} else {

			// make a full pit/bulge to vprFloor:

			// up/down ...
			lastIdx := len(line) - 1
			line = append(line, Point{})
			line[lastIdx+1].x = line[lastIdx].x
			line[lastIdx+1].y = vprFloor

			// ... and east
			lastIdx = len(line) - 1
			line = append(line, Point{})
			line[lastIdx+1].x = vpr2
			line[lastIdx+1].y = line[lastIdx].y

		}

		// fictitious vertical upward
		lastIdx := len(line) - 1
		line = append(line, Point{})
		line[lastIdx+1].x = line[lastIdx].x
		line[lastIdx+1].y = line[lastIdx].y - fictSegLen

		// fictitious eastw horiz
		lastIdx = len(line) - 1
		line = append(line, Point{})
		line[lastIdx+1].x = line[lastIdx].x + fictSegLen
		line[lastIdx+1].y = line[lastIdx].y
	}

	//
	//
	if vpr1 > xm1 {
		// todo: cropping existing outline
	}
	if xm2 > vpr2 {
		// todo: cropping existing outline
	}

	return line
}
