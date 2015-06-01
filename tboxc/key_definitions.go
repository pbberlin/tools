package tboxc

import "github.com/pbberlin/termbox-go"

var fcmap = []string{
	"CTRL2~0",
	"CTRL+A",
	"CTRL+B",
	"CTRL+C",
	"CTRL+D",
	"CTRL+E",
	"CTRL+F",
	"CTRL+G",
	"CTRL+H, BACKSPACE",
	"CTRL+I, TAB",
	"CTRL+J",
	"CTRL+K",
	"CTRL+L",
	"CTRL+M, ENTER",
	"CTRL+N",
	"CTRL+O",
	"CTRL+P",
	"CTRL+Q",
	"CTRL+R",
	"CTRL+S",
	"CTRL+T",
	"CTRL+U",
	"CTRL+V",
	"CTRL+W",
	"CTRL+X",
	"CTRL+Y",
	"CTRL+Z",
	"CTRL+3, ESC, CTRL+[",
	"CTRL+4, CTRL+\\",
	"CTRL+5, CTRL+]",
	"CTRL+6",
	"CTRL+7, CTRL+/, CTRL+_",
	"SPACE",
}

var fkmap = []string{
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"F8",
	"F9",
	"F10",
	"F11",
	"F12",
	"INSERT",
	"DELETE",
	"HOME",
	"END",
	"PGUP",
	"PGDN",
	"ARROW UP",
	"ARROW DOWN",
	"ARROW LEFT",
	"ARROW RIGHT",
}

func FuncKeyMap(k termbox.Key) string {
	if k == termbox.KeyCtrl8 {
		return "CTRL+8, BACKSPACE 2" /* 0x7F */
	} else if k >= termbox.KeyArrowRight && k <= 0xFFFF {
		return fkmap[0xFFFF-k]
	} else if k <= termbox.KeySpace {
		return fcmap[k]
	}
	return "UNKNOWN"
}