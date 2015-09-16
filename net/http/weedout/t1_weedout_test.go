// +build weed1
// go test -tags=weed1

package weedout

import (
	"bytes"
	"path"
	"testing"

	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/os/osutilpb"
	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	remoteHostname := "www.welt.de"
	remoteHostname = "www.welt.de/politik/ausland"

	dirs1, _, msg, err := fileserver.GetDirContents(repo.RepoURL, remoteHostname)
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

		p := path.Join(remoteHostname, v1)
		dirs2, fils2, msg, err := fileserver.GetDirContents(repo.RepoURL, p)
		_ = dirs2
		if err != nil {
			lge(err)
			lg("%s", msg)
		}
		// lg("  dirs2 %v", stringspb.IndentedDump(dirs2))
		// lg("  fils2 %v", stringspb.IndentedDump(fils2))

		for _, v2 := range fils2 {
			least3Files = append(least3Files, path.Join(remoteHostname, v1, v2))
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

	logdir := osutilpb.PrepareLogDir()

	iter := make([]int, numTotal)

	//
	// domclean
	for i, _ := range iter {

		surl := spf("%v/%v", repo.RepoURL, least3Files[i])

		fNamer := domclean2.FileNamer(logdir, i)
		fNamer() // first call yields key

		resBytes, inf, err := fetch.UrlGetter(nil, fetch.Options{URL: surl})
		if err != nil {
			lge(err)
			return
		}
		lg("fetched %4.1fkB from %v", float64(len(resBytes))/1024, stringspb.ToLenR(inf.URL.String(), 60))
		opts := domclean2.CleaningOptions{Proxify: true, Beautify: true}
		// opts.FNamer = fNamer
		opts.AddOutline = true
		opts.RemoteHost = remoteHostname
		doc, err := domclean2.DomClean(resBytes, opts)

		osutilpb.Dom2File(fNamer()+".html", doc)

	}

	//
	// Textify with brute force
	for i, _ := range iter {

		fNamer := domclean2.FileNamer(logdir, i)
		fNamer() // first call yields key

		bts := osutilpb.BytesFromFile(fNamer() + ".html")
		doc, err := html.Parse(bytes.NewReader(bts))
		lge(err)

		textifyBruteForce(doc)

		var buf bytes.Buffer
		err = html.Render(&buf, doc)
		lge(err)

		b := buf.Bytes()
		b = bytes.Replace(b, []byte("[br]"), []byte("\n"), -1)

		osutilpb.Bytes2File(fNamer()+"_raw"+".txt", b)
	}

	//
	// Textify with more finetuning.
	// Save result to memory.
	textsByArticOutl := map[string][]*TextifiedTree{}
	for i, _ := range iter {

		fNamer := domclean2.FileNamer(logdir, i)
		fnKey := fNamer() // first call yields key

		bts := osutilpb.BytesFromFile(fNamer() + ".html")
		doc, err := html.Parse(bytes.NewReader(bts))
		lge(err)

		fNamer() // one more

		//
		mp, bts := BubbledUpTextExtraction(doc, fnKey)
		osutilpb.Bytes2File(fNamer()+".txt", bts)

		mpSorted, dump := orderByOutline(mp)
		osutilpb.Bytes2File(fNamer()+".txt", dump)
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

		fNamer := domclean2.FileNamer(logdir, 0)
		fnKey := fNamer() // first call yields key

		levelsToProcess = map[int]bool{weedStage: true}
		frags := similarTextifiedTrees(textsByArticOutl, skipPrefixes, map[string]bool{fnKey: true})

		similaritiesToFile(logdir, frags, weedStage)

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
	fNamer := domclean2.FileNamer(logdir, 0)
	fNamer() // first call yields key

	bts := osutilpb.BytesFromFile(fNamer() + ".html")
	doc, err := html.Parse(bytes.NewReader(bts))
	lge(err)

	weedoutApply(doc, skipPrefixes)

	domclean2.DomFormat(doc)

	osutilpb.Dom2File(fNamer()+".html", doc)

	pf("MapSimiliarCompares: %v SimpleCompares: %v LevenstheinComp: %v\n", breakMapsTooDistinct, appliedLevenshtein, appliedCompare)
	pf("correct finish\n")

}

// func weedoutFilename(articleId, weedoutStage int) (string, string) {
// 	stagedFn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
// 	prefix := fmt.Sprintf("outp_%03v", articleId)
// 	return stagedFn, prefix
// }
