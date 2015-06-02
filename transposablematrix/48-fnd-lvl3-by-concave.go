package transposablematrix

//
// Find the perfect fit for a given edge x-y-x
func StairyPerfectConcave(ar *Reservoir,
	fs Fusion) (chosen *Amorph, baseShift Point) {

	var x1, y, x2 = fs.xyx[0], fs.xyx[1], fs.xyx[2]

	var stepdown bool
	if y < 0 {
		stepdown = true
		x1, x2 = x2, x1
		y = -y
	}

	_, chosen = ar.ByStairyEdge(fs, x1, y, x2, 0, westw, grow)

	// switch back
	if stepdown && chosen != nil {
		chosen.Edge[0], chosen.Edge[2] = chosen.Edge[2], chosen.Edge[0]
		chosen.Edge[1] = -chosen.Edge[1]
		x1, x2 = x2, x1
		y = -y
	}

	return chosen, Point{}

}

// We need to distinguish four cases
// 		stepup				or	stepdown
//		LowStat-HighDyn		or	LowDyna-HighSta-
//
// First, recollect following equivalence:
//
//  XXXXX-<-                ->-XXXXX
//  XXX          equiv           XXX
//  XXX          mirror          XXX
//  3;2;2...               ...2;-2;3
//  westw                      eastw
//
// Then consider switching of the shrink *direction* too:
//
//  ->-XXXXX                  XXXXX-<-
//  ->-X          equiv           X-<-
//  ->-X          mirror          X-<-
//  3..1;2;4                  4;-2;3..1
//  eastw                     westw
//
func StairyShrinkyConcave(ar *Reservoir,
	fs Fusion) (chosen *Amorph, baseShift Point) {

	if fs.xyx[1] > 0 {
		chosen, baseShift = ar.AdapterStepdownAndDirection(fs, "LowStat-HighDyn-", eastw)
		if chosen == nil {
			chosen, baseShift = ar.AdapterStepdownAndDirection(fs, "LowDyna-HighSta-", westw)
		}
	} else {
		chosen, baseShift = ar.AdapterStepdownAndDirection(fs, "LowStat-HighDyn-", westw)
		if chosen == nil {
			chosen, baseShift = ar.AdapterStepdownAndDirection(fs, "LowDyna-HighSta-", eastw)
		}
	}

	return

}

func (ar *Reservoir) AdapterStepdownAndDirection(fs Fusion,
	desc string, direction VariDirection) (chosen *Amorph, baseShift Point) {

	var x1, y, x2 = fs.xyx[0], fs.xyx[1], fs.xyx[2]

	var stepdown bool
	if y < 0 {
		stepdown = true
		x1, x2 = x2, x1
		y = -y
		direction.SwitchHoriz()
	}

	if direction == westw {
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
			_, chosen = ar.ByStairyEdge(fs, xdyn, y, x2, maxOffs, westw, shrink)
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
			_, chosen = ar.ByStairyEdge(fs, x1, y, xdyn, maxOffs, eastw, shrink)
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
		direction.SwitchHoriz()
	}

	// if stepdown && chosen != nil {
	if chosen != nil {
		baseShift.x += fs.xyx[0] - chosen.Edge[0]
		// baseShift.y--
	}

	return

}
