package transposablematrix

// import "github.com/pbberlin/tools/tboxc"

func MainTestNoEasternNeighbor() {

	w := 55
	h := 65
	M = NewTransposableMatrix(w, h)
	M.SetPerspective(perspect1)
	f := func() {
		M.Render2Termbox(-16, 16, -13, 13)
	}

	ar := NewReservoir()
	ar.AddAmorphs([]Amorph{Amorph{Cols: 7, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 3}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 4}})

	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 3}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 2}})

	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 3}})

	f()

	M.CastStitch(Point{-5, -3}, ar.Amorphs[0])
	M.CastStitch(Point{2, -3}, ar.Amorphs[1])
	M.CastStitch(Point{6, -4}, ar.Amorphs[2])
	M.CastStitch(Point{9, -5}, ar.Amorphs[3])

	M.CastStitch(Point{-5, -5}, ar.Amorphs[4])
	M.CastStitch(Point{-1, -5}, ar.Amorphs[5])
	M.CastStitch(Point{6, -7}, ar.Amorphs[6])
	M.CastStitch(Point{9, -7}, ar.Amorphs[7])

	M.CastStitch(Point{-5, -7}, ar.Amorphs[8])
	M.CastStitch(Point{-1, -8}, ar.Amorphs[9])

	ar.GenerateRandomAmorphs(25)
	ar.SmallestDesirableHeight = 3

	M.SetCenter(w/2-7, h/2)

	// M.CastStitch( Point{ 1,  -5}, ar.Amorphs[5])
	// M.CastStitch( Point{ 4,  -5}, ar.Amorphs[6])

	M.ViewportRestriction(-6, -4, 11, 1)
	f()

	l1, _ := M.PartialOutline()
	M.DrawLine(l1, "*", false)
	f()

	M.DrawLine(l1, "", false)

	_, err := M.HeuristicsApply(&ar)
	if err != nil {
		pf(spf("%v", err))
	}
	f()

}
