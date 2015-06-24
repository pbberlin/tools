package parse2

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pbberlin/tools/pblog"
	"golang.org/x/net/html"
)

func dom2File(fn string, node *html.Node) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	if err != nil {
		log.Fatal(err)
	}
	bytes2File(fn, b.Bytes())
}

func bytes2File(fn string, b []byte) {
	var err error
	err = ioutil.WriteFile(fn, b, 0)
	if err != nil {

		err = os.MkdirAll(filepath.Dir(fn), os.ModePerm)
		pblog.LogE(err, "directory creation failed: %v")

		err = ioutil.WriteFile(fn, b, 0)
		pblog.LogE(err)

	}
}

func bytesFromFile(fn string) []byte {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println(err)
	}
	return b
}
