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

type EditOpExt struct {
	op       EditOp
	src, dst int
}

type TEditScrpt []EditOpExt

func (es TEditScrpt) Print() {
	fmt.Printf("EditScript - rowToCol - %v steps\n", len(es))
	fmt.Printf("%v", strings.Repeat(" ", cl))

	fmt2 := fmt.Sprintf("%s-%vv", "%", cl)

	for k, _ := range es {
		fmt.Printf(fmt2, k)
	}
	fmt.Printf("\n")

	sumIns := 0
	sumDel := 0

	fmt.Printf("%v", strings.Repeat(" ", cl))
	for _, v := range es {
		v.src++
		v.dst++

		s := fmt.Sprintf("%v-%v-%v", v.op, v.src+sumIns-sumDel, v.dst)
		// s := fmt.Sprintf("%v-%v", v.op, v.dst)
		fmt.Printf(fmt2, s)

		if v.op == Ins {
			sumIns++
		}
		if v.op == Del {
			sumDel++
		}

	}
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("\n")
}
