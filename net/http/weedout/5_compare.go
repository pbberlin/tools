package weedout

import (
	"bytes"
	"sort"
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
			text := cleanseTextForComparisonOnly(se.Text)
			fr := TextifiedTree{fnKey, lvl, se.Outl, text, []Similar{}}
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
			outl, text := se.Outl, se.Text

			lvl := strings.Count(outl, ".") + 1

			if lvl > src.Lvl+levelsTolerance {
				break // since we are now sorted by lvl, we can this is safe
			}

			if lvl == src.Lvl ||
				(lvl > src.Lvl && lvl <= src.Lvl+levelsTolerance) {
				// proceed
			} else {
				continue
			}

			text = cleanseTextForComparisonOnly(text)
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
				sd := string(text[:util.Min(2*excerptLen, len(text)-1)])
				sd = stringspb.ToLen(sd, 2*excerptLen+1)
				_ = sd
				pf("%12v %v %4v %5.2v   ", outl, sd, absDist, relDist)

				cntr++
				br = false

				sim := Similar{}
				sim.ArticleUrl = fnKey
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

func cleanseTextForComparisonOnly(text []byte) []byte {
	// text = bytes.Replace(text, []byte(" hbr"), []byte{}, -1)
	// text = bytes.Replace(text, []byte(" sbr"), []byte{}, -1)
	// text = bytes.Replace(text, []byte(`[img] `), []byte{}, -1)
	// text = bytes.Replace(text, []byte(`[a] `), []byte{}, -1)

	text = bytes.Replace(text, []byte{46}, []byte{32}, -1) // dot
	text = bytes.Replace(text, []byte{44}, []byte{}, -1)   // comma
	text = bytes.Replace(text, []byte{45}, []byte{32}, -1) // hyphen
	text = bytes.Replace(text, []byte{47}, []byte{32}, -1) // forward slash

	words := bytes.Split(text, []byte{byte(32)})

	mp := map[string]int{}
	for _, word := range words {
		mp[string(word)]++
	}

	keys := make([]string, 0, len(mp))
	for k, _ := range mp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := []byte{32}
	for _, key := range keys {
		if len(key) > 1 {
			buf = append(buf, []byte(key)...)
			buf = append(buf, byte(32))
		}
	}

	return buf
}
