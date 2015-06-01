package transposablematrix

import "math/rand"

func (m *TransposableMatrix) SeedCenteredBlocks(deterministic bool) {
	sdx := []int{0, -6, -3, 0, 3, 6, 3, 0, -3}
	sdy := []int{0, 0, -2, -4, -2, 0, 2, 4, 2, 0}
	nobj := len(sdx)

	reverse := false
	if rand.Intn(2) > 0 {
		reverse = true
	}

	for kx := 0; kx < nobj; kx++ {

		k := kx
		if reverse && !deterministic {
			k = nobj - kx - 1
		}

		dx := 3 + rand.Intn(9)
		dy := 2 + rand.Intn(4)

		if deterministic {
			dx = 8
		}

		for i := sdx[k] - dx/2; i < sdx[k]-dx/2+dx; i++ {
			for j := sdy[k] - dy/2; j < sdy[k]-dy/2+dy; j++ {
				s := spf("%2v", kx)
				m.SetLabel(i, j, Slot{Label: s})
			}
		}

	}

}

func (m *TransposableMatrix) SeedRandomCenteredBlocks() {
	nobj := 5
	for k := 0; k < nobj; k++ {
		n := 5
		dev := 8
		x, y := -dev, -dev
		for i := 0; i < n; i++ {
			x += rand.Intn(2*dev+2) / n
			y += rand.Intn(2*dev+2) / n
		}
		dx := 2 + rand.Intn(5)
		dy := 2 + rand.Intn(5)
		// pf("%v;%v %v-%v\n", x, y, dx, dy)
		for i := x; i < x+dx; i++ {
			for j := y; j < y+dy; j++ {
				m.SetLabel(i, j, Slot{Label: " X"})
			}
		}
	}
}

func (m *TransposableMatrix) SeedA(x0, y0 int) {

	for x := 1; x <= 5; x++ {
		for y := 1; y < 7; y++ {
			m.SetLabel(x0+x, y0+y, Slot{Label: " X"})
		}
	}

	m.SetLabel(x0+2, y0+2, Slot{})
	m.SetLabel(x0+4, y0+2, Slot{})
	m.SetLabel(x0+2, y0+3, Slot{})
	m.SetLabel(x0+4, y0+3, Slot{})

	m.SetLabel(x0+2, y0+5, Slot{})
	m.SetLabel(x0+4, y0+5, Slot{})
	m.SetLabel(x0+2, y0+6, Slot{})
	m.SetLabel(x0+4, y0+6, Slot{})
}

func (m *TransposableMatrix) SeedArrow(x0, y0, perspective int) {

	for x := 0; x <= 4; x++ {
		m.SetLabel(x0+x, y0+2, Slot{Label: " X"})
	}

	m.SetLabel(x0+5, y0+0, Slot{Label: " X"})
	m.SetLabel(x0+5, y0+1, Slot{Label: " X"})
	m.SetLabel(x0+5, y0+2, Slot{Label: " X"})
	m.SetLabel(x0+5, y0+3, Slot{Label: " X"})
	m.SetLabel(x0+5, y0+4, Slot{Label: " X"})

	m.SetLabel(x0+6, y0+1, Slot{Label: " X"})
	m.SetLabel(x0+6, y0+2, Slot{Label: " X"})
	m.SetLabel(x0+6, y0+3, Slot{Label: " X"})

	m.SetLabel(x0+7, y0+2, Slot{Label: " X"})

}
