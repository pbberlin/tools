package transposablematrix

import (
	"math"
	"math/rand"
	"time"
)

const (
	initSize = 100 // preallocated slices
)

// Now it seems, Spent/Map is not neccessary,
// as long as we clean up
//   1.) After seeded first 5 amorphs => RemoveSpentEdges()
//   1.) After ApplyHeuristic has succeeded => RemoveSpentEdges()
var Spent = map[int]bool{} // this should belong into Reservoir, but must be accessed from Matrix-Methods - I hesitate to insert additional arguments into all the methods, befor I am certain

type Amorph struct {
	IdxArticle     int // reference 'upward', since this type is also used in the 'matrix' later on
	NElements      int
	Cols           int // width
	Rows           int // height
	Slack          int
	Edge           []int
	AestheticValue int // 100 - abs((ratio -  goldenRatio)*100)

	Padded int // Elements added to suit into an otherwise unfillable gap
}

type Reservoir struct { // Reservoir of Amorphs
	Amorphs []Amorph

	AvgWidth                float64
	SdvWidth                float64
	AvgHeight               float64
	SdvHeight               float64
	SmallestDesirableWidth  int
	SmallestDesirableHeight int

	// Width...
	HistoWidths []int // serves for WidthOffsets, serves to sum squares for average and standard dev
	SumSqWidths int

	// Height...
	HistoHeights []int
	SumSqHeights int

	// Index structs:
	// =========================
	// M(ap)Elements - named to distinguish from common N(um)Elements
	// Also acts as counter to measure spent amorphs.
	MElements map[int]map[int]bool

	Edge3          map[int]map[int]bool
	EdgeCond       map[int]map[int]bool // condensed, 5,6,...9 => condensed to 5
	EdgesSlackless map[int]map[int]bool // only slack==0
}

func NewReservoir() Reservoir {

	AR := Reservoir{}

	AR.HistoWidths = make([]int, initSize)
	AR.HistoHeights = make([]int, initSize)

	AR.MElements = map[int]map[int]bool{}

	AR.Edge3 = map[int]map[int]bool{}
	AR.EdgeCond = map[int]map[int]bool{}
	AR.EdgesSlackless = map[int]map[int]bool{}

	return AR
}

// RootCeiling gives an efficient limit
// for edge permutations of rectangles.
// Assuming rectangles have integer areal.
// Todo: We could/should build a fixed-limited or a synchronized global lookup table.
func RootCeiling(nElements int) (rootCeil int) {
	rootCeil = int(math.Ceil(math.Sqrt(float64(nElements))))
	return
}

// OtherSide tries the most efficient way
// to compute the integer value *second* edge of a rectangle,
// given the first edge and rectangle area
func OtherSide(nElements, cols int) (rows, slack int) {
	if nElements%cols == 0 {
		rows = nElements / cols
	} else {
		rows = (nElements / cols) + 1 // deliberately unidiomatic, saving one op in 50pct of cases
	}
	slack = cols*rows - nElements
	return
}

// AestheticValue() yields how "beautiful" a rectangle looks.
// Higher return values mean hotter looks.
// While values greater zero are almost equally acceptable,
// return values smaller zero indicate considerable ugliness.
func AestheticValue(width, height int) int {

	// Everything is multiplied by 10 to keep integer format.
	// Portrait or landscape ratio gte 3:1 gives positive number,
	// i.e. 3cm width, 1cm height or vice versa yield plus 1.
	// More extreme ratios turn negative
	const goodRatioApex = 301

	ratio := 0
	if width > height {
		ratio = int(100 * float32(width) / float32(height))
	} else {
		ratio = int(100 * float32(height) / float32(width))
	}
	ret := -(ratio - goodRatioApex)

	//
	// degression for "extreme ugly" values
	if ret < -50 {
		ret = -50 + ret/10
	}
	if ret < -99 {
		ret = -99
	}
	return ret

}

// GenerateRandomAmorphs creates a slice of integers
// 50 percent are in the range of 1...8
// 50 percent are in the range of 4...11
func (ar *Reservoir) GenerateRandomAmorphs(numExampleAmorphs int) {

	nums := make([]int, numExampleAmorphs)
	for i := 0; i < numExampleAmorphs; i++ {
		num := 1 + rand.Intn(9)
		if rand.Intn(4) > 2 {
			num = 4 + rand.Intn(12)
		}
		nums[i] = num
	}
	ar.GenerateSpecificAmorphs(nums)

}

