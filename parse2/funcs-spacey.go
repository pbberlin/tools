package parse2

import "strings"

var replNewLines = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ")

func isSpacey(sarg string) bool {
	s := sarg
	s = replNewLines.Replace(s)
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return false
}
