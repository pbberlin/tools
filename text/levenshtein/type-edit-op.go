package levenshtein

import (
	"fmt"
	"strings"
)

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

type TEditScrpt []EditOp

func (es TEditScrpt) Print() {
	fmt.Printf("EditScript - rowToCol - %v steps\n", len(es))
	fmt.Printf("%v", strings.Repeat(" ", cl))

	fmt2 := fmt.Sprintf("%s-%vv", "%", cl)

	for k, _ := range es {
		fmt.Printf(fmt2, k)
	}
	fmt.Printf("\n")

	fmt.Printf("%v", strings.Repeat(" ", cl))
	for _, v := range es {
		fmt.Printf(fmt2, v)
	}
	fmt.Printf("\n")
	fmt.Printf("\n")
}
