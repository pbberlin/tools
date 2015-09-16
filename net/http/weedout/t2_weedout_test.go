// +build weed2
// go test -tags=weed2

package weedout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"appengine"
	"appengine/aetest"

	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

var logDir = "c:/tmp/weedout2/"

var memMapFileSys = memfs.New(memfs.DirSort("byDateDesc")) // package variable required as "persistence"

func GetFS(c appengine.Context, whichType int) (fs fsi.FileSystem) {
	switch whichType {
	case 0:
		// re-instantiation would delete contents
		fs = fsi.FileSystem(memMapFileSys)
	case 1:
		// must be re-instantiated for each request
		dsFileSys := dsfs.New(dsfs.DirSort("byDateDesc"), dsfs.MountName("mntTest"), dsfs.AeContext(c))
		fs = fsi.FileSystem(dsFileSys)
	case 2:

		osFileSys := osfs.New(osfs.DirSort("byDateDesc"))
		fs = fsi.FileSystem(osFileSys)
		os.Chdir(logDir)
	default:
		panic("invalid whichType ")
	}

	return
}

func Test2(t *testing.T) {

	lg, b := loghttp.BuffLoggerUniversal(nil, nil)
	_ = b
	// closureOverBuf := func(bUnused *bytes.Buffer) {
	// 	loghttp.Pf(nil, nil, b.String())
	// }
	// defer closureOverBuf(b) // the argument is ignored,

	remoteHostname := "www.welt.de"

	c, err := aetest.NewContext(nil)
	lg(err)
	if err != nil {
		return
	}
	defer c.Close()

	fs := GetFS(c, 2)

	var urls = []string{
		"www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html",
		"www.economist.com/news/britain/21663648-hard-times-hard-hats-making-britain-make-things-again-proving-difficult",
		"www.economist.com/news/americas/21661804-gender-equality-good-economic-growth-girl-power",
	}

	fullURL := fmt.Sprintf("https://%s%s?%s=%s&cnt=%v", cTestHostDev, repo.UriFetchSimilar,
		routes.URLParamKey, urls[0], numTotal-1)
	lg("lo sending to URL:")
	lg("lo %v", fullURL)

	fo := fetch.Options{}
	fo.URL = fullURL
	bResp, inf, err := fetch.UrlGetter(nil, fo)
	_ = inf
	lg(err)
	if err != nil {
		return
	}
	if len(bResp) == 0 {
		lg("empty bResp")
		return
	}

	var mp map[string][]byte
	err = json.Unmarshal(bResp, &mp)
	lg(err)
	if err != nil {
		if _, ok := mp["msg"]; ok {
			lg("%s", mp["msg"])
		} else {
			lg("%s", bResp)
		}
		return
	}

	smaxFound := string(mp["lensimilar"])
	maxFound := util.Stoi(smaxFound)
	if maxFound < numTotal-1 {
		lg("not enough files returned by FetchSimilar 1 - mp[lensimilar] too small: %s", mp["lensimilar"])
		return
	}
	least3Files := make([]repo.FullArticle, maxFound+1)

	_, ok1 := mp["url_self"]
	_, ok2 := mp["mod_self"]
	_, ok3 := mp["bod_self"]
	if ok1 && ok2 && ok3 {
		least3Files[0].Url = string(mp["url_self"])
		least3Files[0].Mod, err = time.Parse(http.TimeFormat, string(mp["mod_self"]))
		lg(err)
		least3Files[0].Body = mp["bod_self"]
		if len(least3Files[0].Body) < 200 {
			if !bytes.Contains(least3Files[0].Body, []byte(fetch.MsgNoRdirects)) {
				lg("found base but its a redirect")
				return
			}
		}
	}
	lg("found base")

	for k, v := range mp {
		if k == "msg" {
			continue
		}
		if strings.HasSuffix(k, "self") {
			continue
		}

		if strings.HasPrefix(k, "url__") {
			sval := strings.TrimPrefix(k, "url__")
			val := util.Stoi(sval)
			// lg("%v %v %s", sval, val, v)
			least3Files[val+1].Url = string(v)
		}
		if strings.HasPrefix(k, "mod__") {
			sval := strings.TrimPrefix(k, "mod__")
			val := util.Stoi(sval)
			// lg("%v %v %s", sval, val, v)
			least3Files[val+1].Mod, err = time.Parse(http.TimeFormat, string(v))
			lg(err)
		}

		if strings.HasPrefix(k, "bod__") {
			sval := strings.TrimPrefix(k, "bod__")
			val := util.Stoi(sval)
			least3Files[val+1].Body = v //html.EscapeString(string(v)
		}

	}

	lg("found %v similar", maxFound)

	for _, v := range least3Files {
		lg("%v %v", v.Url, len(v.Body))
	}

	//
	// domclean
	for i := 0; i < numTotal; i++ {

		fNamer := domclean2.FileNamer(logDir, i)
		fNamer() // first call yields key

		lg("cleaning %4.1fkB from %v", float64(len(least3Files[i].Body))/1024,
			stringspb.ToLenR(least3Files[i].Url, 60))

		opts := domclean2.CleaningOptions{Proxify: true, Beautify: true}
		// opts.FNamer = fNamer
		opts.AddOutline = true
		opts.RemoteHost = remoteHostname
		doc, err := domclean2.DomClean(least3Files[i].Body, opts)
		lg(err)

		fileDump(lg, fs, doc, fNamer, ".html")

	}

	//
	// Textify with brute force
	for i := 0; i < numTotal; i++ {

		fNamer := domclean2.FileNamer(logDir, i)
		fNamer() // first call yields key

		bts, err := common.ReadFile(fs, fNamer()+".html")
		lg(err)
		doc, err := html.Parse(bytes.NewReader(bts))
		lg(err)

		textifyBruteForce(doc)

		var buf bytes.Buffer
		err = html.Render(&buf, doc)
		lg(err)

		b := buf.Bytes()
		b = bytes.Replace(b, []byte("[br]"), []byte("\n"), -1)

		fileDump(lg, fs, b, fNamer, "_raw.txt")
	}

	//
	// Textify with more finetuning.
	// Save result to memory.
	textsByArticOutl := map[string][]*TextifiedTree{}
	for i := 0; i < numTotal; i++ {

		fNamer := domclean2.FileNamer(logDir, i)
		fnKey := fNamer() // first call yields key

		bts, err := common.ReadFile(fs, fNamer()+".html")

		doc, err := html.Parse(bytes.NewReader(bts))
		lg(err)

		fNamer() // one more

		//
		mp, bts := BubbledUpTextExtraction(doc, fnKey)
		fileDump(lg, fs, bts, fNamer, ".txt")

		mpSorted, dump := orderByOutline(mp)
		fileDump(lg, fs, dump, fNamer, ".txt")
		textsByArticOutl[fnKey] = mpSorted

		// for k, v := range mpSorted {
		// 	if k%33 != 0 {
		// 		continue
		// 	}
		// 	log.Printf("%3v: %v %14v  %v\n", k, v.SourceID, v.Outline, v.Lvl)
		// }

	}

	//
	//
	// We progress from level 1 downwards.
	// Lower levels skip weeded out higher levels,
	// to save expensive levenshtein comparisons
	var skipPrefixes = map[string]bool{}
	for weedStage := 1; weedStage <= stageMax; weedStage++ {

		fNamer := domclean2.FileNamer(logDir, 0)
		fnKey := fNamer() // first call yields key

		levelsToProcess = map[int]bool{weedStage: true}
		frags := similarTextifiedTrees(textsByArticOutl, skipPrefixes, map[string]bool{fnKey: true})

		similaritiesToFile(logDir, frags, weedStage)

		for _, frag := range frags {
			if len(frag.Similars) >= numTotal-1 &&
				frag.SumRelLevenshtein/(numTotal-1) < 0.2 {
				skipPrefixes[frag.Outline+"."] = true
			}
		}
		b := new(bytes.Buffer)
		for k, _ := range skipPrefixes {
			b.WriteString(k)
			b.WriteByte(32)
		}
		// log.Printf("%v\n", b.String())

	}

	//
	// Apply weedout
	fNamer := domclean2.FileNamer(logDir, 0)
	fNamer() // first call yields key

	bts, err := common.ReadFile(fs, fNamer()+".html")
	lg(err)
	doc, err := html.Parse(bytes.NewReader(bts))
	lg(err)

	weedoutApply(doc, skipPrefixes)

	domclean2.DomFormat(doc)

	fileDump(lg, fs, doc, fNamer, ".html")

	pf("MapSimiliarCompares: %v SimpleCompares: %v LevenstheinComp: %v\n", breakMapsTooDistinct, appliedLevenshtein, appliedCompare)
	pf("correct finish\n")

}

func fileDump(lg loghttp.FuncBufUniv, fs fsi.FileSystem,
	content interface{}, fNamer func() string, secondPart string) {

	if fNamer != nil {
		fn := fNamer() + secondPart
		switch casted := content.(type) {
		case *html.Node:
			var b bytes.Buffer
			err := html.Render(&b, casted)
			lg(err)
			if err != nil {
				return
			}
			common.WriteFile(fs, fn, b.Bytes())
		case []byte:
			common.WriteFile(fs, fn, casted)
		}

	}

}
