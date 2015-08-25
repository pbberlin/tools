package weedout

type SortEl struct {
	Outl string
	Text []byte
}

type TextifiedTree struct {
	SourceID  string
	Lvl       int
	Outline   string
	NumTokens int
	Text      []byte

	Similars          []Similar
	SumAbsLevenshtein int
	SumRelLevenshtein float64
}

// Similarity relationship towards another TextifiedTree
type Similar struct {
	SourceID string
	Lvl      int
	Outline  string
	Text     []byte

	AbsLevenshtein int
	RelLevenshtein float64
}
