package upload

import (
	"fmt"
	"io"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

var baseDirX = "c:\\TEMP\\"
var baseDirY = "c:/TEMP/"