package domclean2

import (
	"regexp"
	"strings"
)

var replNewLines = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ")
var replTabs = strings.NewReplacer("\t", " ")

var doubleSpaces = regexp.MustCompile("([ ]+)")

// isSpacey detects all possible occurrences of whitespace
func isSpacey(sarg string) bool {
	s := strings.TrimSpace(sarg) // TrimSpace removes leading-trailing \n \r\n
	if s == "" {
		return true
	}
	return false
}

// All kinds of newlines, tabs and double spaces
// are reduced to single space.
// It paves the way for later beautification.
func textNormalize(s string) string {
	s = replNewLines.Replace(s)
	s = replTabs.Replace(s)
	s = doubleSpaces.ReplaceAllString(s, " ")
	return s
}
