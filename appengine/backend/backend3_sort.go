package backend

type Order struct {
	IdxSrc int    // index to the base slice
	ByI    int    // the base for sorting by int
	ByS    string // the base for sorting by string
}
type ByInt []Order
type ByStr []Order

// 	We could condense the  .Len() and .Swap() methods.
func (s ByInt) Len() int           { return len(s) }
func (s ByInt) Less(i, j int) bool { return s[i].ByI < s[j].ByI }
func (s ByInt) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s ByStr) Len() int           { return len(s) }
func (s ByStr) Less(i, j int) bool { return s[i].ByS < s[j].ByS }
func (s ByStr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
