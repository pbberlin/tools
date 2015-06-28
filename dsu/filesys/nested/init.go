package nested

import (
	"fmt"

	"github.com/pbberlin/tools/dsu/filesys"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf

var spf func(format string, a ...interface{}) string = fmt.Sprintf

var t string

func init() {
	fo := filesys.Directory{}
	t = fmt.Sprintf("%T", fo) // "kind"
	t = "nested"
}
