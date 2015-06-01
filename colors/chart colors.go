package colors

import "fmt"

var GraphColors [][]uint8

// http://ksrowell.com/blog-visualizing-data/2012/02/02/optimal-colors-for-graphs/
var Colors = [][]uint8{
	[]uint8{114, 147, 203},
	[]uint8{225, 151, 76},
	[]uint8{132, 186, 91},
	[]uint8{211, 94, 96},

	[]uint8{128, 133, 133},
	[]uint8{144, 103, 157},
	[]uint8{171, 104, 87},
	[]uint8{204, 194, 15},

	[]uint8{57, 106, 177},
	[]uint8{218, 124, 48},
	[]uint8{62, 150, 81},
	[]uint8{204, 37, 41},

	[]uint8{83, 81, 84},
	[]uint8{107, 76, 154},
	[]uint8{146, 36, 40},
	[]uint8{148, 139, 61},
}

func init() {
	GraphColors = Colors

	// colorizer := Colorizer(1)
	// for i := 0; i < 5; i++ {
	// 	pf("%v \n", colorizer())
	// }

}

func PreventOverFlow(base uint8, summand int) uint8 {

	if int(base)+summand < 0 {
		return uint8(0)
	}
	if int(base)+summand > 255 {
		return uint8(255)
	}

	if summand > 0 {
		return base + uint8(summand)
	} else {
		return base - uint8(-summand)
	}

}

func AlternatingColorShades(idxColor int, idx int) string {

	idxCol := (idxColor + 4) % len(Colors)

	const step = 25
	const steps = 4
	variation := 0

	idx = idx % steps // revolve

	// alternate even - uneven
	if idx%2 == 0 {
		variation = (idx + 0) * step
	} else {
		variation = (idx + 2) * step // stronger distincting between alternating elements
	}

	// %x is the hex format, %2.2x makes padding zeros
	col := fmt.Sprintf("%2.2x%2.2x%2.2x",
		PreventOverFlow(Colors[idxCol][0], variation),
		PreventOverFlow(Colors[idxCol][1], variation),
		PreventOverFlow(Colors[idxCol][2], variation),
	)

	return col

}
