package tboxc

// Terminal Box Comfort

import "github.com/pbberlin/termbox-go"
import "fmt"

type Rect struct {
	X1, Y1, X2, Y2 int
}

func PriRuneAt(x, y int, r rune) {
	termbox.SetCell(x, y, r, fgc, brushBG)
}

func PriLineAt(x, y int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	for _, rne := range s {
		termbox.SetCell(x, y, rne, fgc, brushBG)
		x++
	}
}

func PriColoredLineAt(x, y int, fg termbox.Attribute, format string,
	args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	for _, rne := range s {
		termbox.SetCell(x, y, rne, fg, brushBG)
		x++
	}
}

// ultimately: throw away direct PriParAt
func PriToRect(r Rect, format string, args ...interface{}) {
	PriParAt(r.X1, r.Y1, r.X2-r.X1, r.Y2-r.Y1, fgc, format, args...)
}

// priParAt - printParagraphAt prints a formatted string at given coords.
// It starts a new line, each time param width is surpassed.
// Forthcoming lines are indented.
// priParAt() blanks out previous printings to the extent of width x height.
// Previous prints, that were overflowing height remain.
func PriParAt(x, y, width, height int, fg termbox.Attribute, format string,
	args ...interface{}) {

	for i := x; i <= x+width; i++ {
		for j := y; j <= y+height; j++ {
			// termbox.SetCell(i, j, ' ', fg, termbox.ColorGreen)
			termbox.SetCell(i, j, ' ', fg, brushBG)
		}
	}

	xInit := x + 1
	indent := 1
	s := fmt.Sprintf(format, args...)
	for _, rne := range s {
		if rne != '\n' {
			termbox.SetCell(x, y, rne, fg, brushBG)
		}
		x++
		if x > xInit+width-1 || rne == '\n' {
			y++
			x = xInit + indent - 1
		}
	}

}

// i.e. demoPrintCombinedEdges(x1+1, 2)
func demoPrintCombinedEdges(x, y int) {
	PriRuneAt(x+0, y, 0x251C+0*8) // vert, branch right
	PriRuneAt(x+1, y, 0x251C+1*8) // vert, branch left
	PriRuneAt(x+2, y, 0x251C+2*8) // horz, branch down
	PriRuneAt(x+3, y, 0x251C+3*8) // horz, branch top
}

func PainHorzLine(x, y1, y2 int) {
	for i := y1; i <= y2; i++ {
		PriRuneAt(i, x, 0x2500)
	}
}

func PainFatVertLine(y, x1, x2 int) {
	for i := x1; i <= x2; i++ {
		PriRuneAt(y, i, 0x2588)
	}
}

// blank with current brush
func BlankOut(r Rect) {
	for i := r.X1; i <= r.X2; i++ {
		for j := r.Y1; j <= r.Y2; j++ {
			PriRuneAt(i, j, rune(32))
		}
	}
}

// painRect blanks out the area.
// painRect returns the *inner* rect coords
func PainRect(x1, y1, x2, y2, attachedTo int) (xi1, yi1, xi2, yi2 int) {

	BlankOut(Rect{x1, y1, x2, y2})

	// edges
	PriRuneAt(x1, y1, 0x250C+0*4) // left top
	PriRuneAt(x2, y1, 0x250C+1*4) // right top
	PriRuneAt(x1, y2, 0x250C+2*4) // bottom left
	PriRuneAt(x2, y2, 0x250C+3*4) // br

	switch attachedTo {
	case 1: // attach to left
		PriRuneAt(x1, y1, 0x251C+2*8)
		PriRuneAt(x1, y2, 0x251C+3*8)
	case 2: // attach to top
		PriRuneAt(x1, y1, 0x251C+0*8)
		PriRuneAt(x2, y1, 0x251C+1*8)
	}

	// horizontal lines (inner, because edges are set)
	for i := x1 + 1; i < x2; i++ {
		PriRuneAt(i, y1, 0x2500)
		PriRuneAt(i, y2, 0x2500)
	}

	// vertical lines (inner, because edges are set)
	for i := y1 + 1; i < y2; i++ {
		PriRuneAt(x1, i, 0x2502)
		PriRuneAt(x2, i, 0x2502)
	}

	return x1 + 1, y1 + 1, x2 - 1, y2 - 1

}
