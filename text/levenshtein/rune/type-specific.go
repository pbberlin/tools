package rune

import lscore "github.com/pbberlin/tools/text/levenshtein"

type Token rune

func (internal Token) Matches(t1, t2 interface{}) bool {
	return t1 == t2
}

func toTokenInterface(sl1 []Token) []lscore.Token {
	var ret []lscore.Token
	for _, v := range sl1 {
		cnv := lscore.Token(v)
		ret = append(ret, cnv)
	}
	return ret
}
