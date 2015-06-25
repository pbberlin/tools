// Package tboxc contains comfort functions for go terminal box (termbox-go).
package tboxc

import "github.com/pbberlin/termbox-go"

var fallbackBG = termbox.ColorBlue // "default"
var brushBG = termbox.ColorBlue
var fgc = termbox.ColorWhite | termbox.AttrBold

var excludeColorIdxs = map[int]bool{
	0:               true, // termbox.ColorDefault
	1:               true, // black
	int(fallbackBG): true, // 5, blue, 5+8 / 5+512 is allowed (lightblue)
	16:              true, // white, equiv. to 8+512 (grey+bold)
}

const minColIdx = 2
const maxColIdx = 16
const brightThreshold = 8

var revolving = minColIdx

var ArticleIDsToForegrounds = map[int]termbox.Attribute{}

func RevolveValidIdx() int {
	revolving++
	if revolving > maxColIdx {
		revolving = minColIdx
	}
	for {
		if excludeColorIdxs[revolving] {
			revolving++
		} else {
			break
		}
	}
	return revolving
}

//
func ColorByInt(idx int, exclude bool) termbox.Attribute {
	var newCol termbox.Attribute

	cntr := 0
	if exclude {
		for {
			remainder := idx % (maxColIdx + 1)
			if excludeColorIdxs[remainder] {
				increment := 1         // leads to accumulation of the first non-excluded color
				increment = idx%11 + 1 // more variation, still deterministic
				idx += increment
			} else {
				break
			}
			cntr++
			if cntr > 100 {
				panic("colors forever...")
			}
		}
	}

	idx1 := idx % (maxColIdx + 1)
	idx2 := uint16(idx1)
	// if idx2 >= brightThreshold{
	// 	fmt.Printf("c%v++|", idx2)
	// }

	switch {
	case idx2 <= brightThreshold:
		newCol = termbox.Attribute(idx2)
	default:
		newCol = termbox.Attribute(idx2 - brightThreshold)
		newCol = newCol | termbox.AttrBold
	}
	return newCol
}

func BrushReset() {
	brushBG = fallbackBG
	fgc = termbox.ColorWhite | termbox.AttrBold
}

func BrushByArticleID(ArticleId int) (termbox.Attribute, termbox.Attribute) {
	brushBG = ColorByInt(ArticleId, true)
	if establFGC, ok := ArticleIDsToForegrounds[ArticleId]; ok {
		fgc = establFGC
	} else {
		fgc = shiftingComplement(brushBG)
		ArticleIDsToForegrounds[ArticleId] = fgc
	}
	fgc = termbox.ColorBlack // all black
	return fgc, brushBG
}

// Since we dont have enough background colors to differentiate
// we also vary the foreground color.
func shiftingComplement(bgc termbox.Attribute) (compl termbox.Attribute) {
	idx := RevolveValidIdx()
	if idx == int(bgc) || idx == int(bgc)-int(termbox.AttrBold) {
		revolving++ // neccessary to jump over ranges of several exclusions?
		idx = RevolveValidIdx()
	}
	compl = ColorByInt(idx, false)
	return
}
