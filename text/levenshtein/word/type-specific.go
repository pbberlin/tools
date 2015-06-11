package word

import ls_core "github.com/pbberlin/tools/text/levenshtein"

type Token string // we could use []rune instead of string

func (tk1 Token) Equal(compareTo interface{}) bool {
	tk2, ok := compareTo.(Token)
	if !ok {
		panic("Not the same type")
	}
	return tk1 == tk2
}

func convertToCore(sl1 []Token) []ls_core.Equaler {
	var ret []ls_core.Equaler
	for _, v := range sl1 {
		cnv := ls_core.Equaler(v)
		ret = append(ret, cnv)
	}
	return ret
}
