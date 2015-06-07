package parse2

import (
	"regexp"
	"strings"
)

var replNewLines = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ")
var replTabs = strings.NewReplacer("\t", " ")

var doubleSpaces = regexp.MustCompile("([ ]+)")

func isSpacey(sarg string) bool {
	s := sarg
	s = replNewLines.Replace(s)
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return false
}
