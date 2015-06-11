package levenshtein

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
