package transposablematrix

// import "github.com/pbberlin/tools/tboxc"

func Main1() {

	w := 55
	h := 65
	M = NewTransposableMatrix(w, h)
	M.SetPerspective(perspect1)
	f := func() {
		M.Render2Termbox(-16, 16, -13, 13)
	}

	ar := NewReservoir()
	ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 1}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 6}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 5, Rows: 2}})

	f()

	ar.SmallestDesirableHeight = 3

	M.CastStitch(Point{-5, -1}, ar.Amorphs[0])
	M.CastStitch(Point{-2, -3}, ar.Amorphs[1])
	M.CastStitch(Point{0, -3}, ar.Amorphs[2])
	M.CastStitch(Point{2, -3}, ar.Amorphs[3])
	M.CastStitch(Point{4, -3}, ar.Amorphs[4])
	M.CastStitch(Point{6, 0}, ar.Amorphs[6])

	M.ViewportRestriction(-6, -4, 11, 1)
	f()

	l1, _ := M.PartialOutline()
	PrintOutline(l1)

	M.DrawLine(l1, "*", false)
	f()

	M.FilledMinMax(true)
	f()
	// M.TotalMinMax( true)
	// f()

}
