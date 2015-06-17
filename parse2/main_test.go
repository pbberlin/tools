package parse2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"testing"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/subsort"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

func main() {

	for i := 0; i < len(tests); i++ {
		fn := fmt.Sprintf(docRoot+"/handelsblatt.com/article%02v.html", i+4)
		bytes2File(fn, []byte(tests[i]))
	}

	//
	// ================================================
	for i := 1; i <= 4; i++ {
		var doc *html.Node
		url := fmt.Sprintf("http://localhost:4000/static/handelsblatt.com/article0%v.html", i)
		fn1 := fmt.Sprintf("outpI%v_1S.txt", i)
		fn2 := fmt.Sprintf("outpI%v_2T.txt", i)
		fn3 := fmt.Sprintf("outpI%v_3.html", i)
		_, resBytes, err := pbfetch.UrlGetter(url, nil, true)
		resBytes = globFixes(resBytes)
		doc, err = html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}

		cleanseDom(doc, 0)

		for i := 0; i < 6; i++ {
			convEmptyElementLeafs(doc, 0)
			physicalNodeRemoval(NdX{doc, 0})
		}

		maxLvlPrev := 0
		for i := 0; i < 48; i++ {
			lpMax := maxTreeDepth(doc, 0)
			if lpMax != maxLvlPrev {
				fmt.Printf("i%2v: maxL %2v\n", i, lpMax)
				maxLvlPrev = lpMax
			}
			condenseNestedDivs(doc, 0)
		}

		dumpXPath(doc, 0)

		textExtraction(doc, 0, 0)

		mpb := []byte{}
		keys := make([]string, 0, len(mp))
		for k := range mp {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if true || len(mp[key]) > 30 {
				row := fmt.Sprintf("%8v: %s\n", key, mp[key])
				mpb = append(mpb, row...)
			}
		}
		bytes2File(fn2, mpb)

		TraverseVertIndent(doc, 0)

		bytes2File(fn1, xPathDump)
		dom2File(fn3, doc)
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
