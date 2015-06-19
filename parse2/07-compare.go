package parse2

import (
	"strings"

	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/word"
)

var opt = levenshtein.Options{1, 1, 1}

const cMaxLvl = 2

func rangeOverTexts() {
	for articleId, atexts := range articleTexts {
		pf("%v\n", articleId)
		cntr := 0
		for outl, text := range atexts {
			lvl := strings.Count(outl, ".")
			if lvl > cMaxLvl {
				continue
			}
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

	src := string(srcText)
	srcE := word.WrapAsEqualer(src, true) // src as Equaler

	pf(" cmp  l%v - o%v - len%v  %v  \n", srcLvl, strings.TrimSpace(srcOutl), len(src), pbstrings.Ellipsoider(src, 10))

	for articleId, atexts := range articleTexts {
		if articleId == srcArticle {
			pf("    to %v SKIP\n", articleId)
			continue
		}
		pf("    to %v\n", articleId)
		cntr, br := 0, true
		for outl, text := range atexts {

			lvl := strings.Count(outl, ".")
			if lvl > cMaxLvl {
				continue
			}

			if br {
				pf("\t")
			}
			s := string(text)

			dstE := word.WrapAsEqualer(s, true) // destinations as Equaler
			m := levenshtein.New(srcE, dstE, levenshtein.DefaultOptions)
			absDist, relDist := m.Distance()

			sd := pbstrings.Ellipsoider(s, 10)
			sd = pbstrings.ToLen(sd, 21)

			pf("%v %v %2v %5.2v   ", pbstrings.ToLen(outl, 11), sd, absDist, relDist)

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