// GenerateSpecificAmorphs takes a slice of desired square areas.
// It then creates the most equal-sided edges for such square area.
func (ar *Reservoir) GenerateSpecificAmorphs(nums []int) {

	newAmorphs := make([]Amorph, len(nums))
	for i := 0; i < len(nums); i++ {
		lp := &newAmorphs[i]
		lp.NElements = nums[i]
		lp.Cols = RootCeiling(lp.NElements) // => avgWidth slightly larger than avgHeight
		lp.Rows, lp.Slack = OtherSide(lp.NElements, lp.Cols)
	}
	ar.AddAmorphs(newAmorphs)

}

// AddAmorphs() computes derived values, we want to rely on
// And it updates all sorts of statistics and indexes,
// we want to query.
func (ar *Reservoir) AddAmorphs(newA []Amorph) {

	// append the amorphs themselves
	lenPrev := len(ar.Amorphs)
	ar.Amorphs = append(ar.Amorphs, newA...)

	// update the base structs for summing and sorting
	for i := lenPrev; i < len(ar.Amorphs); i++ {
		lp := &ar.Amorphs[i]

		if lp.IdxArticle == 0 {
			lp.IdxArticle = i
		}

		if lp.AestheticValue == 0 {
			lp.AestheticValue = AestheticValue(lp.Cols, lp.Rows)
		}

		if lp.NElements == 0 {
			lp.NElements = lp.Cols*lp.Rows - lp.Slack
		}

		ar.HistoWidths[lp.Cols]++
		ar.HistoHeights[lp.Rows]++
		ar.SumSqWidths += lp.Cols * lp.Cols
		ar.SumSqHeights += lp.Rows * lp.Rows

		ar.PermutateShapes(*lp, false)
		SubmapAddRemove(false, ar.MElements, lp.NElements, lp.IdxArticle)

	}

	// update global statistics
	ar.AvgWidth, ar.SdvWidth = avgSdvExplicit(ar.HistoWidths, ar.SumSqWidths)
	ar.AvgHeight, ar.SdvHeight = avgSdvExplicit(ar.HistoHeights, ar.SumSqHeights)
	ar.SmallestDesirableWidth = ar.updateDesirableGap(0)
	ar.SmallestDesirableHeight = ar.updateDesirableGap(1)

}

// transpose puts the *inner* coordinates of the matrix slots into perspective
func (a *Amorph) transpose(xm, ym, perspective int) (xb, yb int) {
	switch perspective {
	case 0:
		xb, yb = xm, ym
	case 1:
		xb, yb = ym, a.Cols-xm-1
	case 2:
		xb, yb = a.Cols-xm-1, a.Rows-ym-1
	case 3:
		xb, yb = a.Rows-ym-1, xm
	}
	return
}

// Recursively computing avg and std deviation
// http://math.stackexchange.com/questions/374881/recursive-formula-for-variance
// However, we denominate n, n-1 => our divisior is therefore n
//
// As we maintain current sums (in the form of histograms)
// and sum of squares as fields in the Reservoir,
// we only use explicit computation
func unused_avgSdvRecursive(n, newVal int, avgPrev, stdPrev float64) (avgNew float64, newStddev float64) {

	nfl := float64(n) // n as float

	nvfl := float64(newVal) // new value as float

	avgNew = avgPrev + (nvfl-avgPrev)/nfl

	sqNewStddev := (stdPrev * stdPrev) +
		(avgPrev * avgPrev) -
		(avgNew * avgNew) +
		(nvfl*nvfl-stdPrev*stdPrev-avgPrev*avgPrev)/nfl
	newStddev = math.Sqrt(sqNewStddev)

	return
}

// avgSdvExplicit() computes average and standard deviation from scratch.
// It relies on an up-to-date histogram and sum of squares.
func avgSdvExplicit(histo []int, sumSquares int) (avg float64, sdv float64) {
	count := 0
	sumV := 0
	for i := 0; i < len(histo); i++ {
		count += histo[i]
		sumV += i * histo[i]
	}
	avg = float64(sumV) / float64(count)
	sdv = math.Sqrt(float64(sumSquares)/float64(count) - (avg * avg))
	return
}

