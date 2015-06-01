package transposablematrix

//
// Find the perfect fit for a given edge x-y-x
func StairyPerfectConcave(ar *Reservoir, fs Fusion) (amorphBlocks [][]Amorph, chosen *Amorph, baseShift Point) {

	var x1, y, x2 = fs.xyx[0], fs.xyx[1], fs.xyx[2]

	var stepdown bool
	if y < 0 {
		stepdown = true
		x1, x2 = x2, x1
		y = -y
	}

	var fitting [][]Amorph
	fitting, _ = ar.ByStairyEdge(x1, y, x2, 0, westw, grow)

	chosen = abundantHeightMatch(fitting, fs)

	// switch back
	if stepdown && chosen != nil {
		chosen.Edge[0], chosen.Edge[2] = chosen.Edge[2], chosen.Edge[0]
		chosen.Edge[1] = -chosen.Edge[1]
		x1, x2 = x2, x1
		y = -y
	}

	return nil, chosen, Point{}

}

// We need to distinguish four cases
// 		stepup				or	stepdown
//		LowStat-HighDyn		or	LowDyna-HighSta-
// Todo: limit height
func StairyShrinkyConcave(ar *Reservoir, fs Fusion) (amorphBlocks [][]Amorph, chosen *Amorph, baseShift Point) {

	var x1, y, x2, directionIdx, maxOffs = fs.xyx[0], fs.xyx[1], fs.xyx[2], fs.dirIdx, fs.maxOffs
	_, _ = directionIdx, maxOffs

	if y < 0 {
		chosen = AdapterStepdownAndDirection(ar, "LowStat-HighDyn-", x1, y, x2, 0)

		if chosen == nil {
			chosen = AdapterStepdownAndDirection(ar, "LowDyna-HighSta-", x1, y, x2, 1)
		}

	} else {
		chosen = AdapterStepdownAndDirection(ar, "LowStat-HighDyn-", x1, y, x2, 1)

		if chosen == nil {
			chosen = AdapterStepdownAndDirection(ar, "LowDyna-HighSta-", x1, y, x2, 0)
		}
	}

	if chosen != nil {
		baseShift.x += x1 - chosen.Edge[0]
		// baseShift.y--
	}

	return

}

// First, recollect following equivalence:
//
//  XXXXX-<-                ->-XXXXX
//  XXX          equiv           XXX
//  XXX          mirror          XXX
//  3;2;2...               ...2;-2;3
//  westw                      eastw
//
// The shrink *direction* is switched too
//
//  ->-XXXXX                  XXXXX-<-
//  ->-X          equiv           X-<-
//  ->-X          mirror          X-<-
//  3..1;2;4                  4;-2;3..1
//  eastw                     westw
//
func AdapterStepdownAndDirection(ar *Reservoir, desc string, x1, y, x2 int, direction int) (chosen *Amorph) {

	var stepdown bool
	if y < 0 {
		stepdown = true
		x1, x2 = x2, x1
		y = -y
		if direction == 0 {
			direction = 1
		} else {
			direction = 0
		}
	}

	if direction == 0 {
		if x1 >= wideGapMin*ar.SmallestDesirableWidth {
			pf("gap%v wider   than%v*%v =>    %vStairyShrinky ", x1, wideGapMin, ar.SmallestDesirableWidth, desc)
			// leave at least SmallestDesirableWidth for further fill
			xdyn := x1 - ar.SmallestDesirableWidth
			if xdyn > widestDynStair+1 {
				xdyn = widestDynStair
			}
			if x2 > widestStatStair+1 {
				x2 = widestStatStair
			}
			// we want to stop at narrowestStair, thus:
			maxOffs := xdyn - narrowestStair
			pf("%v...%v\n", xdyn, xdyn-maxOffs)
			_, chosen = ar.ByStairyEdge(xdyn, y, x2, maxOffs, westw, shrink)
		} else {
			pf("gap%v narrowr than%v*%v => no %vStairyShrinky \n", x1, wideGapMin, ar.SmallestDesirableWidth, desc)
		}
	} else {

		if x2 >= wideGapMin*ar.SmallestDesirableWidth {
			pf("gap%v wider   than%v*%v =>    %vStairyShrinky ", x2, wideGapMin, ar.SmallestDesirableWidth, desc)
			// leave at least SmallestDesirableWidth for further fill
			xdyn := x2 - ar.SmallestDesirableWidth
			if xdyn > widestDynStair+1 {
				xdyn = widestDynStair
			}
			if x1 > widestStatStair+1 {
				x1 = widestStatStair
			}
			// we want to stop at narrowestStair, thus:
			maxOffs := xdyn - narrowestStair
			pf("%v...%v\n", xdyn, xdyn-maxOffs)
			_, chosen = ar.ByStairyEdge(x1, y, xdyn, maxOffs, eastw, shrink)
		} else {
			pf("gap%v narrowr than%v*%v => no %vStairyShrinky \n", x2, wideGapMin, ar.SmallestDesirableWidth, desc)
		}

	}

	// switch back
	if stepdown && chosen != nil {
		chosen.Edge[0], chosen.Edge[2] = chosen.Edge[2], chosen.Edge[0]
		chosen.Edge[1] = -chosen.Edge[1]
		x1, x2 = x2, x1
		y = -y
		if direction == 0 {
			direction = 1
		} else {
			direction = 0
		}
	}

	return

}
