package transposablematrix

type GrowShrink int

const (
	shrink     = GrowShrink(-1)
	grow       = GrowShrink(1)
	growShrink = GrowShrink(2)
)

func (gd GrowShrink) String() string {
	switch gd {
	case grow:
		return "grow  "
	case shrink:
		return "shrink"
	case growShrink:
		return "grwShr"
	}
	return ""
}

type VariDirection int

const (
	westw VariDirection = iota
	eastw
	northw
)

func (dir VariDirection) String() string {
	switch dir {
	case westw:
		return "westwd"
	case eastw:
		return "eastwd"
	case northw:
		return "northwd"
	}
	return ""
}

func (pdir *VariDirection) SwitchHoriz() {
	dir := *pdir
	if dir == westw {
		dir = eastw
	} else {
		dir = westw
	}
	*pdir = dir
}
