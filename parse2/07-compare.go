package parse2

import (
	"strings"

	"github.com/pbberlin/tools/pbstrings"
)

func rangeOverTexts() {
	for articleId, atexts := range articleTexts {
		pf("%v\n", articleId)
		cntr := 0
		for outl, text := range atexts {
			lvl := strings.Count(outl, ".") + 1
			rangeOverTexts2(articleId, lvl, outl, text)
			cntr++
			if cntr > 20 {
				pf("  over 20\n")
				break
			}
		}
	}
}

func rangeOverTexts2(srcArticle string, srcLvl int, srcOutl string, srcText []byte) {

	src := pbstrings.Ellipsoider(string(srcText), 10)

	pf(" cmp  %v - l%v - %v - %v  \n", srcArticle, srcLvl, srcOutl, src)

	for articleId, atexts := range articleTexts {
		if articleId == srcArticle {
			pf("    to %v SKIP\n", articleId)
			continue
		}
		pf("    to %v\n", articleId)
		cntr, br := 0, true
		for outl, text := range atexts {
			if br {
				pf("\t")
			}
			s := string(text)
			s = pbstrings.Ellipsoider(s, 10)
			s = pbstrings.ToLen(s, 21)
			pf("%v %v    ", pbstrings.ToLen(outl, 11), s)
			cntr++
			br = false
			if cntr%3 == 0 || cntr > 20 {
				pf("\n")
				br = true
			}
			if cntr > 20 {
				break
			}
		}
		if !br {
			pf("\n")
		}
	}

}
