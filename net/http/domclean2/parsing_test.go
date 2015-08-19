// +build parsing
// go test -tags=parsing

package domclean2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"appengine"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/fetch_rss"
	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/osutilpb"
	"github.com/pbberlin/tools/sort/sortmap"
	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

var numTotal = 0 // comparable html docs
const stageMax = 3

const cTestHostDev = "localhost:8085"
const cTestHostOwn = "localhost:63222"

var hostWithPref = cTestHostDev + fetch_rss.UriMountNameY

func prepare(t *testing.T, c appengine.Context) {

	serveFile := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		fs1 := fetch_rss.GetFS(c)
		fileserver.FsiFileServer(fs1, fetch_rss.UriMountNameY, w, r)
	}
	http.HandleFunc(fetch_rss.UriMountNameY, loghttp.Adapter(serveFile))

	go func() {
		log.Fatal(
			http.ListenAndServe(cTestHostOwn, nil),
		)
	}()

}

func Test1(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	// c, err := aetest.NewContext(nil)
	// if err != nil {
	// 	lge(err)
	// 	t.Fatal(err)
	// }
	// defer c.Close()

	// prepare(t, c)

	lg("waiting for webserver")
	time.Sleep(2 * time.Millisecond)

	expandingPath := "www.welt.de"

	dirs1, _, msg, err := fileserver.GetDirContents(hostWithPref, expandingPath)
	if err != nil {
		lge(err)
		lg("%s", msg)
	}

	lg("dirs1 %v", stringspb.IndentedDump(dirs1))

	least3Files := []string{}
	for _, v1 := range dirs1 {

		dirs2, fils2, msg, err := fileserver.GetDirContents(hostWithPref, path.Join(expandingPath, v1))
		_ = dirs2
		if err != nil {
			lge(err)
			lg("%s", msg)
		}
		// lg("  dirs2 %v", stringspb.IndentedDump(dirs2))
		// lg("  fils2 %v", stringspb.IndentedDump(fils2))

		if len(fils2) > 2 {
			for i2, v2 := range fils2 {
				least3Files = append(least3Files, path.Join(expandingPath, v1, v2))
				if i2 == 2 {
					break
				}
			}
			break
		}
	}

	if len(least3Files) < 3 {
		lg("not enough files in rss fetcher cache")
		return
	}

	lg("  fils2 %v", stringspb.IndentedDump(least3Files))

	logdir := osutilpb.DirOfExecutable()
	logdir = filepath.Join(".", "outp")
	lg("logdir is %v ", logdir)
	err = os.Mkdir(logdir, 0755)
	if err != nil && !os.IsExist(err) {
		lge(err)
		return
	}
	err = os.RemoveAll(logdir)
	if err != nil {
		lge(err)
		return
	}

	iter := []int{0, 1, 2}

	weedoutFilename := func(articleId, weedoutStage int) (string, string) {
		stagedFn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
		prefix := fmt.Sprintf("outp_%03v", articleId)
		return stagedFn, prefix
	}

	for i, _ := range iter {

		var doc *html.Node
		url := spf("%v/%v", hostWithPref, least3Files[i])
		fn1 := spf("outp_%03v_xpath.txt", i)
		fn2 := spf("outp_%03v_texts.txt", i)
		fn3, fnKey := weedoutFilename(i, 0)

		resBytes, effUrl, err := fetch.UrlGetter(nil, fetch.Options{URL: url})
		if err != nil {
			lge(err)
			return
		}
		lg("fetched %4.1fkB from %v", float64(len(resBytes))/1024, stringspb.ToLenR(effUrl.String(), 60))
		continue

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
		osutilpb.Bytes2File(fn2, textsBytes)
		textsByArticOutl[fnKey] = textsSorted

		reIndent(doc, 0)

		osutilpb.Bytes2File(fn1, xPathDump)
		osutilpb.Dom2File(fn3, doc)

		numTotal++
	}

	return

	// statistics on elements and attributes
	sorted1 := sortmap.SortMapByCount(attrDistinct)
	sorted1.Print(6)
	fmt.Println()
	sorted2 := sortmap.SortMapByCount(nodeDistinct)
	sorted2.Print(6)

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
		osutilpb.Bytes2File(spf("outp_wd_%v.txt", weedStage), bb)

		for i, _ := range iter {
			fnInn, _ := weedoutFilename(i, weedStage-1)
			fnOut, fnKey := weedoutFilename(i, weedStage)

			resBytes := osutilpb.BytesFromFile(fnInn)
			doc, err := html.Parse(bytes.NewReader(resBytes))
			if err != nil {
				log.Fatal(err)
			}
			weedoutApply(weedoutMap[fnKey], doc)
			osutilpb.Dom2File(fnOut, doc)
		}
	}

	for i, _ := range iter {
		fnInn, _ := weedoutFilename(i, stageMax)
		fnOut, _ := weedoutFilename(i, stageMax+1)

		resBytes := osutilpb.BytesFromFile(fnInn)
		doc, err := html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			log.Fatal(err)
		}
		flattenTraverse(doc)
		osutilpb.Dom2File(fnOut, doc)
	}

	{
		const jsonStream = `{
		  	"\\": 11,
		  	"\\panorama\\aus-aller-welt": 1,
		  	"\\politik\\international": 3,
		  	"\\politik\\konjunktur\\nachrichten": 1,
		  	"\\sport\\fussball": 1,
		  	"\\unternehmen\\dienstleister\\werber-rat": 1,
		  	"\\unternehmen\\handel-konsumgueter": 1
		}`
		// type Message struct {
		// 	Name, Text string
		// }

		dec := json.NewDecoder(strings.NewReader(jsonStream))
		var mp map[string]int
		err := dec.Decode(&mp)
		logif.E(err)
		fmt.Printf("%v\n", mp)
	}

	pf("correct finish\n")

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}
