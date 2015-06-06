package parse2

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

func main() {

	s1 := `	<p>Links:
				<span>p1</span>
				<span>p2</span>
			</p>
			<ul>
				<li id='332' ><a   href="foo">Linktext1 <span>inside</span></a>
				<li><a   href="/bar/baz">BarBaz</a>
			</ul>`

	s2 := `	<p>
				Ja so sans<br/>
				Ja die sans.
			</p>
			<ul>
				<li>die ooolten Rittersleut</li>
			</ul>`

	var doc1, doc2 *html.Node
	_, _ = doc1, doc2
	var err error

	doc1, err = html.Parse(strings.NewReader(s1))
	if err != nil {
		log.Fatal(err)
	}

	doc2, err = html.Parse(strings.NewReader(s2))
	if err != nil {
		log.Fatal(err)
	}

	// _, resBytes, err := fetch.UrlGetter("http://localhost:4000/static/handelsblatt.com/article01.html", nil, true)
	// doc3, err = html.Parse(bytes.NewReader(resBytes))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	TraverseVert(doc1, 0)
	TraverseVert(doc2, 0)

	// TraverseHori(Tx{doc1, 0})

	//
	ioutil.WriteFile("outp.txt", xPathDump, 0)

	dom2File(doc1, "outp1.html")
	dom2File(doc2, "outp2.html")

}
