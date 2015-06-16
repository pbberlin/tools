package levenshtein

// Options control the way in which the levenshtein matrix
// is computed.
type Options struct {
	InsCost int
	DelCost int
	SubCost int
}

// DefaultOptions.
// Giving slight preference of insert/delete over substitution.
var DefaultOptions Options = Options{
	InsCost: 1,
	DelCost: 1,
	SubCost: 2, // SubCost can be 1, but not lower than InsCost - otherwise Editscript fails
}