// updateDesirableGap() gives a dynamic look into the remaining
// amorphs available for cast and stitch.
// We try giving back the smallest edge length or gap size
// for which a large quantity of amorphs is still available.
// I initially tampered with the standard deviation,
// but found no way to reasonably integrate it.
// Todo: adapt avg and stddev after consumption of amorphs
func (ar *Reservoir) updateDesirableGap(perspective int) int {

	// pf("%3.2v %3.2v - %3.2v %3.2v\n", ar.AvgWidth, ar.SdvWidth, ar.AvgHeight, ar.SdvHeight)

	const scale = 0.9

	tmp := 0
	if perspective%2 == 1 {
		tmp = int(ar.AvgHeight * scale)
	} else {
		tmp = int(ar.AvgWidth * scale)
	}

	if tmp == 0 {
		return 1
	}

	tmp = 2

	return tmp
}

// Takes any slice of amorphs and prints them.
// Parameter "mode" differentiates between print formats
func DumpAmorphs(ams []Amorph, mode int) {

	itemsPerLine := 15

	for i := 0; i < len(ams); i++ {
		switch mode {
		case 1: // idx, cols, rows, slack
			itemsPerLine = 8
			id := spf("#%v", ams[i].IdxArticle)
			pf("%4v: %2vx%02v%%%02v |", id, ams[i].Cols, ams[i].Rows, ams[i].Slack)
		case 2: // idx, cols, rows, slack, edge
			itemsPerLine = 6
			id := spf("#%v", ams[i].IdxArticle)
			pf("%4v: %2vx%02v%%%02v %v |", id, ams[i].Cols, ams[i].Rows, ams[i].Slack, ams[i].Edge)
		case 3: // idx, cols, rows, slack, aestetic val
			itemsPerLine = 7
			id := spf("#%v", ams[i].IdxArticle)
			pf("%4v: %2vx%02v%%%02v %3d |", id, ams[i].Cols, ams[i].Rows, ams[i].Slack, ams[i].AestheticValue)
		case 4: // no idx
			pf("%2vx%02v %2v |", ams[i].Cols, ams[i].Rows, ams[i].Slack)
		case 7: // only idx Article + NElements
			itemsPerLine = 12
			pf("%2v %02v |", ams[i].IdxArticle, ams[i].NElements)
		case 2222:
			itemsPerLine = 4
			if false {
				break // this breaks the *switch*, not the *for*; undocumented but necessary
			}
			fallthrough
		default:
			pf("%2vx%02v!%2v |", ams[i].Cols, ams[i].Rows, ams[i].Slack)
		}
		if i%itemsPerLine == itemsPerLine-1 && i != len(ams)-1 {
			pf("\n")
		}
	}
	if len(ams) > 0 {
		pf("\n")
	}
}

func (ar *Reservoir) PrintAmorphSummary(label string) {

	argMp := ar.MElements

	cntr := 0
	numCols := 10
	newLineTerm := true

	const gran = 5
	const delim = 5

	for i := 0; i < gran*delim; i++ {
		key := Enc(0, 0, i)
		if v, ok := argMp[key]; ok {
			if newLineTerm {
				pf("%s", label)
			}
			pf("%2v:", i)
			pf("%2v|", len(v))
			if cntr%numCols == numCols-1 {
				pf("\n")
				newLineTerm = true
			} else {
				newLineTerm = false
			}
			cntr++
		}
	}

	for i := delim; i < 20; i++ {
		found := 0
		for j := i * gran; j < (i+1)*gran; j++ {
			key := Enc(0, 0, j)
			if v, ok := argMp[key]; ok {
				found += len(v)
			}
		}
		if found > 0 {
			if newLineTerm {
				pf("%s", label)
			}

			pf("%2v-%2v:", i*gran, (i+1)*gran)
			pf("%2v|", found)
			if cntr%numCols == numCols-1 {
				pf("\n")
				newLineTerm = true
			} else {
				newLineTerm = false
			}
			cntr++
		}
	}
	if len(argMp) > 0 && !newLineTerm {
		pf("\n")
	}

}

func (f *Amorph) ReverseEdge() {
	if f.Edge != nil && len(f.Edge) > 2 {
		f.Edge[0], f.Edge[2] = f.Edge[2], f.Edge[0]
		f.Edge[1] = -f.Edge[1]
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
