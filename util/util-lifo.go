package util

import "fmt"

// LIFO
type Stack []string

func (s *Stack) Push(v string) {
	*s = append(*s, v)
}

func (s *Stack) Len() int {
	return len(*s)
}

func (s *Stack) Pop() string {
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
}

func (s Stack) String() string {
	return s.StringExt(false)
}

func (s Stack) StringExt(leafOnly bool) string {
	sep := "/"
	ret := ""
	for i := 0; i < len(s); i++ {
		if leafOnly {
			if i == len(s)-1 {
				ret = fmt.Sprintf("%v%s%-5v", ret, sep, s[i])
			} else {
				ret = fmt.Sprintf("%v%s%-5v", ret, " ", "")
			}
		} else {
			ret = fmt.Sprintf("%v%s%v", ret, sep, s[i])
		}
	}
	return ret
}
