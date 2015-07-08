// Package osutilpb reads and writes files.
package osutilpb

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pbberlin/tools/logif"
	"golang.org/x/net/html"
)

func Dom2File(fn string, node *html.Node) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	logif.F(err)
	Bytes2File(fn, b.Bytes())
}

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

func BytesFromFile(fn string) []byte {
	b, err := ioutil.ReadFile(fn)
	logif.E(err)
	return b
}
