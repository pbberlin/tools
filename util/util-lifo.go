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
	ret := ""
	for i := 0; i < len(s); i++ {
		ret = fmt.Sprintf("%v / %v", ret, s[i])
	}
	return ret
}
