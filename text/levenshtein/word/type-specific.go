// Package levenshtein/word does tokenization on word level.
package word

import (
	"sort"

	"github.com/pbberlin/tools/pbstrings"
	ls_core "github.com/pbberlin/tools/text/levenshtein"
)

type Token string // we could use []rune instead of string

func (tk1 Token) Equal(compareTo interface{}) bool {
	tk2, ok := compareTo.(Token)
	if !ok {
		panic("Not the same type")
	}
	return tk1 == tk2
}

// WrapAsEqualer breaks string into a slice of strings.
// Each string is then converted to <Token> to <Equaler>.
// []<Equaler> can then be pumped into the generic core.
// We could as well create slices of Equalers in the first place
// but first leads to a var farTooUglyLiteral =
//   []ls_core.Equaler{ls_core.Equaler(Token("trink")), ls_core.Equaler(Token("nicht"))}
func WrapAsEqualer(s string, sorted bool) []ls_core.Equaler {
	ss := pbstrings.SplitByWhitespace(s)
	if sorted {
		sort.Strings(ss)
	}
	ret := make([]ls_core.Equaler, 0, len(ss))
	for _, v := range ss {
		cnv := ls_core.Equaler(Token(v))
		ret = append(ret, cnv)
	}
	return ret
}
