package transposablematrix

// ByStairyEdge finds amorphs with a fitting edge.
// First slice contains amorphs with a perfect fit.
// If variation > 0:
// 		Change flankToVary by variation steps.
// 		Change according to param growShrink => increase/decr/both
func (ar *Reservoir) ByStairyEdge(fs Fusion, x1, y, x2 int,
	variation int, flankToVary VariDirection,
	growShrink GrowShrink) (amorphBlocks [][]Amorph, chosen *Amorph) {

	pfTmp := intermedPf(pf)
	defer func() { pf = pfTmp }()

	if variation < 0 {
		variation = -variation
		panic("variation assumed positive")
	}

	pf("srch %v;%v;%v ", x1, y, x2)

	limit := 1000 * 1000 // limit restriction dropped, but kept
	cntrAll := 0         // how many found

	firstBatch := ar.exactStairyEdge(x1, y, x2, limit)
	amorphBlocks = append(amorphBlocks, firstBatch)
	pf("fnd%v ", len(firstBatch))

	cntrAll += len(firstBatch)
	printNL := false

	if cntrAll < limit && variation != 0 {
		pf("| %v-%v-±%v | ", flankToVary, growShrink, variation)

		for i := 1; i <= variation; i++ {

			more1 := []Amorph{}
			more2 := []Amorph{}

			switch growShrink {
			case grow, shrink:
				offset := i * int(growShrink)
				pf("%+2v ", offset)
				switch flankToVary {
				case westw:
					more1 = ar.exactStairyEdge(x1+offset, y, x2, limit-cntrAll)
					pf("%v;%v;%v ", x1+offset, y, x2)
				case eastw:
					more1 = ar.exactStairyEdge(x1, y, x2+offset, limit-cntrAll)
					pf("%v;%v;%v ", x1, y, x2+offset)
				case northw:
					more1 = ar.exactStairyEdge(x1, y+offset, x2, limit-cntrAll)
				default:
					panic("unknown enum value for flankToVary")
				}
			case growShrink:

				offset1 := i * int(grow)
				offset2 := i * int(shrink)
				pf("±%v ", offset1)
				switch flankToVary {
				case westw:
					more1 = ar.exactStairyEdge(x1+offset1, y, x2, limit-cntrAll)
					pf("%v;%v;%v ", x1+offset1, y, x2)
					more2 = ar.exactStairyEdge(x1+offset2, y, x2, limit-cntrAll)
					pf("%v;%v;%v ", x1+offset2, y, x2)
				case eastw:
					more1 = ar.exactStairyEdge(x1, y, x2+offset1, limit-cntrAll)
					pf("%v;%v;%v ", x1, y, x2+offset1)
					more2 = ar.exactStairyEdge(x1, y, x2+offset2, limit-cntrAll)
					pf("%v;%v;%v ", x1, y, x2+offset2)
				case northw:
					more1 = ar.exactStairyEdge(x1, y+offset1, x2, limit-cntrAll)
					more2 = ar.exactStairyEdge(x1, y+offset2, x2, limit-cntrAll)
				default:
					panic("unknown enum value for flankToVary")
				}
			}

			if len(more1) > 0 { // dont append empty slices
				amorphBlocks = append(amorphBlocks, more1) // effects copying of the amorph
			}
			if len(more2) > 0 { // dont append empty slices
				amorphBlocks = append(amorphBlocks, more2) // effects copying of the amorph
			}

			cntrAll += len(more1)
			cntrAll += len(more2)

			pf("fnd%v  ", len(more1)+len(more2))

			if i%6 == 5 {
				pf("\n")
				printNL = true
			} else {
				printNL = false
			}

			if cntrAll >= limit {
				if printNL == false {
					pf("\n")
					printNL = true
				}
				pf("reached limit %2v in outer func. ret\n", limit)
				return
			}

		}
	}

	if printNL == false {
		pf("\n")
		printNL = true
	}

	chosen = activeFilter.Filter(amorphBlocks, fs)

	return
}
