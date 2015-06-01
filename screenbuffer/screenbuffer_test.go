package screenbuffer

import "testing"

func TestBase2Map(t *testing.T) {

	w := 55
	h := 44
	m := NewTransposableMatrix(w, h)

	sX := []int{-6, 0, 8}
	sY := []int{-5, 0, 7}

	for i := 0; i < 4; i++ {
		perspective := i
		for i := 0; i < len(sX); i++ {
			x := sX[i]
			for i := 0; i < len(sY); i++ {
				y := sY[i]
				xb, yb := m.transposeMapped2Base(x, y, perspective) // Set() Get()

				xm, ym := m.transposeBase2Mapped(xb, yb, perspective) // Set() Get()
				if xm != x || ym != y {
					t.Errorf("%v: %3v %3v - %3v %3v - %3v %3v \n", perspective, x, y, xb, yb, xm, ym)
				}
			}

		}
		// pf("\n")
	}

}

func TestScreenbuffer(t *testing.T) {

	sb := new(ScreenBuffer)

	sb.PrintAppend("abc\nxyz")
	sb.PrintAppend("def")
	sb.PrintAppend("\n")
	sb.PrintAppend("\n")
	sb.PrintAppend("123")

	sb.PrintAt(10, "and let the merry breezes blow\nmy dust to where")
	sb.PrintAt(5, "some flowers\ngrow")
	sb.PrintAt(0, "      this is\n   my last\n         and final\n      will")
	sb.PrintAt(4, "         good luck to all of you")

	// sb.Dump()

	got := spf("%s", sb)
	want := `&[abc      this is xyzdef   my last          and final 123      will          good luck to all of you some flowers grow    and let the merry breezes blow my dust to where]`

	if got != want {
		t.Errorf("got  %v\nwant %v\n", got, want)
		t.Fail()
	}

}

func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func TestMsg(t *testing.T) {
	// pf("screenbuffer testing done...\n")
}
