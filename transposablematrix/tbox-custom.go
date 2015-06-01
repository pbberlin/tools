package transposablematrix

import "github.com/pbberlin/termbox-go"
import "github.com/pbberlin/tools/tboxc"
import "strings"

type key struct {
	x  int
	y  int
	ch rune
}

var currStage = 0
var appStageTBDs = [][]termbox.Cell{} // Terminal Box Dumps
var appStageLogs = []string{""}       // always one larger than appStageTBDs, containing tail till the next dump
// both slices are appended at the end of Render2Termbox

var ReStatusTo = tboxc.Rect{4, 2, 4 + 80, 2 + 1}       // rectangle status area total
var ReStatusLe = tboxc.Rect{4, 2, 4 + 40, 2}           // left column
var ReStatusRi = tboxc.Rect{4 + 40, 2, 4 + 40 + 40, 2} // right column

var ReStatusLeBo = tboxc.Rect{4, 3, 4 + 40, 3}
var ReStatusRiBo = tboxc.Rect{4 + 40, 3, 4 + 40 + 40, 3}

var RectOutp1 = tboxc.Rect{10, 6, 10 + 40, 6 + 5}
var RectOutp2 = tboxc.Rect{10, 6, 10 + 40, 6 + 5}

var cntrMouseEvents = 0

func painAppStageData() {

	if currStage < 0 {
		currStage = 0
	}

	idxTBD := 0
	if len(appStageTBDs) > 0 {
		idxTBD = currStage % len(appStageTBDs)
	}

	msg := ""

	if len(appStageTBDs) > 0 {
		if len(appStageTBDs[idxTBD]) > 0 {
			frozenBuf := appStageTBDs[idxTBD]
			currntBuf := termbox.CellBuffer()
			dimX, dimY := termbox.Size()
			_, _ = dimX, dimY
			if x1, x2 := termbox.Size(); x1*x2 != len(currntBuf) {
				panic(spf("should be %v - really is %v", x1*x2, len(currntBuf)))
			}
			msg += spf("taking #%v of %v.", idxTBD, len(appStageTBDs)-1)
			// msg += spf("Dump size %v (%v*%v)", len(currntBuf), dimX, dimY)
			copy(currntBuf, frozenBuf)
			termbox.Flush()
		} else {
			// dump is empty;
			// probably the first element of slice
			tboxc.BlankOut(RectOutp2)
			msg += spf("empty Dump #%v", idxTBD)
		}
	} else {
		msg += "no appStageTBDs"
	}

	// logoutput limps one ahead of termbox dumps.
	// i.e. three dumps - four logs
	// With last TBD, we want last log AND the log since

	logTillTBD := spf("%v\n", appStageLogs[idxTBD+0])
	if strings.HasSuffix(logTillTBD, "\n\n") {
		logTillTBD = logTillTBD[0 : len(logTillTBD)-1]
	}
	logSince := "" // since last TBC
	if idxTBD+1 == len(appStageLogs)-1 {
		logSince = spf("log since TBD:%v\n", appStageLogs[idxTBD+1])
	}
	tboxc.PriToRect(RectOutp1, "AD%v %v\n%v%v", idxTBD, msg, logTillTBD, logSince)

}

func paintMain() {

	var desktopBG = termbox.ColorBlack

	termbox.Clear(termbox.ColorDefault, desktopBG)

	x1, y1 := 2, 1
	x2, y2 := 79, 16
	x2 = 140

	xi1, yi1, xi2, yi2 := tboxc.PainRect(x1, y1, x2, y2, 0)

	RectOutp1.X1 = xi1 + 3
	RectOutp1.Y1 = yi1 + 3
	RectOutp1.X2 = xi2 - 3
	RectOutp1.Y2 = yi2 - 3

	eastExtension := 10
	southExtension := 36

	tboxc.PainFatVertLine(xi1+1, 5, 11)
	tboxc.PainFatVertLine(xi2-1, 5, 11)

	tboxc.PainHorzLine(4, xi1, xi2)

	tboxc.PainRect(x2, y1, x2+eastExtension, y2, 1)
	xi1, yi1, xi2, yi2 = tboxc.PainRect(x1, y2, x2+eastExtension, y2+southExtension, 2)

	RectOutp2.X1 = xi1
	RectOutp2.Y1 = yi1 + 1
	RectOutp2.X2 = xi2 - 1
	RectOutp2.Y2 = yi2 - 2

	painAppStageData()

}

func printKeyPress(ev *termbox.Event) {

	tboxc.PriToRect(ReStatusLe, "Key: %-2v %-4s", ev.Key, tboxc.FuncKeyMap(ev.Key))
	tboxc.PriToRect(ReStatusRi, "Char:    %-2d %2s", ev.Ch, string(ev.Ch))

	modifier := "none"
	if ev.Mod != 0 {
		modifier = "termbox.ModAlt"
	}
	tboxc.PriToRect(ReStatusLeBo, "Mod: %s", modifier)

	im := termbox.SetInputMode(termbox.InputCurrent)
	tboxc.PriToRect(ReStatusRiBo, "InpMode: %v", im)

}

func printResizeEvent(ev *termbox.Event) {
	tboxc.PriToRect(ReStatusTo, "Resize event: %d x %d", ev.Width, ev.Height)
}

func printMouseEvent(ev *termbox.Event) {
	button := ""
	switch ev.Key {
	case termbox.MouseLeft:
		button = "MouseLeft"
	case termbox.MouseMiddle:
		button = "MouseMiddle"
	case termbox.MouseRight:
		button = "MouseRight"
	}
	cntrMouseEvents++
	tboxc.PriToRect(ReStatusLe, "Mouse event: %d x %d", ev.MouseX, ev.MouseY)
	tboxc.PriToRect(ReStatusRi, "Key: %s, #%v", button, cntrMouseEvents)
}

func crossHair(mx, my int) {

	cb := termbox.CellBuffer()
	dx, dy := termbox.Size()
	_, _ = dx, dy

	for x := RectOutp2.X1; x < dx; x++ {
		for y := RectOutp2.Y1; y < RectOutp2.Y2+1; y++ {
			if x == mx || y == my {
				idx := y*dx + x
				if idx < len(cb)-1 {
					cell := cb[idx]
					if cell.Ch != rune(32) || true {
						cell.Bg = termbox.ColorCyan | termbox.AttrBold
						termbox.SetCell(x, y, cell.Ch,
							termbox.ColorBlack, // cell.Fg
							termbox.ColorBlack|termbox.AttrBold)
					}
				}
			}

		}
	}

}
