// Package osutilpb reads and writes files with maximum convenience.
package osutilpb

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pbberlin/tools/logif"
	"golang.org/x/net/html"
)

// Dom2File writes DOM to file
func Dom2File(fn string, node *html.Node) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	logif.F(err)
	Bytes2File(fn, b.Bytes())
}

// Bytes2File writes bytes; creates path if neccessary
// and logs any errors even to appengine log
func Bytes2File(fn string, b []byte) {
	var err error
	err = ioutil.WriteFile(fn, b, 0)
	if err != nil {

		err = os.MkdirAll(filepath.Dir(fn), os.ModePerm)
		logif.E(err, "directory creation failed: %v")

		err = ioutil.WriteFile(fn, b, 0)
		logif.E(err)

	}
}

// BytesFromFile reads bytes and logs any
// errors even to appengine log.
func BytesFromFile(fn string) []byte {
	b, err := ioutil.ReadFile(fn)
	logif.E(err)
	return b
}
