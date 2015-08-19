package util

import (
	"os"
	"os/exec"
	"runtime"
)

// better use
// go get -u github.com/nsf/termbox-go
func ClearScreen() {

	if runtime.GOOS != "windows" {
		panic("this command only runs on windows")
	}

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

}
