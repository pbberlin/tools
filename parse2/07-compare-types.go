package parse2

type SortEl struct {
	Outl string
	Text []byte
}

type similarity struct {
	ArticleUrl     string
	Lvl            int
	Outline        string
	AbsLevenshtein int
	RelLevenshtein float64
	Text           []byte
}

type fragment struct {
	ArticleUrl string
	Lvl        int
	Outline    string
	Text       []byte
	Similars   []similarity
}
