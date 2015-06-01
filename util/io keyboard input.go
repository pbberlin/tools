package util

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"strings"
	"time"
	// "syscall" disabled for gae
)

var StdInputOutsideTest *os.File // interactive testing impossible - stdin simply does not work under testing

func init() {
	// disabled for gae
	// StdInputOutsideTest = os.NewFile(uintptr(syscall.Stdin), "/dev/stdin") // tried 4 others - none works during testing
}


// https://groups.google.com/forum/#!topic/golang-nuts/w6TTJpw9stA
// Sadly, all shown methods require an onerous "ENTER".
func getKeyboardInput() {

	ck := make(chan bool)

	reader := bufio.NewReader(StdInputOutsideTest)
	go func() {
		for {
			fmt.Print("tell! ")
			// line, err := reader.ReadString('\n')
			oneRune, size, err := reader.ReadRune()
			line := string(oneRune)
			if err != nil {
				if err == io.EOF {
					fmt.Print("EOF ")
				} else {
					fmt.Print("error neq EOF: ", err)
					break
				}
			} else {
				line = strings.TrimSpace(line)
				fmt.Printf("u said %q s%v\n", line, size)
				if line == "q" {
					break
				}
			}
			time.Sleep(1200 * time.Millisecond)
		}
		ck <- true
	}()

	<-ck

}
