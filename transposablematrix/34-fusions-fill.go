package transposablematrix

// no receiver (*Reservoir)Find
// dont know how to apply interface
type AmorphFinder interface {
	AmorphFind(*Reservoir, int, int, int, int, int) ([][]Amorph, *Amorph, Point)
}

type AmorphFinderFunc func(*Reservoir, Fusion) ([][]Amorph, *Amorph, Point)

// Stack of heuristics - Stack No 1
var concaveHeuristics = []AmorphFinderFunc{
	StairyPerfectConcave, // 0
	StairyShrinkyConcave, // 1
	StraightPerfect,      // 2
	StraightShrinky,      // 3
	ByNumElementsWrap,    // 4
}

// Stack of heuristics - Stack No 2
var allHeuristics = []AmorphFinderFunc{
	StairyPerfect,     // 0
	StairyShrinky,     // 1
	StraightPerfect,   // 2
	StraightShrinky,   // 3
	ByNumElementsWrap, // 4
}

// Combinining the stacks, depending the curve description
var applicable = map[CurveDesc][]AmorphFinderFunc{
	cncave: concaveHeuristics,
	stairW: allHeuristics,
	stairE: allHeuristics,
	convex: allHeuristics,
}

// var restrictHeuristicsByIndex = map[int]bool{0: true, 1: true, 2: true}
// this map is no longer serving its purpose well
// We should instead modify the "applicable" map
var restrictHeuristicsByIndex = map[int]bool{}

var reduceHeightForTestingCheat bool

func (m *TransposableMatrix) HeuristicsApply(ar *Reservoir) (hits []int, err error) {

	clns, l, errCNLS := m.ConcavestLowestNarrowestSections()
	if errCNLS != nil {
		err = errCNLS
		return
	}

	ar.PrintAmorphSummary("# el:")
	pf("l   : ")
	PrintOutline(l)
	PrintCLNS(clns)

	if len(ar.MElements) < 1 {
		err = AmorphsExhaustedError
	} else {
		_, errTry := m.IterateFusedSections(ar, clns, l)
		if errTry != nil {
			err = errTry
		}
	}

	return
}

func (m *TransposableMatrix) IterateFusedSections(ar *Reservoir,
	clns [][]int, l []Point) (SearchResult, error) {

	hits := make([]int, 2)

	fs, err := m.FuseAllSections(ar, clns, l)
	if err != nil {
		return nothing, err
	}

	// for i := 0; i < len(fs); i++ {
	// 	fs[i].Print()
	// }

loopFusedSections:
	for i := 0; i < len(fs); i++ {

		var heuristics []AmorphFinderFunc

		heuristics = applicable[fs[i].curveDesc]

		// var curveRestrictions map[CurveDesc]bool
		// curveRestrictions = map[CurveDesc]bool{cncave: true}
		// if len(curveRestrictions) > 0 && !curveRestrictions[fs[i].curveDesc] {
		// 	pf("fusion curve %v restricted =>  continue\n", fs[i].curveDesc)
		// 	return nothing, err
		// }

		fs[i].Print()

		hits[0]++

		for iHeur := 0; iHeur < len(heuristics); iHeur++ {

			if len(restrictHeuristicsByIndex) > 0 && !restrictHeuristicsByIndex[iHeur] {
				pf("heuristic %v restricted =>  continue\n", iHeur)
				continue
			}

			_, chosen, baseShift := heuristics[iHeur](ar, fs[i])

			if chosen != nil {

				hits[1]++
				fs[i].base.x += baseShift.x
				fs[i].base.y += baseShift.y
				// base.y += -2

				m.SetAmorphSnapRightwise(fs[i].base, *chosen)
				ar.RemoveSpentEdges(*chosen)

				pf("!! iHeur%v: success (id%v els%v %v*%vs%v) %v| \n", iHeur, chosen.IdxArticle, chosen.NElements, chosen.Edge, chosen.Rows, chosen.Slack, baseShift)

				// we *must break upon each successfull inscripture
				// and recalculate the new outline
				break loopFusedSections
			} else {
				if iHeur > 0 {
					pf("!! iHeur%v: failure | \n", iHeur)
				} else {
					pf("sameline !! iHeur%v: failure | \n", iHeur)
				}
			}
		}

	}

	result := nothing
	if hits[0] > 0 && hits[1] == 0 {
		result = failure
	} else if hits[0] > 0 && hits[1] > 0 {
		result = success
	}
	return result, nil
}
