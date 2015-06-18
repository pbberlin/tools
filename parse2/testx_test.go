package parse2

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"testing"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/subsort"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

func main() {

	texts := []map[string][]byte{}

	//
	// ================================================
	for i := 3; i <= 5; i++ {
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

		texts = append(texts, mp)

		reIndent(doc, 0)

		computeOutline(doc, 0, []int{0})

		bytes2File(fn1, xPathDump)
		dom2File(fn3, doc)
	}

	sorted1 := subsort.SortMapByCount(attrDistinct)
	sorted1.Print()
	fmt.Println()
	sorted2 := subsort.SortMapByCount(nodeDistinct)
	sorted2.Print()

	pf("testing\n")
	for k1, v1 := range texts {
		pf(" %v\n", k1)
		for k2, v2 := range v1 {
			pf(" cmp  %v - %v\n  to ", pbstrings.Ellipsoider(string(v2), 10), k2)
			cols := 0
			for k3, v3 := range texts {
				if k1 == k3 {
					continue
				}
				for _, v4 := range v3 {
					pf(" %v |", pbstrings.ToLen(pbstrings.Ellipsoider(string(v4), 10), 20))
					cols++
					if cols%4 == 0 {
						pf("\n     ")
					}
				}
			}
			pf("\n")
		}
	}

	return

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}
