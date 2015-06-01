package transposablematrix

import "strings"
import "github.com/pbberlin/tools/tboxc"
import "github.com/pbberlin/tools/util"
import "github.com/pbberlin/termbox-go"

var margL = 4
var margT = 3

// Render2Termbox prints a transposable matrix into a termbox
// x1  ... matrix coords
// xo  ... output coords
// xoo ... output coords plus header margins
func (m *TransposableMatrix) Render2Termbox(x1, x2, y1, y2 int) {

	m.renderHeaderTB(x1, y1)

	xo1 := RectOutp2.X1
	xo2 := RectOutp2.X2
	yo1 := RectOutp2.Y1
	yo2 := RectOutp2.Y2

	xoo := xo1 + margL + 1 //
	yoo := yo1 + margT + 0 // dont know why

	col := 0
	row := 0

	// pf("\ndumpstrt  %3v %3v | \n",x1+col, y1+row)
	for x := xoo; x <= xo2-2; x += 2 {
		row = 0
		for y := yoo; y <= yo2; y++ {

			sl := m.Get(x1+col, y1+row)
			runes := []rune{fullBlock, fullBlock}

			tboxc.BrushReset()
			if sl.AmX != nil {
				tboxc.BrushByArticleID(sl.AmX.IdxArticle)
				// s := spf("%-2v", sl.AmX.IdxArticle)
				// pf("xfnd %v %v %v|",x,y,s)
			}

			if sl.Label == "" {
				runes = []rune{space, space}
			}
			if sl.Label != "" {
				// the only way to extract *runes*, not bytes from a string
				// is by ranging over it; compare http://blog.golang.org/strings
				//   and look for "range loops", nihongo
				runeCntr := 0 // byteIdx is *not* the number of runes, but the *byte* index
				for byteIdx, runeI := range sl.Label {
					if runeCntr > len(runes)-1 {
						s := spf("lp%v> len(runes)%v! lbl:%q % x \n", byteIdx, len(runes),
							sl.Label, sl.Label)
						pf(s)
						// panic(s)
						break
					}
					runes[runeCntr] = runeI
					runeCntr++
				}
			}

			for i := 0; i < len(runes); i++ {
				tboxc.PriRuneAt(x+i, y, runes[i])
			}
			row++
		}
		col++
	}
	// pf("\ndumpstop  %3v %3v | \n",x1+col, y1+row)

	//
	// conserve the resulting screenbuffer
	cb := termbox.CellBuffer()
	cbc := make([]termbox.Cell, len(cb), len(cb))
	copy(cbc, cb)
	appStageTBDs = append(appStageTBDs, cbc)
	appStageLogs = append(appStageLogs, "")
	currStage++

}

func (m *TransposableMatrix) renderHeaderTB(xscale, yscale int) {

	xo1 := RectOutp2.X1
	xo2 := RectOutp2.X2
	yo1 := RectOutp2.Y1
	yo2 := RectOutp2.Y2

	xoo := xo1 + margL
	yoo := yo1 + margT

	// row1
	tboxc.PriLineAt(xo1, yo1+0, "  x|")
	tboxc.PriLineAt(xo1, yo1+1, "   |")

	for x := xoo; x <= xo2-2; x += 2 {
		if xscale%10 == 0 {
			colHead1 := spf("%2d", (xscale-xscale%10)/10)
			tboxc.PriLineAt(x, yo1+0, colHead1)
		} else {
			tboxc.PriLineAt(x, yo1+0, "  ")
		}
		colHead2 := spf("%2d", util.Abs(xscale%10))
		tboxc.PriLineAt(x, yo1+1, colHead2)
		xscale += 1
	}

	// row3
	tboxc.PriLineAt(xo1, yo1+2, "  y|")
	tboxc.PriLineAt(xoo, yo1+2, spf(strings.Repeat("=", xo2-xo1-margL+1)))

	for y := yoo; y <= yo2; y++ {
		tboxc.PriLineAt(xo1, y, spf("%3d|", yscale))
		yscale++
	}

}
