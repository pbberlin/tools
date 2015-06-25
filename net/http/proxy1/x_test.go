package proxy1

import (
	"net/http"
	"testing"
)

func Test2(t *testing.T) {
	http.ListenAndServe("localhost:4000", nil)
}
