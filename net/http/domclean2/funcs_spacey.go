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
