package parse2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/pbberlin/tools/fetch"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

func main() {

	s1 := `<!DOCTYPE html><html><head>
		<script type="text/javascript" src="./article01_files/empty.js"></script>
		<link href="./article01_files/vendor.css" rel="stylesheet" type="text/css"/>
		</head><body><p>Links:
				<span>p1</span>
				<span>p2</span>
				<span>p3</span>
			</p>
			<style> p {font-size:17px}</style>
			<ul>
				<li id='332' ><a   href="foo">Linktext1 <span>inside</span></a>
				<li><a   href="/bar/baz">BarBaz</a>
			</ul></body></html>`

	s2 := `	<p>
				Ja so sans<br/>
				Ja die sans.
			</p>
			<ul>
				<li>die ooolten Rittersleut</li>
			</ul>`

	var doc1, doc2, doc3 *html.Node
	_, _, _ = doc1, doc2, doc3
	var err error

	doc1, err = html.Parse(strings.NewReader(s1))
	if err != nil {
		log.Fatal(err)
	}

	TraverseVertCleanse(doc1, 0)
	TraverseVertIndent(doc1, 0)
	TraverseVert(doc1, 0)
	ioutil.WriteFile("outp1.txt", xPathDump, 0)
	dom2File(doc1, "outp1.html")

	// ================================================
	doc2, err = html.Parse(strings.NewReader(s2))
	if err != nil {
		log.Fatal(err)
	}
	TraverseVertCleanse(doc2, 0)
	TraverseVertIndent(doc2, 0)
	TraverseVert(doc2, 0)
	ioutil.WriteFile("outp2.txt", xPathDump, 0)
	dom2File(doc2, "outp2.html")

	// ================================================
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("http://localhost:4000/static/handelsblatt.com/article0%v.html", i)
		fn1 := fmt.Sprintf("outpL%v.txt", i)
		fn2 := fmt.Sprintf("outpL%v.html", i)
		_, resBytes, err := fetch.UrlGetter(url, nil, true)
		resBytes = globFixes(resBytes)
		doc3, err = html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}
		TraverseVertCleanse(doc3, 0)
		TraverseHoriRemoveNodes(Tx{doc3, 0})
		TraverseVertIndent(doc3, 0)
		TraverseVert(doc3, 0)
		ioutil.WriteFile(fn1, xPathDump, 0)
		dom2File(doc3, fn2)
	}

	//
	for k, val := range attrDistinct {
		fmt.Printf("%12v %v\n", k, val)
	}
}

func globFixes(b []byte) []byte {
	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}

func dom2File(node *html.Node, fn string) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(fn, b.Bytes(), 0)

}
