package util

import "fmt"
import "net/http"

var Counter chan int = make(chan int)
var cntr = 0

/*
	a counter - globally synchronized across all goroutines
*/

func init() {
	go func() {
		for {
			Counter <- cntr
			cntr++
		}
	}()
	http.HandleFunc("/counter/get", GetCounter)
	http.HandleFunc("/counter/reset", ResetCounter)
}

func ResetCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "cntr was %v<br>\n", cntr)
	cntr = 0
	fmt.Fprintf(w, "cntr resetted to %v<br>\n", cntr)

}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	x := <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)

	x = <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)

	x = <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)
}
