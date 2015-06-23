// +build parsing
// go test -tags=parsing

package parse2

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/subsort"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	main()
}

var numTotal = 0 // comparable html docs
const stageMax = 2

func main() {

	//
	// ================================================
	iter := []int{1, 2, 3}

	for _, i := range iter {
		var doc *html.Node
		url := fmt.Sprintf("http://localhost:4000/static/handelsblatt.com/art0%v.html", i)
		fn1 := fmt.Sprintf("outp_%03v_xpath.txt", i)
		fn2 := fmt.Sprintf("outp_%03v_texts.txt", i)
		fn3, fnKey := getFN(i, 0)
		resBytes, err := pbfetch.UrlGetter(url, nil, false)
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

		pf("-----------------\n")

		maxLvlPrev := 0
		for i := 0; i < 1; i++ {
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

		textsBytes, textsSorted := orderByOutline(textsByOutl)
		bytes2File(fn2, textsBytes)
		textsByArticOutl[fnKey] = textsSorted

		reIndent(doc, 0)

		bytes2File(fn1, xPathDump)
		dom2File(fn3, doc)

		numTotal++
	}

	// statistics on elements and attributes
	sorted1 := subsort.SortMapByCount(attrDistinct)
	sorted1.Print()
	fmt.Println()
	sorted2 := subsort.SortMapByCount(nodeDistinct)
	sorted2.Print()

	for weedStage := 1; weedStage <= stageMax; weedStage++ {

		levelsToProcess = map[int]bool{weedStage: true}
		frags := rangeOverTexts()

		similaritiesToFile(frags, weedStage)

		weedoutMap := map[string]map[string]bool{}
		for _, i := range iter {
			_, fnKey := getFN(i, weedStage)
			weedoutMap[fnKey] = map[string]bool{}
		}
		weedoutMap = assembleWeedout(frags, weedoutMap)

		bb := pbstrings.IndentedDumpBytes(weedoutMap)
		bytes2File(spf("outp_wd_%v.txt", weedStage), bb)

		for _, i := range iter {
			fnInn, _ := getFN(i, weedStage-1)
			fnOut, fnKey := getFN(i, weedStage)

			resBytes := bytesFromFile(fnInn)
			doc, err := html.Parse(bytes.NewReader(resBytes))
			if err != nil {
				log.Fatal(err)
			}
			weedoutApply(weedoutMap[fnKey], doc)
			dom2File(fnOut, doc)
		}
	}

	for _, i := range iter {
		fnInn, _ := getFN(i, stageMax)
		fnOut, _ := getFN(i, stageMax+1)

		resBytes := bytesFromFile(fnInn)
		doc, err := html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}
		flattenTraverse(doc)
		dom2File(fnOut, doc)
	}

	pf("correct finish\n")

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}

func getFN(articleId, weedoutStage int) (string, string) {
	fn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
	prefix := fmt.Sprintf("outp_%03v", articleId)
	return fn, prefix
}
