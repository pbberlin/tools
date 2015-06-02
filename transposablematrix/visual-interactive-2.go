package transposablematrix

// import "github.com/pbberlin/tools/tboxc"

func Main2() {

	w := 55
	h := 65
	M = NewTransposableMatrix(w, h)
	M.SetPerspective(perspect1)

	// M.FilledMinMax(0, true)
	// M.TotalMinMax()

	ar := NewReservoir()
	ar.AddAmorphs([]Amorph{Amorph{Cols: 6, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 9}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 2}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 13}})

	// because *set* SetAmorphSnap[Right/Left]wise() do not erase
	// the amorphs from the pools of edges, we do it ourselves
	for i := 0; i <= 4; i++ {
		ar.RemoveSpentEdges(ar.Amorphs[i])
	}

	ar.GenerateSpecificAmorphs([]int{8})
	ar.GenerateSpecificAmorphs([]int{2})

	ams := make([]Amorph, 5)
	ams[0] = ar.Amorphs[0]
	ams[1] = ar.Amorphs[1]
	ams[2] = ar.Amorphs[2]
	ams[3] = ar.Amorphs[3]
	ams[4] = ar.Amorphs[4]

	xmin, xmax := 0, 0
	ymin, ymax := 0, 0
	_, _, _ = ymin, ymax, xmax // sigh, sometimes ...
	l1 := []Point{}

	f := func() {
		M.Render2Termbox(-16, 16, -13, 13)
	}

	M.SetAmorphSnapRightwise(Point{-7, ams[0].Rows/2 + 1}, ams[0])
	// M.CastStitch( Point{0, 0}, ams[0])

	xmin, ymin, xmax, ymax = M.TotalMinMax(false)
	M.SetAmorphSnapLeftwise(Point{xmin, ams[1].Rows/2 + 1}, ams[1])
	f()

	xmin, ymin, xmax, ymax = M.TotalMinMax(false)
	M.SetAmorphSnapRightwise(Point{xmax + 1, ams[2].Rows/2 + 1}, ams[2])

	xmin, ymin, xmax, ymax = M.TotalMinMax(false)
	M.SetAmorphSnapRightwise(Point{xmax + 1, ams[3].Rows/2 + 1}, ams[3])

	xmin, ymin, xmax, ymax = M.TotalMinMax(false)
	M.SetAmorphSnapRightwise(Point{xmax + 1, ams[4].Rows/2 + 1}, ams[4])
	f()

	// fm := M.FilledMinMax( true)
	// f()

	// M.TotalMinMax( true)
	// f()

	l1, _ = M.PartialOutline()
	M.DrawLine(l1, "*", false)
	f()

	restrictHeuristicsByIndex = map[int]bool{}
	M.HeuristicsApply(&ar)
	f()

	restrictHeuristicsByIndex = map[int]bool{3: true, 4: true}
	M.HeuristicsApply(&ar)
	f()

	// l1 = M.PartialOutline(perspect1)
	// M.DrawLine( l1, "*", false)
	// f()

	// // M.ViewportRestriction( -11, -4, 20, 2)
	// // M.ViewportRestrictionForce(true)

	l1, _ = M.PartialOutline()
	M.DrawLine(l1, "*", false)
	M.SetPerspective(perspect2)
	M.Render2Termbox(-16, 16, -13, 13)

}
