package shared_structs

/*
	Order is filled with data from other structs.
	The slice index of the origin data records.
	And either an int or a string to sort

	Usage
	var sortedByX = make( []Order , len(sliceSrc))
	for i,v  := range sliceSrc {
		sortedByX[i].IdxSrc = i
		sortedByX[i].byI = v.Size  // or Salary ...
		sortedByX[i].byS = v.Lastname
	}
	bySize := ByInt(sortedByX)
	byLastname := ByStr(sortedByX)

	sort.Sort(byLastname)
	sort.Sort(bySize)		// last one wins

	The copying is necessary, since sort.Sort
	otherwise reorders the sliceSrc

*/
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
