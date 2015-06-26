package osmaps

import "fmt"

func insertLinkedList(nd *node, fct1 func(string, string) bool, key, value string) bool {

	inserted := false

	nd.kv = append(nd.kv, kvt{key, value, noSuccessor})

	idxLastEl := len(nd.kv) - 1
	isNewMax := false

	fmt.Printf("added %v as %v\n", key, idxLastEl)

	// first insert
	if len(nd.kv) == 1 {
		nd.min = key
		nd.max = key
		nd.minIdx = 0
		report(nd)
		return true
	}

	// new min
	if fct1(key, nd.min) {

		fmt.Printf("\t\tfound %v<%v  - newMin \n", key, nd.min)
		nd.kv[idxLastEl].succ = nd.minIdx
		nd.min = key
		nd.minIdx = idxLastEl

	} else {

		// in between or new max

		ec := 0 // emergency counter
		i := nd.minIdx
		iPrev := nd.minIdx
		for {

			ec++
			if ec > cFanout+1 {
				fmt.Printf("overflow in iteration linked list; %v\n", ec)
				report(nd)
				return false
			}

			// in between
			fmt.Printf("\t\texamine %v<%v<%v  - iter%v\n", nd.kv[iPrev].key, key, nd.kv[i].key, ec)

			if fct1(nd.kv[iPrev].key, key) && fct1(key, nd.kv[i].key) {
				fmt.Printf("\t\t TRUE \n")
				prev := nd.kv[iPrev].succ
				nd.kv[iPrev].succ = idxLastEl
				nd.kv[idxLastEl].succ = prev
				break

			}

			// new max
			if nd.kv[i].succ == noSuccessor {
				isNewMax = true
				fmt.Printf("\t\tno successor, new max %v\n", key)
				break
			}

			//
			iPrev = i
			i = nd.kv[i].succ

		}

		if isNewMax {
			nd.kv[i].succ = idxLastEl
			nd.max = key
		}

	}

	pointerHash(&key)
	inserted = true

	report(nd)

	return inserted
}

func report(nd *node) {
	fmt.Printf("\t\t\t")
	for i := 0; i < len(nd.kv); i++ {
		fmt.Printf("%v %v  | ", nd.kv[i].key, nd.kv[i].succ)
	}
	fmt.Printf("\n")

	fmt.Printf("\t\tmin - max: %v %v \n", nd.min, nd.max)
	fmt.Printf("\t\tminIdx-vl: %v %v \n", nd.minIdx, nd.kv[nd.minIdx].key)
	sortedIter(nd)
	fmt.Printf("\n")

}

func sortedIter(nd *node) {

	fmt.Printf("\t\t")
	i := nd.minIdx
	ec := 0 // emergency counter
	for {
		ec++
		if ec > cFanout+1 {
			fmt.Printf("overflow in iteration linked list; %v\n", ec)
			break
		}

		fmt.Printf("%v ", nd.kv[i].key)
		if nd.kv[i].succ == noSuccessor {
			break
		}
		i = nd.kv[i].succ
	}
	fmt.Printf("\n")

}
