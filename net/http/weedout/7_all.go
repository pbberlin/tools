package weedout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// Puttting it all together
func WeedOut(least3Files []repo.FullArticle, lg loghttp.FuncBufUniv, fs fsi.FileSystem) *html.Node {

	//
	// domclean
	for i := 0; i < len(least3Files); i++ {

		fNamer := domclean2.FileNamer(logDir, i)
		fNamer() // first call yields key

		lg("cleaning %4.1fkB from %v", float64(len(least3Files[i].Body))/1024,
			stringspb.ToLenR(least3Files[i].Url, 60))

		opts := domclean2.CleaningOptions{Proxify: true, Beautify: true}
		// opts.FNamer = fNamer
		opts.AddOutline = true
		opts.RemoteHost = fetch.HostFromStringUrl(least3Files[i].Url)
		doc, err := domclean2.DomClean(least3Files[i].Body, opts)
		lg(err)

		fileDump(lg, fs, doc, fNamer, ".html")

	}

	if false {
		//
		// Textify with brute force
		for i := 0; i < len(least3Files); i++ {

			fNamer := domclean2.FileNamer(logDir, i)
			fNamer() // first call yields key

			bts, err := fs.ReadFile(fNamer() + ".html")
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
	}

	//
	// Textify with more finetuning.
	// Save result to memory.
	textsByArticOutl := map[string][]*TextifiedTree{}
	for i := 0; i < len(least3Files); i++ {

		fNamer := domclean2.FileNamer(logDir, i)
		fnKey := fNamer() // first call yields key

		bts, err := fs.ReadFile(fNamer() + ".html")

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

		similaritiesToFile(fs, logDir, frags, weedStage)

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

	bts, err := fs.ReadFile(fNamer() + ".html")
	lg(err)
	doc, err := html.Parse(bytes.NewReader(bts))
	lg(err)

	weedoutApply(doc, skipPrefixes)

	domclean2.DomFormat(doc)

	return doc
}

func FetchAndDecodeJSON(r *http.Request, surl, knownProtocol string, lg loghttp.FuncBufUniv, fs fsi.FileSystem) []repo.FullArticle {

	fullURL := fmt.Sprintf("%s%s?%s=%s&cnt=%v&prot=%v", routes.AppHost01, routes.FetchSimilarURI,
		routes.URLParamKey, surl, numTotal-1, knownProtocol)

	// fullURL = fmt.Sprintf("%s%s?%s=%s&cnt=%v", r.URL.Host, repo.routes.FetchSimilarURI,
	// 	routes.URLParamKey, surl, numTotal-1)

	lg("lo fetching %v", fullURL)
	start := time.Now()

	fo := fetch.Options{}
	fo.URL = fullURL
	bJSON, inf, err := fetch.UrlGetter(r, fo)
	_ = inf
	lg(err)
	if err != nil {
		lg("msg %v", inf.Msg)
		return nil
	}
	if len(bJSON) == 0 {
		lg("empty bJSON")
		return nil
	}

	lg("\t\tfetch resp complete after %4.2v secs", time.Now().Sub(start).Seconds())

	var mp map[string][]byte
	err = json.Unmarshal(bJSON, &mp)
	lg(err)
	if err != nil {
		if _, ok := mp["msg"]; ok {
			lg("%s", mp["msg"])
		} else {
			lg("%s", bJSON)
		}
		return nil
	}

	smaxFound := string(mp["lensimilar"])
	maxFound := util.Stoi(smaxFound)
	if maxFound < numTotal-1 {
		lg("not enough files returned by FetchSimilar 1 - mp[lensimilar] too small: %s", mp["lensimilar"])
		return nil
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
				return nil
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

	lg("found %v similar; decoding complete after %4.2v secs", maxFound, time.Now().Sub(start).Seconds())

	for _, v := range least3Files {
		lg("%v %v", v.Url, len(v.Body))
	}

	return least3Files

}
