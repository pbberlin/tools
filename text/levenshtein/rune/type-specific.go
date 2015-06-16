// Package levenshtein/rune tokenizes on utf8 "codepoint" level.
package rune

import ls_core "github.com/pbberlin/tools/text/levenshtein"

type Token rune

func (tk1 Token) Equal(compareTo interface{}) bool {
	tk2, ok := compareTo.(Token)
	if !ok {
		panic("Not the same type")
	}
	return tk1 == tk2
}

// wrapAsEqualer wraps slice of tokens into interface type Equaler.
// Since our core implementation requires such slices.
func wrapAsEqualer(sl1 []Token) []ls_core.Equaler {
	var ret []ls_core.Equaler
	for _, v := range sl1 {
		cnv := ls_core.Equaler(v)
		ret = append(ret, cnv)
	}
	return ret
}
