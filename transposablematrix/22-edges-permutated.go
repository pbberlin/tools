package transposablematrix

import (
	"sort"
	"time"
)

// RemoveSpentEdges() removes all shapeshifts from the maps,
// so that the maps are kept up to date
// and show a true picture of abundance
// i.e.
// for i := 0; i < len(AR.Amorphs)-12; i++ {
// 	AR.RemoveSpentEdges(AR.Amorphs[i])
// }
func (ar *Reservoir) RemoveSpentEdges(am Amorph) {
	ch1 := make(chan bool)
	// var chTimeOut <-chan time.Time // the <- is required - but I do not understand it at all
	chTimeOut := time.After(2 * time.Second)
	go func() {
		for {
			select {
			case <-ch1:
				// pf("received ch1\n")
				return
			case <-chTimeOut:
				pf("RemoveSpentEdges timed out\n")
				pf("%+v\n", am)
				// default:
			}
		}

	}()
	// pf("%+v\n", am)

	SubmapAddRemove(true, ar.MElements, am.NElements-am.Padded, am.IdxArticle)
	ar.PermutateShapes(am, true)
	ch1 <- true
}

// PermutateShapes shapeshifts an amorph in all conceivable ways.
// Only restriction: Maximum *one* stair
// Results are stored into various maps for quick retrieval
func (ar *Reservoir) PermutateShapes(amArg Amorph, doRemoval bool) {

	plugs := [][]int{} // all possible sections of edges

	nElements := amArg.NElements - amArg.Padded

	// try different widths and heights
	//
	// first compute starting length of edge
	// starting not with 1, but ceil(root)
	rootCeil := RootCeiling(nElements)
	cols, rows, slack := 0, 0, 0

	//
	// now try all side lengths.
	// Look at output of PrintPermutations() to understand
	// why we run from rootCeil...nElements, instead of 2...rootCeil
	for cols = rootCeil; cols <= nElements; cols++ {
		rows, slack = OtherSide(nElements, cols)
		if slack == 0 {
			SubmapAddRemove(doRemoval, ar.EdgesSlackless, Enc(cols, rows, 0), amArg.IdxArticle)
			SubmapAddRemove(doRemoval, ar.EdgesSlackless, Enc(rows, cols, 0), amArg.IdxArticle)
		}
		plugs = PermutateEdges(cols, rows, slack, plugs)
	} // next shapeshift

	//
	// Superedges
	// In addition we also want to compute amorphs *larger* than x*y<=nElements
	// Uneven nElements are distributable at most root(n/2)+1 - root(n/2)+1
	//		5 els => 3*3
	// Even   nElements are distributable at most root(n/2)+1 - root(n/2)
	//		6 els => 4*3
	superEdge1 := nElements/2 + 1
	superEdge2 := superEdge1
	if nElements%2 == 0 {
		superEdge2 = nElements / 2
	}
	for cols := rootCeil; cols <= superEdge2; cols++ {
		for rows := cols; rows <= superEdge1; rows++ { // note rows == COLS prevents 4*3 after 3*4
			slack := (cols * rows) - nElements
			plugs = PermutateEdges(cols, rows, slack, plugs)
		}
	}

	//
	//
	// Now attach the collected kerfs to the global structs
	for i := 0; i < len(plugs); i++ {
		k := plugs[i]
		SubmapAddRemove(doRemoval, ar.Edge3, Enc(k[0], k[1], k[2]), amArg.IdxArticle) // k[2], k[1], k[0] would be WRONG
		enc := Enc(Granularize(k[0]), Granularize(k[1]), Granularize(k[2]))
		SubmapAddRemove(doRemoval, ar.EdgeCond, enc, amArg.IdxArticle)
	}

}

//
// Factorize determines all prime factors of parameter prod (for product)
// We stupidly try each number.
// Since we quickly break the loop, the function is still - barely - acceptable.
// Todo: We could/should build a fixed-limited or a synchronized global lookup table.
func Factorize(prod int) []int {
	const maxFactor = 100
	if prod > maxFactor*maxFactor {
		panic("Factorize() is only intended for smaller arguments")
	}
	ret := make([]int, 0, 2)
	for i1 := 2; i1 < maxFactor; i1++ {
		if i1*i1 > prod {
			// pf("broken after %v;  ", i)
			break
		}
		i2 := prod / i1
		if i2*i1 == prod { // equivalent prod%i1 == 0
			if i2 == i1 {
				ret = append(ret, i1)
			} else {
				ret = append(ret, i1, i2)
			}
		}
	}
	if len(ret) > 2 {
		sort.Ints(ret)
	}
	return ret
}

