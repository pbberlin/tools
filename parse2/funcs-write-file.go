package parse2

import (
	"bytes"
	"io/ioutil"
	"log"

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
	err := ioutil.WriteFile(fn, b, 0)
	if err != nil {
		log.Println(err)
	}
}
