package levenshtein

import "fmt"

type EditOp int // EditOperation

const (
	Ins EditOp = iota
	Del
	Sub
	Match
)

func (o EditOp) String() string {
	if o == Match {
		return "match"
	} else if o == Ins {
		return "ins"
	} else if o == Sub {
		return "sub"
	}
	return "del"
}

type EditOpExt struct {
	op       EditOp
	src, dst int
}

func (op EditOpExt) String() string {
	return fmt.Sprintf("%v-%v-%v", op.op, op.src, op.dst)
}

type TEditScrpt []EditOpExt

func (es TEditScrpt) Print() {

	fmt2 := fmt.Sprintf("%s-%vv", "%", cl)

	fmt.Printf(fmt2, "EditScr:")
	// fmt.Printf("%v", strings.Repeat(" ", cl))

	for k, _ := range es {
		fmt.Printf(fmt2, k)
	}
	fmt.Printf("\n")

	fmt.Printf(fmt2, " ")
	for _, v := range es {
		fmt.Printf(fmt2, v)
	}
	fmt.Printf("\n")

}
