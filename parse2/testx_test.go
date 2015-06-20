package parse2

import (
	"bytes"
	"fmt"
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

var articleTexts = map[string]map[string][]byte{}

func main() {

	//
	// ================================================
	// for i := 3; i <= 5; i++ {
	for i := 1; i <= 3; i++ {
		var doc *html.Node
		url := fmt.Sprintf("http://localhost:4000/static/handelsblatt.com/article0%v.html", i)
		fn1 := fmt.Sprintf("outpI%v_1S.txt", i)
		fn2sh := fmt.Sprintf("outpI%v_2Tsh.txt", i)
		fn2lg := fmt.Sprintf("outpI%v_2Tlg.txt", i)
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

		computeOutline(doc, 0, []int{0})
		nodeCountHoriz(doc, 0, 1)

		textExtraction(doc, 0)

		textsLong := sortByKey(mpLg)
		bytes2File(fn2lg, textsLong)

		textsShrt := sortByKey(mpSh)
		bytes2File(fn2sh, textsShrt)

		articleTexts[fn3] = mpSh

		reIndent(doc, 0)

		bytes2File(fn1, xPathDump)
		dom2File(fn3, doc)
	}

	rangeOverTexts()

	// bfrags := pbstrings.IndentedDumpBytes(frags)
	bfrags := []byte{}
	for _, v := range frags {
		bfrags = append(bfrags, v.ArticleUrl...)
		bfrags = append(bfrags, ' ')
		bfrags = append(bfrags, fmt.Sprintf("%v", v.Lvl)...)
		bfrags = append(bfrags, ' ')
		bfrags = append(bfrags, string(v.Outline)...)
		bfrags = append(bfrags, '\n')
		bfrags = append(bfrags, string(v.Text)...)
		bfrags = append(bfrags, '\n')
		for _, v1 := range v.Similars {
			bfrags = append(bfrags, string(v1)...)
			bfrags = append(bfrags, '\n')
		}
	}
	bytes2File("outp_frags.txt", bfrags)

	sorted1 := subsort.SortMapByCount(attrDistinct)
	sorted1.Print()
	fmt.Println()
	sorted2 := subsort.SortMapByCount(nodeDistinct)
	sorted2.Print()

}

func sortByKey(mp1which map[string][]byte) []byte {

	keys := make([]string, 0, len(mp1which))
	for k := range mp1which {
		keys = append(keys, k)
	}

	// sort.Strings(keys)
	sort.Sort(sortByOutline(keys))

	ret := []byte{}
	for _, key := range keys {
		row := fmt.Sprintf("%-12v: %s\n", key, mp1which[key])
		ret = append(ret, row...)
	}
	return ret

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}
