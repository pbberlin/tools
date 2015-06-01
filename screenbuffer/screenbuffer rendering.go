package screenbuffer

import "strings"

var suffx1 = "    "

func (m *TransposableMatrix) RenderMatrixFrom(perspective int, sb *ScreenBuffer,
	x1, x2, y1, y2 int) {
	m.renderHeader(sb, x1, x2)
	sbLine := 1
	for y := y1; y < y2; y++ {
		sbLine++
		sb.PrintAt(sbLine, spf("%3d|", y))
		for x := x1; x < x2; x++ {
			if x == 0 && y == 0 {
				sb.PrintAt(sbLine, "++")
				continue
			}
			sl := m.Get(x, y, perspective)
			if sl.Label != "" {
				sb.PrintAt(sbLine, sl.Label)
			} else {
				sb.PrintAt(sbLine, "  ")
			}
		}
		sb.PrintAt(sbLine, suffx1)
	}
}

func (m *TransposableMatrix) renderHeader(sb *ScreenBuffer, x1, x2 int) {

	// row1
	sb.PrintAt(0, "  x|")
	for x := x1; x < x2; x++ {
		colHead := spf("%2d", x%10)
		sb.PrintAt(0, colHead)
	}
	sb.PrintAt(0, suffx1)

	// row2
	sb.PrintAt(1, "  y|")
	sb.PrintAt(1, spf(strings.Repeat("=", 2*(x2-x1))))
	sb.PrintAt(1, suffx1)
}
