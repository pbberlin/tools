package util

import "fmt"
import "net/http"

var Counter chan int = make(chan int)
var cntr = -1

/*
	a counter - globally synchronized across all goroutines
*/

func init() {
	go func() {
		for {
			cntr++
			Counter <- cntr
		}
	}()
	http.HandleFunc("/counter/demo", counterDemo)
	http.HandleFunc("/counter/reset", counterReset)
	http.HandleFunc("/counter/decr", counterDecrement)
}

func CounterLast() int {
	return cntr - 1
}

func counterReset(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "cntr was %v<br>\n", cntr)
	cntr = 0
	fmt.Fprintf(w, "cntr resetted to %v<br>\n", cntr)
}
func counterDecrement(w http.ResponseWriter, r *http.Request) {
	cntr--
	fmt.Fprintf(w, "cntr decremented to %v<br>\n", cntr)
}

func counterDemo(w http.ResponseWriter, r *http.Request) {
	x := <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)

	x = <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)

	x = <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)
}
