package extract

import (
	"fmt"
	"net/http"
	"testing"
)

func TestExtraction(t *testing.T) {

	r, _ := http.NewRequest("GET", "some_url", nil)
	u := User{}

	Extract(r, &u)

	fmt.Printf("%v \n", u)
}
