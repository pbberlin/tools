package levenshtein

type Options struct {
	InsCost int
	DelCost int
	SubCost int
}

// DefaultOptions is the default options: insertion cost is 1, deletion cost is
// 1, substitution cost is 2, and two Token match iff they are the same.
var DefaultOptions Options = Options{
	InsCost: 1,
	DelCost: 1,
	SubCost: 2,
}
