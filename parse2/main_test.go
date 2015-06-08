package parse2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/pbberlin/tools/fetch"
	"github.com/pbberlin/tools/subsort"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

func main() {

	tests := make([]string, 2)
	var err error

	tests[0] = `<!DOCTYPE html><html><head>
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

	tests[1] = `	<p>
				Ja so sans<br/>
				Ja die sans.
			</p>
			<ul>
				<li>die ooolten Rittersleut</li>
			</ul>`

	for i := 0; i < len(tests); i++ {
		var doc1 *html.Node
		fn1 := fmt.Sprintf("outpT%v.txt", i)
		fn2 := fmt.Sprintf("outpT%v.html", i)
		doc1, err = html.Parse(strings.NewReader(tests[i]))
		if err != nil {
			log.Fatal(err)
		}
		TraverseVertConvert(doc1, 0)
		TraverseVertIndent(doc1, 0)
		TravVertStats(doc1, 0)
		ioutil.WriteFile(fn1, xPathDump, 0)
		dom2File(doc1, fn2)
	}

	// ================================================
	for i := 1; i <= 3; i++ {
		var doc *html.Node
		url := fmt.Sprintf("http://localhost:4000/static/handelsblatt.com/article0%v.html", i)
		fn1 := fmt.Sprintf("outpL%v.txt", i)
		fn2 := fmt.Sprintf("outpL%v.html", i)
		_, resBytes, err := fetch.UrlGetter(url, nil, true)
		resBytes = globFixes(resBytes)
		doc, err = html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}

		TraverseVertConvert(doc, 0)

		for i := 0; i < 5; i++ {
			TravVertConvertEmptyLeafs(doc, 0)
			TravHoriRemoveCommentAndSpaces(Tx{doc, 0})
		}

		for i := 0; i < 88; i++ {
			TraverseVert_ConvertDivDiv(doc, 0)
			TravHoriRemoveCommentAndSpaces(Tx{doc, 0})
		}

		TravVertStats(doc, 0)

		TraverseVertIndent(doc, 0)

		ioutil.WriteFile(fn1, xPathDump, 0)
		dom2File(doc, fn2)
	}

	sorted1 := subsort.SortMapByCount(attrDistinct)
	sorted1.Print()
	fmt.Println()
	sorted2 := subsort.SortMapByCount(nodeDistinct)
	sorted2.Print()

	return

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

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
