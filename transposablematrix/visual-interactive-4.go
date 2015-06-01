package transposablematrix

func Main4() {

	w := 55
	h := 65
	M = NewTransposableMatrix(w, h)
	M.SetPerspective(perspect1)

	ar := NewReservoir()
	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 3}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 7, Rows: 1}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 4, Rows: 4}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 1}})
	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 3}})
	// Because *set* SetAmorphSnap[Right/Left]wise() do not erase
	// the amorphs from the pools of edges, we do it ourselves
	for i := 0; i < len(ar.Amorphs); i++ {
		ar.RemoveSpentEdges(ar.Amorphs[i])
	}

	l1 := []Point{}
	f := func() {
		M.Render2Termbox(-16, 16, -13, 13)
	}

	fh := func() {
		_, err := M.HeuristicsApply(&ar)
		checkPrint(err)
		f()

		l1, _ = M.PartialOutline()
		M.DrawLine(l1, "*", false)
		f()
		M.DrawLine(l1, "", false)
	}

	M.ViewportRestriction(-7, -4, 14, 2)

	M.SetAmorphSnapRightwise(Point{-5, 2}, ar.Amorphs[0])
	M.SetAmorphSnapRightwise(Point{-1, 2}, ar.Amorphs[1])
	M.SetAmorphSnapRightwise(Point{6, 2}, ar.Amorphs[2])
	M.SetAmorphSnapRightwise(Point{10, 2}, ar.Amorphs[3])
	M.SetAmorphSnapRightwise(Point{12, 2}, ar.Amorphs[4])
	M.SetAmorphSnapRightwise(Point{12, 2}, ar.Amorphs[4])

	ar.AddAmorphs([]Amorph{Amorph{Cols: 2, Rows: 1}})
	M.SetAmorphSnapRightwise(Point{0, 1}, ar.Amorphs[5])
	ar.AddAmorphs([]Amorph{Amorph{Cols: 1, Rows: 1}})
	M.SetAmorphSnapRightwise(Point{3, 1}, ar.Amorphs[6])

	l1, _ = M.PartialOutline()
	M.DrawLine(l1, "*", false)
	f()

	l1, _ = M.PartialOutline()
	M.DrawLine(l1, "", false)

	// ar.GenerateSpecificAmorphs([]int{8})
	// _,errH = M.HeuristicsApply( &ar)
	f()

	ar.GenerateSpecificAmorphs([]int{5, 7, 8, 9, 10})
	ar.GenerateSpecificAmorphs([]int{22})
	// ar.GenerateSpecificAmorphs([]int{19})
	ar.GenerateSpecificAmorphs([]int{16})
	// ar.GenerateSpecificAmorphs([]int{11, 12, 13, 14, 15})
	// ar.GenerateSpecificAmorphs([]int{21, 22, 23, 24, 25, 26})
	// ar.GenerateSpecificAmorphs([]int{29})

	fh()

	ar.GenerateSpecificAmorphs([]int{11})
	ar.GenerateSpecificAmorphs([]int{31})

	for i := 0; i < 10; i++ {
		fh()
	}

}
