package transposablematrix

// unused, we use MatchFunc instead
type Matcher interface {
	AmorphFind(*Reservoir, Fusion) (*Amorph, Point)
}

type MatcherFunc func(*Reservoir, Fusion) (*Amorph, Point)

var AmorphsExhaustedError = epf("Amorphs exhausted")

var Heuristics = []MatcherFunc{
	StairyPerfect,     // 0
	StairyShrinky,     // 1
	StraightPerfect,   // 2
	StraightShrinky,   // 3
	ByNumElementsWrap, // 4
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
		var errTry error
		hits, errTry = m.IterateFusedSections(ar, clns, l)
		if errTry != nil {
			err = errTry
		}
	}

	return
}

func (m *TransposableMatrix) IterateFusedSections(ar *Reservoir,
	clns [][]int, l []Point) ([]int, error) {

	hits := make([]int, 2) // first: heuristics applied, second: number of results

	fs, err := m.FuseAllSections(ar, clns, l)
	if err != nil {
		return hits, err
	}

loopFusedSections:
	for i := 0; i < len(fs); i++ {

		// var curveRestrictions map[CurveDesc]bool
		// curveRestrictions = map[CurveDesc]bool{cncave: true}
		// if len(curveRestrictions) > 0 && !curveRestrictions[fs[i].curveDesc] {
		// 	pf("fusion curve %v restricted =>  continue\n", fs[i].curveDesc)
		// 	return nothing, err
		// }

		fs[i].Print()

		hits[0]++

		for iHeur := 0; iHeur < len(Heuristics); iHeur++ {

			if len(restrictHeuristicsByIndex) > 0 && !restrictHeuristicsByIndex[iHeur] {
				pf("heuristic %v restricted =>  continue\n", iHeur)
				continue
			}

			chosen, baseShift := Heuristics[iHeur](ar, fs[i])

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

	return hits, nil
}
