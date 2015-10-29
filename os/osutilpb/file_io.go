package osutilpb

/*
	use fsi/common/convenience funcs instead
*/

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pbberlin/tools/net/http/loghttp"
	"golang.org/x/net/html"
)

func WriteBytesToFilename(filename string, ptrB *bytes.Buffer) {
	Bytes2File(filename, ptrB.Bytes())
}

// Dom2File writes DOM to file
func Dom2File(fn string, node *html.Node) {
	lg, _ := loghttp.BuffLoggerUniversal(nil, nil)

	var b bytes.Buffer
	err := html.Render(&b, node)
	lg(err)
	Bytes2File(fn, b.Bytes())
}

// Bytes2File writes bytes; creates path if neccessary
// and logs any errors even to appengine log
func Bytes2File(fn string, b []byte) {
	lg, _ := loghttp.BuffLoggerUniversal(nil, nil)

	var err error
	err = ioutil.WriteFile(fn, b, 0)
	if err != nil {

		err = os.MkdirAll(filepath.Dir(fn), os.ModePerm)
		lg(err, "directory creation failed: %v")

		err = ioutil.WriteFile(fn, b, 0)
		lg(err)

	}
}

// BytesFromFile reads bytes and logs any
// errors even to appengine log.
func BytesFromFile(fn string) []byte {
	lg, _ := loghttp.BuffLoggerUniversal(nil, nil)
	b, err := ioutil.ReadFile(fn)
	lg(err)
	return b
}
