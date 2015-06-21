package parse2

import (
	"strings"

	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/wordb"
	"github.com/pbberlin/tools/util"
)

var opt = levenshtein.Options{1, 1, 1}

// var processLevels = map[int]bool{1: true, 2: true}
var processLevels = map[int]bool{1: true}
var lvlTolerance = 1
var frags = []fragment{}

const excerptLen = 20

func rangeOverTexts() {
	for articleId, atexts := range articleTexts {
		pf("%v\n", articleId)
		cntr := 0
		for outl, text := range atexts {
			lvl := strings.Count(outl, ".") + 1
			if !processLevels[lvl] {
				continue
			}
			fr := fragment{articleId, lvl, outl, text, []similarity{}}
			pf("  cmp %5v lvl%v - len%v   %v \n",
				strings.TrimSpace(fr.Outline), fr.Lvl, len(fr.Text),
				string(fr.Text[:util.Min(len(fr.Text), 3*excerptLen)]))

			rangeOverTexts2(&fr)

			if len(fr.Similars) > 0 {
				frags = append(frags, fr)
			}

			cntr++
			if cntr > 20 {
				// pf("  over 20\n")
				// break
			}
		}
	}
}

func rangeOverTexts2(src *fragment) {

	// srcE := word.WrapAsEqualer(string(src.Text), true) // ssrc as Equaler
	srcE := wordb.WrapAsEqualer(src.Text, true)
	srcLen := float64(len(src.Text))

	for articleId, atexts := range articleTexts {
		if articleId == src.ArticleUrl {
			pf("    to %v SKIP self\n", articleId)
			continue
		}
		pf("    to %v\n", articleId)
		cntr, br := 0, true
		for outl, text := range atexts {

			lvl := strings.Count(outl, ".") + 1
			if lvl == src.Lvl ||
				(lvl > src.Lvl && lvl <= src.Lvl+lvlTolerance) {
				// proceed
			} else {
				continue
			}

			relSize := srcLen / float64(util.Max(1, len(text)))
			if relSize < 0.33 || relSize > 3 {
				continue
			}

			dstE := wordb.WrapAsEqualer(text, true) // destinations as Equaler
			m := levenshtein.New(srcE, dstE, opt)
			absDist, relDist := m.Distance()

			//
			if relDist < 0.26 && absDist < 10 {
				if br {
					pf("\t")
				}
				sd := string(text[util.Min(2*excerptLen, len(text))])
				sd = pbstrings.ToLen(sd, 2*excerptLen+1)
				pf("%12v %v %4v %5.2v   ", outl, sd, absDist, relDist)
				cntr++
				br = false

				sim := similarity{}
				sim.ArticleUrl = articleId
				sim.Lvl = lvl
				sim.Outline = outl
				sim.AbsLevenshtein = absDist
				sim.RelLevenshtein = relDist
				sim.Text = text
				src.Similars = append(src.Similars, sim)

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
