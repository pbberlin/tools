package transposablematrix

const (

	// Threshold for imperfect fills.
	// Leaving room for two or more.
	// Factor for MaxDesirableWidth.
	wideGapMin = 2

	// Dont bridge huge gaps with ByNumElementsWrap
	// Fallback to smaller steps instead
	// Factor for MaxDesirableWidth.
	wideGapCap = 3

	// Maximum limb width of a stairy edge
	// that is imperfectly filled.
	// Could be derived from wideGapCap?
	widestStatStair = 5
	widestDynStair  = 4

	// Minimum limb width of a stairy edge
	// that is imperfectly filled
	narrowestStair = 1 // uncontroversial

)

// For the despair heuristic FindmoreThanXElements().
const (
	minPossibleHeight = 1 // Marks the starting for an all out search
	cMaxDiff          = 20
)

// Search exactStraightEdge needs a constraint on
// how high the found amorphs may be.
// Param mainly for performance tuning.
const cMaxHeight = 8
