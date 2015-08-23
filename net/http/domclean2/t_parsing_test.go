// +build parsing
// go test -tags=parsing

package domclean2

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"appengine/aetest"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/fetch_rss"
	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/osutilpb"
	"github.com/pbberlin/tools/sort/sortmap"
	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

const numTotal = 3 // comparable html docs
const stageMax = 3 // weedstages

const cTestHostDev = "localhost:8085"
const cTestHostOwn = "localhost:63222"

var hostWithPref = cTestHostDev + fetch_rss.UriMountNameY

func prepare(t *testing.T) aetest.Context {

	lg, lge := loghttp.Logger(nil, nil)
	_ = lg

	c, err := aetest.NewContext(nil)
	if err != nil {
		lge(err)
		t.Fatal(err)
	}

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

	return c

}

func Test1(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	// c := prepare(t)
	// defer c.Close()

	lg("waiting for webserver")
	time.Sleep(2 * time.Millisecond)

	remoteHostname := "www.welt.de"

	dirs1, _, msg, err := fileserver.GetDirContents(hostWithPref, remoteHostname)
	if err != nil {
		lge(err)
		lg("%s", msg)
	}

	lg("dirs1")
	for _, v := range dirs1 {
		lg("    %v", v)
	}

	least3Files := []string{}
	for _, v1 := range dirs1 {

		dirs2, fils2, msg, err := fileserver.GetDirContents(hostWithPref, path.Join(remoteHostname, v1))
		_ = dirs2
		if err != nil {
			lge(err)
			lg("%s", msg)
		}
		// lg("  dirs2 %v", stringspb.IndentedDump(dirs2))
		// lg("  fils2 %v", stringspb.IndentedDump(fils2))

		if len(fils2) > numTotal-1 {
			for i2, v2 := range fils2 {
				least3Files = append(least3Files, path.Join(remoteHostname, v1, v2))
				if i2 == numTotal-1 {
					break
				}
			}
			break
		}
	}

	if len(least3Files) < numTotal {
		lg("not enough files in rss fetcher cache")
		return
	}

	lg("fils2")
	for _, v := range least3Files {
		lg("    %v", v)
	}

	logdir := prepareLogDir()

	iter := make([]int, numTotal)

	for i, _ := range iter {

		var doc *html.Node
		surl := spf("%v/%v", hostWithPref, least3Files[i])

		fNamer := FileNamer(logdir, i)
		fnKey := fNamer() // first call yields key

		resBytes, effUrl, err := fetch.UrlGetter(nil, fetch.Options{URL: surl})
		if err != nil {
			lge(err)
			return
		}
		lg("fetched %4.1fkB from %v", float64(len(resBytes))/1024, stringspb.ToLenR(effUrl.String(), 60))

		resBytes = globFixes(resBytes)
		doc, err = html.Parse(bytes.NewReader(resBytes))
		if err != nil {
			lge(err)
			return
		}

		osutilpb.Dom2File(fNamer()+".html", doc)

		//
		//
		cleanseDom(doc, 0)
		removeCommentsAndIntertagWhitespace(NdX{doc, 0})
		reIndent(doc, 0)
		osutilpb.Dom2File(fNamer()+".html", doc)

		//
		//
		{
			removeCommentsAndIntertagWhitespace(NdX{doc, 0})
			condenseTopDown(doc, 0, 0)
			removeEmptyNodes(doc, 0)
		}
		reIndent(doc, 0)
		osutilpb.Dom2File(fNamer()+".html", doc)

		//
		//
		{
			removeCommentsAndIntertagWhitespace(NdX{doc, 0}) // prevent spacey textnodes around singl child images
			breakoutImagesFromAnchorTrees(doc)
		}
		reIndent(doc, 0)
		osutilpb.Dom2File(fNamer()+".html", doc)

		//
		//
		{
			removeCommentsAndIntertagWhitespace(NdX{doc, 0}) // prevent spacey textnodes around singl child images
			// condenseBottomUpV3(doc, 0, 8, map[string]bool{"div": true})
			condenseBottomUpV3(doc, 0, 7, map[string]bool{"div": true})
			condenseBottomUpV3(doc, 0, 6, map[string]bool{"div": true})
			condenseBottomUpV3(doc, 0, 5, map[string]bool{"div": true})
			condenseBottomUpV3(doc, 0, 4, map[string]bool{"div": true})

		}
		removeCommentsAndIntertagWhitespace(NdX{doc, 0}) // prevent spacey textnodes around singl child images
		reIndent(doc, 0)
		osutilpb.Dom2File(fNamer()+".html", doc)

		//
		//
		removeCommentsAndIntertagWhitespace(NdX{doc, 0})
		// proxify(doc, "libertarian-islands.appspot.com", url.Url{Host: remoteHostname})
		proxify(doc, "localhost:8085", &url.URL{Scheme: "http", Host: remoteHostname})
		removeCommentsAndIntertagWhitespace(NdX{doc, 0})
		reIndent(doc, 0)
		osutilpb.Dom2File(fNamer()+".html", doc)

		/*

			//
			//
			{
				removeCommentsAndIntertagWhitespace(NdX{doc, 0}) // prevent id count with textnodes
				addOutlineAttr(doc, 0, []int{0})
				addIdAttr(doc, 0, 1)
			}
			reIndent(doc, 0)
			osutilpb.Dom2File(fNamer()+".html", doc)



		*/

		//
		computeXPathStack(doc, 0)
		osutilpb.Bytes2File(fNamer()+".txt", xPathDump)

		//
		textExtraction(doc, 0)
		textsBytes, textsSorted := orderByOutline(textsByOutl)
		osutilpb.Bytes2File(fNamer()+".txt", textsBytes)
		textsByArticOutl[fnKey] = textsSorted

	}

	// statistics on elements and attributes
	sorted1 := sortmap.SortMapByCount(attrDistinct)
	sorted1.Print(6)
	fmt.Println()
	sorted2 := sortmap.SortMapByCount(nodeDistinct)
	sorted2.Print(6)

	return

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

	pf("correct finish\n")

}

func globFixes(b []byte) []byte {
	// <!--(.*?)-->

	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}

func weedoutFilename(articleId, weedoutStage int) (string, string) {
	stagedFn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
	prefix := fmt.Sprintf("outp_%03v", articleId)
	return stagedFn, prefix
}

func prepareLogDir() string {

	lg, lge := loghttp.Logger(nil, nil)

	logdir := "outp"
	lg("logdir is %v ", logdir)

	// sweep previous
	rmPath := spf("./%v/", logdir)
	err := os.RemoveAll(rmPath)
	if err != nil {
		lge(err)
		os.Exit(1)
	}
	lg("removed %q", rmPath)

	// create anew
	err = os.Mkdir(logdir, 0755)
	if err != nil && !os.IsExist(err) {
		lge(err)
		os.Exit(1)
	}

	return logdir

}

func FileNamer(logdir string, fileNumber int) func() string {
	cntr := -2
	return func() string {
		cntr++
		if cntr == -1 {
			return spf("outp_%03v", fileNumber) // prefix/filekey
		} else {
			fn := spf("outp_%03v_%v", fileNumber, cntr) // filename with stage
			fn = filepath.Join(logdir, fn)
			return fn
		}
	}
}
