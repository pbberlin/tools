// Package screenbuffer is a ultra primitive columnar printing utility;
// abandonded in favor of termbox-go.
package screenbuffer

import "fmt"

type ScreenBuffer [][]byte // NO UNICODE AWARENESS

func (ptrSB *ScreenBuffer) PrintAt(idxLine int, s string) {

	sb := *ptrSB
	latestRow := len(sb) - 1 // init complete

	if latestRow < idxLine { // append_a
		sb = append(sb, make([][]byte, idxLine-latestRow)...)
		*ptrSB = sb             // changing the receiver - NOT ptrSB = &sb
		latestRow = len(sb) - 1 // refresh
	}

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if idxLine == latestRow { // append_b
				sb = append(sb, []byte{})
				*ptrSB = sb             // changing the receiver - NOT ptrSB = &sb
				latestRow = len(sb) - 1 // refresh
			}
			idxLine++
			continue
		}
		sb[idxLine] = append(sb[idxLine], byte(s[i]))
	}

}

func (ptrSB *ScreenBuffer) PrintAppend(s string) {
	sb := *ptrSB
	latestRow := len(sb) - 1
	if latestRow == -1 {
		latestRow = 0
	}
	ptrSB.PrintAt(latestRow, s)
}

func (ptrSB *ScreenBuffer) Dump() {
	sb := *ptrSB
	for i := 0; i < len(sb); i++ {
		// fmt.Printf("%-2d >", i)
		row := string(sb[i])
		fmt.Println(row)
	}
}
