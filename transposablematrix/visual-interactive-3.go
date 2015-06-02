package transposablematrix

func Main3() {

	w := 55
	h := 65
	M = NewTransposableMatrix(w, h)
	M.SetPerspective(perspect1)
	f := func() {
		M.Render2Termbox(-16, 16, -13, 13)
	}

	ar := NewReservoir()
	ar.GenerateRandomAmorphs(25)
	// for i := 0; i < len(ar.Amorphs); i++ {
	// 	ar.RemoveSpentEdges(ar.Amorphs[i])
	// }

	// ar.AddAmorphs([]Amorph{Amorph{Cols: 3, Rows: 1}})
	// M.CastStitch( Point{-5, -1}, ar.Amorphs[len(ar.Amorphs)-1])
	// // ar.RemoveSpentEdges(ar.Amorphs[0])

	// ar.SmallestDesirableHeight = 2
	ar.SmallestDesirableHeight = 3 // clearly the best choice - but why ?
	ar.SmallestDesirableWidth = 3  // clearly the best choice - but why ?

	ar.PrintEdgesSummary()

	M.ViewportRestriction(-6, -4, 11, 1)
	f()

	l1, _ := M.PartialOutline()
	M.DrawLine(l1, "*", false)
	f()

	M.DrawLine(l1, "", false)

	restrictHeuristicsByIndex = map[int]bool{2: true, 3: true, 4: true}

	for i := 0; i < 8; i++ {

		// ar.PrintEdgesSummary()
		_, err := M.HeuristicsApply(&ar)
		checkPrint(err)
		if err != nil {
			f()
			break
		}

		f()

		l1, _ = M.PartialOutline()
		// pf("%v",l1)
		M.DrawLine(l1, "*", false)
		f()

		M.DrawLine(l1, "", false)
		// f()

	}

	// M.ViewportRestrictionForce(false)

	restrictHeuristicsByIndex = map[int]bool{}

	M.SetCenter(-7, 0)

	for i := 0; i < 10; i++ {

		_, err := M.HeuristicsApply(&ar)
		checkPrint(err)
		if err != nil {
			f()
			break
		}

		f()
		l1, _ = M.PartialOutline()

		M.DrawLine(l1, "*", false)
		f()

		M.DrawLine(l1, "", false)
		// f()
	}

	M.FilledMinMax(true)
	f()
	M.TotalMinMax(true)
	f()

	// ----------------------

	M.SetPerspective(perspect2)
	M.FilledMinMax(true)
	M.TotalMinMax(true)
	M.Render2Termbox(-16, 16, -13, 13)

}
