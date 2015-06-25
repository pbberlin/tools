// +build parsing
// go test -tags=parsing

package domclean2

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/pbberlin/tools/io/ioutilpb"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/fetch_rss"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/subsort"
	"golang.org/x/net/html"
)

var baseUrl string
var subdirs []string

func init() {
	baseUrl, subdirs = fetch_rss.Serve()
}

func Test1(t *testing.T) {
	main()
}

var numTotal = 0 // comparable html docs
const stageMax = 3

func main() {

	//
	// ================================================
	// iter := make([]int, len(testDocs))
	iter := append([]int{0, 1}, []int{2, 3, 4}...)

	for i, _ := range iter {
		var doc *html.Node
		url := fmt.Sprintf("%v/%v/art0%v.html", baseUrl, subdirs[0], i)
		fn1 := fmt.Sprintf("outp_%03v_xpath.txt", i)
		fn2 := fmt.Sprintf("outp_%03v_texts.txt", i)
		fn3, fnKey := weedoutFilename(i, 0)
		resBytes, err := fetch.UrlGetter(url, nil, false)
		if err != nil {
			log.Fatal(err)
		}
		resBytes = globFixes(resBytes)
		doc, err = html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}

		cleanseDom(doc, 0)

		physicalNodeRemoval(NdX{doc, 0})

		convEmptyElementLeafs(doc, 0)

		condenseNestedDivs(doc, 0)

		dumpXPath(doc, 0)

		computeOutline(doc, 0, []int{0})
		nodeCountHoriz(doc, 0, 1)

		textExtraction(doc, 0)

		textsBytes, textsSorted := orderByOutline(textsByOutl)
		ioutilpb.Bytes2File(fn2, textsBytes)
		textsByArticOutl[fnKey] = textsSorted

		reIndent(doc, 0)

		ioutilpb.Bytes2File(fn1, xPathDump)
		ioutilpb.Dom2File(fn3, doc)

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
		for i, _ := range iter {
			_, fnKey := weedoutFilename(i, weedStage)
			weedoutMap[fnKey] = map[string]bool{}
		}
		weedoutMap = assembleWeedout(frags, weedoutMap)

		bb := stringspb.IndentedDumpBytes(weedoutMap)
		ioutilpb.Bytes2File(spf("outp_wd_%v.txt", weedStage), bb)

		for i, _ := range iter {
			fnInn, _ := weedoutFilename(i, weedStage-1)
			fnOut, fnKey := weedoutFilename(i, weedStage)

			resBytes := ioutilpb.BytesFromFile(fnInn)
			doc, err := html.Parse(bytes.NewReader(resBytes))
			if err != nil {
				log.Fatal(err)
			}
			weedoutApply(weedoutMap[fnKey], doc)
			ioutilpb.Dom2File(fnOut, doc)
		}
	}

	for i, _ := range iter {
		fnInn, _ := weedoutFilename(i, stageMax)
		fnOut, _ := weedoutFilename(i, stageMax+1)

		resBytes := ioutilpb.BytesFromFile(fnInn)
		doc, err := html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}
		flattenTraverse(doc)
		ioutilpb.Dom2File(fnOut, doc)
	}

	pf("correct finish\n")

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}