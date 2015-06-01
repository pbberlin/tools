package transposablematrix

// import "github.com/pbberlin/tools/tboxc"

func MainStairyShrinky() {

	w := 55
	h := 65
	M = NewTransposableMatrix(w, h)
	M.SetPerspective(perspect1)

	f := func() {
		M.Render2Termbox(-16, 16, -13, 13)
	}

	ar := NewReservoir()

	ar.AddAmorphs([]Amorph{Amorph{Cols: 7, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 3}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 7, Rows: 2}})

	lastAmorphIdx := len(ar.Amorphs) - 1
	lastAmorphIdx = 0
	nextA := func() Amorph {
		am := ar.Amorphs[lastAmorphIdx]
		lastAmorphIdx++
		return am
	}

	f()

	M.CastStitch(Point{-5, -5}, nextA())
	M.CastStitch(Point{2, -4}, nextA())
	M.CastStitch(Point{4, -3}, nextA())

	ar.GenerateRandomAmorphs(8)
	ar.SmallestDesirableHeight = 3

	M.ViewportRestriction(-6, -4, 11, 1)
	f()

	l1, _ := M.PartialOutline()
	M.DrawLine(l1, "*", false)
	f()

	M.DrawLine(l1, "", false)

	restrictHeuristicsByIndex = map[int]bool{1: true, 2: true, 3: true, 4: true}
	_, err := M.HeuristicsApply(&ar)
	if err != nil {
		pf(spf("%v", err))
	}
	f()

	_, err = M.HeuristicsApply(&ar)
	if err != nil {
		pf(spf("%v", err))
	}
	f()

}
