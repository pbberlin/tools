package transposablematrix

import "github.com/pbberlin/termbox-go"

var TermBoxDone = make(chan struct{}) // exported to be used by main()

func InitTB() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)

	paintMain()
	termbox.Flush()

	inpMode := 0
	ctrlSequOpened := false // control sequence opened

	// now open up the keyboard input loop
	// decouple it from init(), since init()s must return immediately
	go func() {
		defer termbox.Close()
		defer close(TermBoxDone)

	lblMainLoop:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ctrlSequOpened && ev.Key == termbox.KeyCtrlQ || string(ev.Ch) == "q" {
					break lblMainLoop
				}
				if ctrlSequOpened && ev.Key == termbox.KeyCtrlS {
					termbox.Sync()
				}
				if ctrlSequOpened && ev.Key == termbox.KeyCtrlM {
					chmap := []termbox.InputMode{
						termbox.InputEsc | termbox.InputMouse,
						termbox.InputAlt | termbox.InputMouse,
						termbox.InputEsc,
						termbox.InputAlt,
					}
					inpMode++
					if inpMode >= len(chmap) {
						inpMode = 0
					}
					termbox.SetInputMode(chmap[inpMode])
				}

				if ev.Key == termbox.KeyPgup || ev.Key == termbox.KeyArrowUp {
					currStage--
					painAppStageData()
					printKeyPress(&ev)
					termbox.Flush()
					break
				}
				if ev.Key == termbox.KeyPgdn || ev.Key == termbox.KeyArrowDown {
					currStage++
					painAppStageData()
					printKeyPress(&ev)
					termbox.Flush()
					break
				}

				if ev.Key == termbox.KeyCtrlX {
					ctrlSequOpened = true
				} else {
					ctrlSequOpened = false
				}

				paintMain()
				printKeyPress(&ev)
				termbox.Flush()
			case termbox.EventResize:
				paintMain()
				printResizeEvent(&ev)
				termbox.Flush()
			case termbox.EventMouse:
				paintMain()
				printMouseEvent(&ev)
				crossHair(ev.MouseX, ev.MouseY)
				termbox.Flush()
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}()
}
