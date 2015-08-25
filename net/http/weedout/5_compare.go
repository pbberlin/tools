package weedout

import (
	"strings"

	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/wordb"
	"github.com/pbberlin/tools/util"
)

var opt = levenshtein.Options{1, 1, 1} // cheap substitution

var levelsToProcess = map[int]bool{1: true}
var levelsTolerance = 0

const excerptLen = 20

func rangeOverTexts(mp map[string][]SortEl) []TextifiedTree {

	pf = pfDevNull
	defer func() { pf = pfRestore }()

	frags := []TextifiedTree{}

	for fnKey, atexts := range mp {
		pf("%v\n", fnKey)

		for _, se := range atexts {
			lvl := strings.Count(se.Outl, ".") + 1
			if !levelsToProcess[lvl] {
				continue
			}

			fr := TextifiedTree{fnKey, lvl, se.Outl, se.Text, []Similar{}}
			pf("  cmp %5v lvl%v - len%v   %v \n",
				strings.TrimSpace(fr.Outline), fr.Lvl, len(fr.Text),
				string(fr.Text[:util.Min(len(fr.Text)-1, 3*excerptLen)]))

			rangeOverTexts2(&fr, mp)

			if len(fr.Similars) > 0 {
				frags = append(frags, fr)
			}

		}
	}

	return frags
}

func rangeOverTexts2(src *TextifiedTree, mp map[string][]SortEl) {

	// srcE := word.WrapAsEqualer(string(src.Text), true) // ssrc as Equaler
	srcE := wordb.WrapAsEqualer(src.Text, true)
	srcLen := float64(len(src.Text))

	for fnKey, atexts := range mp {

		if fnKey == src.ArticleUrl {
			pf("    to %v SKIP self\n", fnKey)
			continue
		}

		pf("    to %v\n", fnKey)

		cntr, br := 0, true
		for _, se := range atexts {
			// outl, text := se.Outl, se.Text

			lvl := strings.Count(se.Outl, ".") + 1

			if lvl > src.Lvl+levelsTolerance {
				break // since we are now sorted by lvl, we can this is safe
			}

			if lvl == src.Lvl ||
				(lvl > src.Lvl && lvl <= src.Lvl+levelsTolerance) {
				// proceed
			} else {
				continue
			}

			relSize := srcLen / float64(util.Max(1, len(se.Text)))
			if relSize < 0.33 || relSize > 3 {
				continue
			}

			dstE := wordb.WrapAsEqualer(se.Text, true) // destinations as Equaler
			m := levenshtein.New(srcE, dstE, opt)
			absDist, relDist := m.Distance()

			//
			if relDist < 0.26 && absDist < 10 {
				if br {
					pf("\t")
				}
				sd := string(se.Text[:util.Min(2*excerptLen, len(se.Text)-1)])
				sd = stringspb.ToLen(sd, 2*excerptLen+1)
				_ = sd
				pf("%12v %v %4v %5.2v   ", se.Outl, sd, absDist, relDist)

				cntr++
				br = false

				sim := Similar{}
				sim.ArticleUrl = fnKey
				sim.Lvl = lvl
				sim.Outline = se.Outl
				sim.AbsLevenshtein = absDist
				sim.RelLevenshtein = relDist
				sim.Text = se.Text
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
