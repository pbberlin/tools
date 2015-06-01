package transposablematrix

import "github.com/pbberlin/tools/util"

// mostAbundant seeks the largest slice.
// It then returns the first amorph.
// Should be replaced by abundantHeightMatch()
func mostAbundant(amorphBlocks [][]Amorph) (chosen *Amorph) {

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
func abundantHeightMatch(amorphBlocks [][]Amorph, fs Fusion) (chosen *Amorph) {

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

// exactStraightEdge returns amorphs with a requested width.
// It singles out the most abundant of heights for
// the return of "chosen"
func exactStraightEdge(ar *Reservoir, x1 int) ([][]Amorph, *Amorph) {

	pfTmp := intermedPf(pf)
	defer func() { pf = pfTmp }()

	minHeight := ar.SmallestDesirableHeight - 1

	fnd := make([][]Amorph, cMaxHeight)

	for i := 0; i < cMaxHeight; i++ {
		enc := Enc(x1, minHeight+i, 0)
		if mp, ok := ar.EdgesSlackless[enc]; ok {
			fnd[i] = []Amorph{}
			for k, _ := range mp {
				lp := ar.Amorphs[k] // effects a copying of the amorph
				lp.Cols = x1
				lp.Rows, lp.Slack = OtherSide(lp.NElements, lp.Cols)
				lp.Edge = nil // reset
				fnd[i] = append(fnd[i], lp)
			}
			pf("fnd%v %vx%v ", len(fnd[i]), x1, minHeight+i)
		}
	}

	chosen := mostAbundant(fnd)

	return fnd, chosen
}

//
// exactStairyEdge returns amorphs
// with exactly the desired edges.
// The edge is also attached to the amorph.
// Param limit restricts the amount of amorphs returned.
func (ar *Reservoir) exactStairyEdge(x1, y, x2 int, limit int) (amorphs []Amorph) {

	enc := Enc(x1, y, x2)
	if _, ok := ar.Edge3[enc]; ok {
		mp := ar.Edge3[enc]

		for k, _ := range mp {

			// pf("found %v \n", k)
			lp := ar.Amorphs[k] // effects a copying of the amorph
			lp.Cols = x1 + x2
			lp.Rows, lp.Slack = OtherSide(lp.NElements, lp.Cols)

			// Increase rows,
			// if the requested edge is a superedge.
			cntr := 0
			for lp.Rows <= util.Abs(y) || lp.Slack < util.Abs(x2*y) {
				lp.Rows++
				lp.Slack = (lp.Cols * lp.Rows) - lp.NElements
				if cntr++; cntr > 100 {
					panic("superedge blowup logic faulty")
				}
			}

			// Attach the edge
			lp.Edge = []int{x1, y, x2}

			amorphs = append(amorphs, lp)

			if len(amorphs) >= limit {
				return
			}

		}
	}
	return
}

// moreThanXElements returns an amorph closest above
// requested number of elements.
// There is no effort to find the most abundant amorph
// or the most height-appropriate among several.
// Its for the most desperate heuristic anyway.
func moreThanXElements(ar *Reservoir, minElements int) (chosen *Amorph) {

	pfTmp := intermedPf(pf)
	defer func() { pf = pfTmp }()

lblFBNE:
	for i := 0; i < cMaxDiff; i++ {
		if i < 2 || i%5 == 0 {
			pf(" srch%v", minElements+i)
		}
		pf("+%v", i)
		if mp, ok := ar.MElements[minElements+i]; ok {
			for k, _ := range mp {
				chosen = &ar.Amorphs[k]
				pf(" fnd%v ID%v\n", chosen.NElements, chosen.IdxArticle)
				break lblFBNE // label points to begin of loop - *but* really terminates it
			}
		}
	}
	return
}
