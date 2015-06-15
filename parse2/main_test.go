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

	tests := make([]string, 2)

	tests[0] = `<!DOCTYPE html><html><head>
		<script type="text/javascript" src="./article01_files/empty.js"></script>
		<link href="./article01_files/vendor.css" rel="stylesheet" type="text/css"/>
		</head><body><p>Links:
				<span>span01</span>
				<span>span02-line1<br>span02-line2</span>
				<span>span03</span>
			</p>
			<style> p {font-size:17px}</style>
			<ul>
				<li id='332' ><a   href="/some/first/page.html">Linktext1 <span>inside</span></a>
				<li><a   href="/snd/page" title="wi-title">LinkT2</a>
			</ul>
			<div>
				<div>div-1-content</div>
				<div>div-2-content</div>
				<p>pararaph in between</p>
				<div>div-3-content with iimmage<img alt="alt-cnt" title='title-cnt' 
				href='some-long-href-some-long-href-some-long-href-some-long-href'>after img</div>
			</div>
			</body></html>`

	tests[1] = `	<p>
				Ja so sans<br/>
				Ja die sans.
			</p>
			<ul>
				<li>die ooolten Rittersleut</li>
			</ul>`

	for i := 0; i < len(tests); i++ {
		fn := fmt.Sprintf(docRoot+"/handelsblatt.com/article%02v.html", i+4)
		ioutil.WriteFile(fn, []byte(tests[i]), 0)
	}

	//
	// ================================================
	for i := 4; i <= 4; i++ {
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

		TraverseVertConvert(doc, 0)

		for i := 0; i < 6; i++ {
			TravVertConvertEmptyLeafs(doc, 0)
			TravHoriRemoveCommentAndSpaces(Tx{doc, 0})
		}

		maxLvlPrev := 0
		for i := 0; i < 48; i++ {
			lpMax := TravVertMaxLevel(doc, 0)
			if lpMax != maxLvlPrev {
				fmt.Printf("i%2v: maxL %2v\n", i, lpMax)
				maxLvlPrev = lpMax
			}
			TraverseVert_CondenseDivStaples(doc, 0)
		}

		TravVertStats(doc, 0)

		TravVertTextify(doc, 0, 0)

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
		ioutil.WriteFile(fn2, mpb, 0)

		TraverseVertIndent(doc, 0)

		ioutil.WriteFile(fn1, xPathDump, 0)
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
	ioutil.WriteFile(fn, b.Bytes(), 0)

}
