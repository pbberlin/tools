package weedout

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/text/levenshtein/wordb"
	"github.com/pbberlin/tools/util"
)

var opt = levenshtein.Options{1, 1, 1} // cheap substitution

var levelsToProcess = map[int]bool{1: true}

// var levelsToProcess = map[int]bool{1: true, 2: true, 3: true}

var levelsTolerance = 0

const excerptLen = 20

var appliedLevenshtein = 0
var appliedCompare = 0
var breakMapsTooDistinct = 0

func similarTextifiedTrees(mp map[string][]*TextifiedTree, skipPrefix map[string]bool, onlyKeys map[string]bool) []TextifiedTree {

	pf = pfDevNull
	defer func() { pf = pfRestore }()

	frags := []TextifiedTree{}

	for fnKey, tts := range mp {
		if !onlyKeys[fnKey] {
			continue
		}
		pf("%v\n", fnKey)

	MarkX:
		for _, tt := range tts {

			if !levelsToProcess[tt.Lvl] {
				continue
			}

			outls := strings.Split(tt.Outline, ".")
			for i := 0; i < len(outls)-1; i++ {
				jn := strings.Join(outls[0:i+1], ".") + "."
				if skipPrefix[jn] {
					// log.Printf("  %-8v contains %-6v => skip\n", tt.Outline, jn)
					continue MarkX
				} else {
					// log.Printf("  %-8v proccessing ...\n", tt.Outline)
				}

			}

			similarTextifiedTrees2(tt, mp, skipPrefix)
			if len(tt.Similars) > 0 {
				frags = append(frags, *tt)
			}
		}
	}

	return frags
}

func similarTextifiedTrees2(src *TextifiedTree, mp map[string][]*TextifiedTree, skipPrefix map[string]bool) {

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

			if src.NumTokens < 1 {
				continue
			}

			if src.NumTokens < 5 && tt.NumTokens > 7 {
				continue
			}

			if HistoBasedDistance(src, tt) > 0.51 {
				breakMapsTooDistinct++
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
				src.SumAbsLevenshtein += absDist
				src.SumRelLevenshtein += relDist

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

func similaritiesToFile(fs fsi.FileSystem, logdir string, frags []TextifiedTree, stage int) {

	// bfrags := stringspb.IndentedDumpBytes(frags)
	b := new(bytes.Buffer)
	for _, v := range frags {
		b.WriteString(fmt.Sprintf("%v %2v ", v.SourceID, v.Lvl))
		b.WriteString(fmt.Sprintf("%-8v             ", v.Outline))
		b.Write(v.Text)
		b.WriteString("\n")
		for _, v1 := range v.Similars {
			b.WriteString(fmt.Sprintf("%v %2v ", v1.SourceID, v1.Lvl))
			b.WriteString(fmt.Sprintf("%-8v    ", string(v1.Outline)))
			b.WriteString(spf("%2v ", v1.AbsLevenshtein))
			b.WriteString(spf("%-5.2v ", v1.RelLevenshtein))
			b.Write(v1.Text)
			b.WriteByte(10)
		}
		b.WriteByte(10)
	}
	common.WriteFile(fs, spf("%v/outp_fragments_st%v.txt", logdir, stage), b.Bytes())

}

// HistoBasedDistance isa cheap alternative to Levenshtein.distance.
// Particularly, since we compute the histogram anyway.
// Tests show, that LevenshteinDistance > HistoBasedDistance
// i.e.             0.36                  0.31
// i.e.             0.6                   0.4
// Thus we can break early if i.e. HistoBasedDistance > 0.5
// implying that LevenshteinDistance is at least >= 0.5
func HistoBasedDistance(src, dst *TextifiedTree) float64 {

	largerOuter := src
	inner := dst
	if src.NumTokens < dst.NumTokens {
		largerOuter = dst
		inner = src
	}

	// Handle division by zero
	if largerOuter.NumTokens == 0 {
		return 0.0
	}

	// inner overlap
	same := 0
	for k, _ := range largerOuter.Histo {
		if _, ok := inner.Histo[k]; ok {
			same++
		}
	}

	distinctBySheerSize := largerOuter.NumTokens - inner.NumTokens
	distinctInner := inner.NumTokens - same

	ret := float64(distinctBySheerSize+distinctInner) / float64(largerOuter.NumTokens)

	// crit1 := inner.NumTokens > 5 && distinctBySheerSize < 5 && distinctInner < 5
	// _ = crit1
	// if ret > 0 && ret < 0.51 {
	// 	fmt.Printf("%3v %3v ; sizediff %3v worddiff %3v =>  %4.2v\n", largerOuter.NumTokens, inner.NumTokens,
	// 		distinctBySheerSize, distinctInner, ret)
	// }

	return ret
}

// slow - only for debug
func mpKeys(mp map[string]int) string {
	ret := ""
	for k, _ := range mp {
		ret += k + " "
	}
	return ret
}
