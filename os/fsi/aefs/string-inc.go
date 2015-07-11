package aefs

import (
	"fmt"
	"unicode/utf8"
)

func IncrementString(s string) string {

	if s == "" {
		panic("Increment String is undefined for an empty string")
	}

	uTFCodePointLastChar, itsSize := utf8.DecodeLastRuneInString(s)
	if uTFCodePointLastChar == utf8.RuneError {
		panic(fmt.Sprint("Following string is invalid utf8: %q", s))
	}
	sReduced := s[:len(s)-itsSize]

	uTFCodePointLastChar++
	oneHigherChar := fmt.Sprintf("%c", uTFCodePointLastChar)

	return sReduced + oneHigherChar

}
