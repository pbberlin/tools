package parse2

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/subsort"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

var articleTexts = map[string]map[string][]byte{}
var numTotal = 0

func main() {

	//
	// ================================================
	iter := []int{1, 2, 3}

	for _, i := range iter {
		var doc *html.Node
		url := fmt.Sprintf("http://localhost:4000/static/handelsblatt.com/art0%v.html", i)
		fn1 := fmt.Sprintf("outp_%v_1S.txt", i)
		fn2sh := fmt.Sprintf("outp_%v_2Tsh.txt", i)
		fn2lg := fmt.Sprintf("outp_%v_2Tlg.txt", i)
		fn3 := fmt.Sprintf("outp_%v_3.html", i)
		_, resBytes, err := pbfetch.UrlGetter(url, nil, true)
		if err != nil {
			log.Fatal(err)
		}
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

		computeOutline(doc, 0, []int{0})
		nodeCountHoriz(doc, 0, 1)

		textExtraction(doc, 0)

		textsLong := recreateOrderedByOutline(mpLg)
		bytes2File(fn2lg, textsLong)

		textsShrt := recreateOrderedByOutline(mpSh)
		bytes2File(fn2sh, textsShrt)

		articleTexts[fn3] = mpSh

		reIndent(doc, 0)

		bytes2File(fn1, xPathDump)
		dom2File(fn3, doc)

		numTotal++
	}

	sorted1 := subsort.SortMapByCount(attrDistinct)
	sorted1.Print()
	fmt.Println()
	sorted2 := subsort.SortMapByCount(nodeDistinct)
	sorted2.Print()

	levelsTolerance = 0
	for stage := 1; stage < 4; stage++ {

		levelsToProcess = map[int]bool{stage: true}
		frags := rangeOverTexts()
		similaritiesToFile(frags, stage)
		weedouts := weedOut(frags)

		for _, i := range iter {
			fnInn := fmt.Sprintf("outp_%v_%v.html", i, stage+2)
			fnOut := fmt.Sprintf("outp_%v_%v.html", i, stage+3)

			resBytes := bytesFromFile(fnInn)
			doc, err := html.Parse(bytes.NewReader(resBytes))
			if err != nil {
				log.Fatal(err)
			}
			weedoutApply(weedouts, doc)
			dom2File(fnOut, doc)
		}
	}

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}
