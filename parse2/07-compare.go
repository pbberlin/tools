package parse2

import (
	"strings"

	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/wordb"
	"github.com/pbberlin/tools/util"
)

var opt = levenshtein.Options{1, 1, 1}

const cMaxLvl = 1

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
			lvl := strings.Count(outl, ".") + 1
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

	pf("  cmp %5v lvl%v - len%v   %v \n",
		strings.TrimSpace(src.Outline), src.Lvl, len(src.Text),
		string(src.Text[:util.Min(len(src.Text), 40)]))

	for articleId, atexts := range articleTexts {
		if articleId == src.ArticleUrl {
			pf("    to %v SKIP self\n", articleId)
			continue
		}
		pf("    to %v\n", articleId)
		cntr, br := 0, true
		for outl, text := range atexts {

			lvl := strings.Count(outl, ".") + 1
			// if lvl > cMaxLvl {
			if lvl != src.Lvl {
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

			pf("%12v %v %4v %5.2v   ", outl, sd, absDist, relDist)

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
