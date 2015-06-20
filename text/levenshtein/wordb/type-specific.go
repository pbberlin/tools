// Package levenshtein/wordb does tokenization on word level,
// using byte slices instead of string, saving conversion cost.
package wordb

import (
	"bytes"
	"sort"

	ls_core "github.com/pbberlin/tools/text/levenshtein"
)

type Token []byte // similar to word, but without string

func (tk1 Token) Equal(compareTo interface{}) bool {
	tk2, ok := compareTo.(Token)
	if !ok {
		panic("Not the same type")
	}
	return bytes.Equal(tk1, tk2) // bytes.EqualFold would make it case insensitive
}

// See word.WrapAsEqualer
func WrapAsEqualer(sb []byte, sorted bool) []ls_core.Equaler {

	sbf := bytes.Fields(sb)
	if sorted {
		sort.Sort(sortBoB(sbf))

		// weed out doublettes
		su, prev := make([][]byte, 0, len(sbf)), []byte{}
		for _, v := range sbf {
			if bytes.Equal(v, prev) {
				continue
			}
			su = append(su, v)
			prev = v
		}
		sbf = su

	}

	ret := make([]ls_core.Equaler, 0, len(sbf))
	for _, v := range sbf {
		cnv := ls_core.Equaler(Token(v))
		ret = append(ret, cnv)
	}
	return ret
}

type sortBoB [][]byte // slice of bytes of bytes

func (sb sortBoB) Len() int {
	return len(sb)
}

func (sb sortBoB) Less(i, j int) bool {
	return bytes.Compare(sb[i], sb[j]) > 0
}

func (sb sortBoB) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}

func sortBBSlice(bs [][]byte) {
	sort.Sort(sortBoB(bs))
}
