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
const excerptLen = 20

type fragment struct {
	ArticleUrl string
	Lvl        int
	Outline    string
	Text       []byte
	Similars   [][]byte
}

var frags = []fragment{}

func rangeOverTexts() {
	for articleId, atexts := range articleTexts {
		pf("%v\n", articleId)
		cntr := 0
		for outl, text := range atexts {
			lvl := strings.Count(outl, ".") + 1
			if lvl > cMaxLvl {
				continue
			}
			fr := fragment{articleId, lvl, outl, text, [][]byte{}}
			rangeOverTexts2(&fr)

			if len(fr.Similars) > 0 {
				frags = append(frags, fr)
			}

			cntr++
			if cntr > 20 {
				pf("  over 20\n")
				break
			}
		}
	}
}

func rangeOverTexts2(src *fragment) {

	// srcE := word.WrapAsEqualer(string(src.Text), true) // ssrc as Equaler
	srcE := wordb.WrapAsEqualer(src.Text, true)

	pf("  cmp %5v lvl%v - len%v   %v \n",
		strings.TrimSpace(src.Outline), src.Lvl, len(src.Text),
		string(src.Text[:util.Min(len(src.Text), 3*excerptLen)]))

	for articleId, atexts := range articleTexts {
		if articleId == src.ArticleUrl {
			pf("    to %v SKIP self\n", articleId)
			continue
		}
		pf("    to %v\n", articleId)
		cntr, br := 0, true
		for outl, text := range atexts {

			lvl := strings.Count(outl, ".") + 1
			if lvl != src.Lvl && lvl != src.Lvl+1 && lvl != src.Lvl-1 {
				continue
			}

			relSize := float64(len(src.Text)) / float64(util.Max(1, len(text)))
			if relSize < 0.33 || relSize > 3 {
				// pf("skipping size ratio %4.2v  ", relSize)
				continue
			}

			dstE := wordb.WrapAsEqualer(text, true) // destinations as Equaler
			m := levenshtein.New(srcE, dstE, opt)
			absDist, relDist := m.Distance()

			//
			if relDist < 0.75 {
				if br {
					pf("\t")
				}
				sd := pbstrings.Ellipsoider(string(text), excerptLen)
				sd = pbstrings.ToLen(sd, 2*excerptLen+1)
				pf("%12v %v %4v %5.2v   ", outl, sd, absDist, relDist)
				cntr++
				br = false

				src.Similars = append(src.Similars, text)

				if cntr%2 == 0 || cntr > 20 {
					pf("\n")
					br = true
				}
				if cntr > 20 {
					break
				}
			}

		}
		if !br {
			pf("\n")
		}
	}

}
