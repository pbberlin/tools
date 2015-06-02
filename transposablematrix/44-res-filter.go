package transposablematrix

import "github.com/pbberlin/tools/util"

type Filterer interface {
	Filter([][]Amorph, Fusion) *Amorph
}

type MostAbundant struct{}
type MostAbundantInProximity struct{}

// May be changed from outside to any type transposablematrix.Filterer
var ActiveFilter Filterer = MostAbundantInProximity{}

// var ActiveFilter Filterer = MostAbundant{}

// mostAbundant seeks the largest slice.
// It then returns the first amorph.
// Should be replaced by abundantHeightMatch()
func (dummy MostAbundant) Filter(amorphBlocks [][]Amorph, fs Fusion) (chosen *Amorph) {

	maxFound := 0
	for i := 0; i < len(amorphBlocks); i++ {
		maxFound = util.Max(maxFound, len(amorphBlocks[i]))
	}

	for i := 0; i < len(amorphBlocks); i++ {
		if len(amorphBlocks[i]) == maxFound {
			if len(amorphBlocks[i]) > 0 {
				return &amorphBlocks[i][0]
			}
		}
	}
	return nil
}

// abundantHeightMatch should replace mostAbundant()
// for the selection of the most appropriate amorph.
// It returns at least an amorph, complying to max height.
// If there are *several* amorphs close to the optimal height,
// then we return one of the most abundant in the interval plus-minus 2
func (dummy MostAbundantInProximity) Filter(amorphBlocks [][]Amorph, fs Fusion) (chosen *Amorph) {

	pfTmp := intermedPf(pf)
	defer func() { pf = pfTmp }()
	pf = pfDevNull

	heightLim := fs.pm[2]
	heightOpt := fs.FillHeightFloor()

	pf("lim%v,opt%v ", heightLim, heightOpt)

	// plus minus 2
	// -2=>0  , -1=>1  , 0=>2  , 1=>3 , 2=>4
	closest := [][]Amorph{
		[]Amorph{}, []Amorph{}, []Amorph{}, []Amorph{}, []Amorph{},
	}

	for i := 0; i < len(amorphBlocks); i++ {

		amorphs := amorphBlocks[i]

		// Scan all blocks of amorphs.
		// Extract one, which is closest to desired height.
		// Also create a histogram of amorphs,
		// which have a height of plus-minus 2
		for j := 0; j < len(amorphs); j++ {

			var lastDist, lpDist int

			lp := amorphs[j]

			if lp.Rows > heightLim {
				continue
			}

			if chosen == nil {
				chosen = &lp
			}
			lastDist = chosen.Rows - heightOpt
			lpDist = lp.Rows - heightOpt
			if util.Abs(lpDist) < util.Abs(lastDist) {
				chosen = &lp
			}

			if util.Abs(lpDist) <= 2 {
				closest[lpDist+2] = append(closest[lpDist+2], lp)
			}

		}
	}

	// Debug output
	if false {
		pf("\n")
		for i := 0; i < len(closest); i++ {
			sa := closest[i]
			pf("h%v fnd%v:  ", heightOpt+i-2, len(sa))
			for j := 0; j < len(sa); j++ {
				pf("%2v %vx%v=%2v | ", sa[j].IdxArticle, sa[j].Rows, sa[j].Cols, sa[j].NElements)
			}
			pf("\n")
		}
	}

	// How many amorphs were close to desired height
	maxBatch := 0
	for i := 0; i < len(closest); i++ {
		if len(closest[i]) > maxBatch {
			maxBatch = len(closest[i])
		}
	}

	// Return one abundant amorph
	// from the range +/- 2
	if maxBatch > 0 {
		for {
			if len(closest[2]) == maxBatch {
				chosen = &closest[2][0]
				break
			}
			if len(closest[1]) == maxBatch {
				chosen = &closest[1][0]
				break
			}
			if len(closest[3]) == maxBatch {
				chosen = &closest[3][0]
				break
			}
			if len(closest[0]) == maxBatch {
				chosen = &closest[0][0]
				break
			}
			if len(closest[4]) == maxBatch {
				chosen = &closest[4][0]
				break
			}
			break
		}
	}

	return
}