func Enc(i1, i2, i3 int) int {
	if i2 < 0 {
		i1, i3 = i3, i1
		i2 = -i2
	}
	enc := i1*multiSortFactorSq + i2*multiSortFactor + i3
	return enc
}
func Dec(enc int) (i1, i2, i3 int, s string) {
	i1 = enc / multiSortFactorSq
	enc -= i1 * multiSortFactorSq
	i2 = enc / multiSortFactor
	enc -= i2 * multiSortFactor
	i3 = enc
	s = spf("%v-%v-%v", i1, i2, i3)
	return
}

func SubmapAddRemove(doRemoval bool, mp map[int]map[int]bool, k1 int, k2 int) {
	_, ok := mp[k1]
	if doRemoval {
		if ok {
			delete(mp[k1], k2)
			// pf("del%v %v|", k1, k2)
		}
		if len(mp[k1]) == 0 {
			delete(mp, k1)
		}
	} else {
		if !ok {
			mp[k1] = map[int]bool{}
		}
		mp[k1][k2] = true
	}
}

func PrintEdgesRestrictTo(label string, argMp map[int]map[int]bool, restrictTo map[int]bool) {
	keys := make([]int, 0, len(argMp))
	for k, _ := range argMp {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	cntr := 0
	newLineTerm := true
	numCols := 8
	for _, k := range keys {
		if ok, _ := restrictTo[k]; !ok && len(restrictTo) > 0 {
			continue
		}

		if newLineTerm {
			pf("%s", label)
		}
		_, _, _, sk := Dec(k)
		pf("%9v: ", sk)
		mp := argMp[k]

		pf("%-4v", len(mp))

		if cntr%numCols == numCols-1 {
			pf("\n")
			newLineTerm = true
		} else {
			newLineTerm = false
		}
		cntr++
	}
	if len(argMp) > 0 && !newLineTerm {
		pf("\n")
	}
}

func (ar *Reservoir) PrintEdgesSummary() {
	// PrintEdgesRestrictTo("cond ", ar.EdgeCond, map[int]bool{})
	// PrintEdgesRestrictTo("no-sl", ar.EdgesSlackless, map[int]bool{})
	PrintEdgesRestrictTo("numel", ar.MElements, map[int]bool{})
}

func Granularize(t1 int) int {
	if t1 > 9 {
		t1 = 9
	} else if t1 > 4 {
		t1 = 4
	}
	return t1
}

// PrintPermutations is a mere demonstration,
// showing how we get more and exotic shapes
// if we run from maxEdge...Area,
// instead of 2...maxEdge.
// Notice, that second set completely contains first set
func PrintPermutations() {
	for i := 12; i < 24; i++ {
		maxEdge := RootCeiling(i)
		pf("%v: ", maxEdge)
		for j := 1; j <= maxEdge; j++ {
			k, slack := OtherSide(i, j)
			pf("%vx%2v-%v=%2v | ", j, k, slack, i)
		}
		pf("\n")
	}

	pf("now the superset:\n")
	cntr := 0
	for i := 12; i < 24; i++ {
		maxEdge := RootCeiling(i)
		pf("%v: ", maxEdge)
		for j := maxEdge; j <= i; j++ {
			k, slack := OtherSide(i, j)
			pf("%2vx%2v-%2v=%2v | ", j, k, slack, i)
			if cntr%10 == 9 {
				pf("\n   ")
			}
			cntr++
		}
		cntr = 0
		pf("\n")
	}
}

// PermutateEdges finds all possibilities to inscribe a slack
// into a given width-height.
// First,  the slack is taken as a *bar* of width 1.
// Second, the slack is fractionized.
// In all cases the slack is rotated by 90 degrees
// and the outer rectangle is also rotated by 90 degrees.
func PermutateEdges(cols, rows, slack int, plugs [][]int) [][]int {

	if slack > 0 {
		y := slack
		x := 1
		if rows > x && cols > y {
			plugs = append(plugs, []int{rows - x, y, x})
			plugs = append(plugs, []int{cols - y, x, y})
		}
		if rows > y && cols > x {
			plugs = append(plugs, []int{rows - y, x, y})
			plugs = append(plugs, []int{cols - x, y, x})
		}

		if !primes[slack] {
			factors := Factorize(slack)
			// for i := 0; i < len(factors)/2+1; i++ {
			for i := 0; i < len(factors); i++ {
				x := factors[i]
				y := slack / x
				// fmt.Printf(" fac%v*%v=%v\n", x, y, slack)
				if rows > x && cols > y {
					plugs = append(plugs, []int{rows - x, y, x})
					plugs = append(plugs, []int{cols - y, x, y})
				}
				if rows > y && cols > x {
					plugs = append(plugs, []int{rows - y, x, y})
					plugs = append(plugs, []int{cols - x, y, x})
				}
			}
		}
	}

	return plugs
}
