package weedout

import (
	"bytes"

	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/wordb"
	"github.com/pbberlin/tools/util"
)

var opt = levenshtein.Options{1, 1, 1} // cheap substitution

var levelsToProcess = map[int]bool{1: true}
var levelsTolerance = 0

const excerptLen = 20

var appliedLevenshtein = 0
var appliedCompare = 0

func similarTextifiedTrees(mp map[string][]*TextifiedTree, onlyFirst bool) []TextifiedTree {

	pf = pfDevNull
	defer func() { pf = pfRestore }()

	frags := []TextifiedTree{}

	for fnKey, tts := range mp {
		pf("%v\n", fnKey)

		for _, tt := range tts {
			if !levelsToProcess[tt.Lvl] {
				continue
			}
			similarTextifiedTrees2(tt, mp)
			if len(tt.Similars) > 0 {
				frags = append(frags, *tt)
			}

		}
		if onlyFirst {
			break
		}
	}

	return frags
}

func similarTextifiedTrees2(src *TextifiedTree, mp map[string][]*TextifiedTree) {

	// srcE := word.WrapAsEqualer(string(src.Text), true) // ssrc as Equaler
	srcE := wordb.WrapAsEqualer(src.Text, true)
	srcLen := float64(len(src.Text))

	for fnKey, tts := range mp {

		if fnKey == src.SourceID {
			pf("    to %v SKIP self\n", fnKey)
			continue
		}

		pf("    to %v\n", fnKey)

		cntr, br := 0, true
		for _, tt := range tts {
			// outl, text := tt.Outl, tt.Text

			if tt.Lvl > src.Lvl+levelsTolerance {
				break // since we are now sorted by lvl, we can this is safe
			}

			if tt.Lvl == src.Lvl ||
				(tt.Lvl > src.Lvl && tt.Lvl <= src.Lvl+levelsTolerance) {
				// proceed
			} else {
				continue
			}

			if src.NumTokens < 2 {
				continue
			}

			if src.NumTokens < 5 && tt.NumTokens > 7 {
				continue
			}

			relSize := srcLen / float64(util.Max(1, len(tt.Text)))
			if relSize < 0.33 || relSize > 3 {
				continue
			}

			absDist, relDist := 0, 0.0

			if tt.NumTokens == src.NumTokens &&
				len(tt.Text) == len(src.Text) &&
				bytes.Equal(tt.Text, src.Text) {
				absDist, relDist = 0, 0.0
				appliedCompare++
			} else {
				dstE := wordb.WrapAsEqualer(tt.Text, true) // destinations as Equaler
				m := levenshtein.New(srcE, dstE, opt)
				absDist, relDist = m.Distance()
				appliedLevenshtein++
			}

			// if relDist < 0.4 && relDist > 0.0 {
			// 	s1 := string(src.Text)
			// 	s2 := string(tt.Text)
			// 	fmt.Printf("%v %14v %4v %5.2v %s %s\n", src.SourceID, tt.Outline, absDist, relDist,
			// 		stringspb.ToLen(s1, 34),
			// 		stringspb.ToLen(s2, 34),
			// 	)
			// }

			//
			if relDist < 0.26 && absDist < 10 {
				if br {
					pf("\t")
				}

				sd := ""
				sd = string(tt.Text[:util.Min(2*excerptLen, len(tt.Text)-1)])
				sd = stringspb.ToLen(sd, 2*excerptLen+1)
				pf("%12v %v %4v %5.2v   ", tt.Outline, sd, absDist, relDist)

				cntr++
				br = false

				sim := Similar{}
				sim.SourceID = fnKey
				sim.Lvl = tt.Lvl
				sim.Outline = tt.Outline
				sim.AbsLevenshtein = absDist
				sim.RelLevenshtein = relDist
				sim.Text = tt.Text
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
