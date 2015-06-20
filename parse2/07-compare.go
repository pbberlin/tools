package parse2

import (
	"strings"

	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/wordb"
)

var opt = levenshtein.Options{1, 1, 1}

const cMaxLvl = 2

type fragment struct {
	ArticleUrl string
	Lvl        int
	Outline    string
	Text       []byte
}

func rangeOverTexts() {
	for articleId, atexts := range articleTexts {
		pf("%v\n", articleId)
		cntr := 0
		for outl, text := range atexts {
			lvl := strings.Count(outl, ".")
			if lvl > cMaxLvl {
				continue
			}
			rangeOverTexts2(fragment{articleId, lvl, outl, text})
			cntr++
			if cntr > 20 {
				pf("  over 20\n")
				break
			}
		}
	}
}

func rangeOverTexts2(src fragment) {

	// srcE := word.WrapAsEqualer(string(src.Text), true) // ssrc as Equaler
	srcE := wordb.WrapAsEqualer(src.Text, true)

	pf(" cmp  l%v - o%v - len%v    \n", src.Lvl, strings.TrimSpace(src.Outline), len(src.Text))

	for articleId, atexts := range articleTexts {
		if articleId == src.ArticleUrl {
			pf("    to %v SKIP self\n", articleId)
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

			dstE := wordb.WrapAsEqualer(text, true) // destinations as Equaler
			m := levenshtein.New(srcE, dstE, levenshtein.DefaultOptions)
			absDist, relDist := m.Distance()

			sd := pbstrings.Ellipsoider(string(text), 10)
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
