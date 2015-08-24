package weedout

type SortEl struct {
	Outl string
	Text []byte
}

// Similarity relationship towards another TextifiedTree
type Similar struct {
	ArticleUrl     string
	Lvl            int
	Outline        string
	AbsLevenshtein int
	RelLevenshtein float64
	Text           []byte
}

type TextifiedTree struct {
	ArticleUrl string
	Lvl        int
	Outline    string
	Text       []byte
	Similars   []Similar
}
