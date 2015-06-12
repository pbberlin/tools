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

func (es TEditScrpt) Print(m Matrix) {

	fmt.Printf("EditScript - rowToCol - %v steps\n", len(es))
	fmt.Printf("%v", strings.Repeat(" ", cl))

	fmt2 := fmt.Sprintf("%s-%vv", "%", cl)

	for k, _ := range es {
		fmt.Printf(fmt2, k)
	}
	fmt.Printf("\n")

	sumIns := 0
	sumDel := 0
	rows2 := make([]Equaler, 0, len(m.rows))
	for _, v := range m.rows {
		rows2 = append(rows2, v)
	}

	const offs = 1
	fmt.Printf("%v", strings.Repeat(" ", cl))
	for _, v := range es {

		s := fmt.Sprintf("%v-%v-%v", v.op, offs+v.src+sumIns-sumDel, offs+v.dst)
		// s := fmt.Sprintf("%v-%v", v.op, v.dst)
		fmt.Printf(fmt2, s)

		if v.op == Ins {
			sumIns++
		}
		if v.op == Del {
			sumDel++
		}

		if v.op == Ins {
			// rows2 = InsertAfter(rows2, util.Min(v.src, len(rows2)-1), m.cols[v.dst])
			rows2 = InsertAfter(rows2, v.src-1, m.cols[v.dst])
		}
		if v.op == Del {
			rows2 = Delete(rows2, v.src)
		}
	}
	fmt.Printf("\n")
	fmt.Printf("%v", rows2)
	fmt.Printf("\n")
	fmt.Printf("\n")
}

// works also for idx == -1
func InsertAfter(s []Equaler, idx int, newVal Equaler) []Equaler {
	if idx > len(s)-1 {
		panic("Cannot insert beyond existing length")
	}
	s = append(s, nil)
	return s
	copy(s[idx+2:], s[idx+1:])
	s[idx+1] = newVal
	return s
}

func Delete(s []Equaler, idx int) []Equaler {
	if idx > len(s)-1 {
		panic("Cannot delete beyond existing length")
	}
	copy(s[idx:], s[idx+1:])
	s = s[:len(s)-1]
	return s
}
