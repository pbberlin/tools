package word

import lscore "github.com/pbberlin/tools/text/levenshtein"

type Token string // we could use []rune instead of string

func (internal Token) Matches(t1, t2 interface{}) bool {
	return t1 == t2
}

func convertToCore(sl1 []Token) []lscore.Token {
	var ret []lscore.Token
	for _, v := range sl1 {
		cnv := lscore.Token(v)
		ret = append(ret, cnv)
	}
	return ret
}
