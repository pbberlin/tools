package util

import "fmt"
import "net/http"

var Counter chan int = make(chan int)

/*
	a counter - globally synchronized across all goroutines
*/
func init() {
	go func() {
		cntr := 0
		for {
			Counter <- cntr
			cntr++
		}
	}()
	http.HandleFunc("/counter", GetCounter)
}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	x := <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)

	x = <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)

	x = <-Counter
	fmt.Fprintf(w, "cntr %v\n", x)
}
