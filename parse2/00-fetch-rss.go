package parse2

import (
	"bytes"
	"encoding/xml"
	"log"
	"time"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/pblog"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/util"
)

var hosts = []string{"www.handelsblatt.com"}

type FullArticle struct {
	Url  string
	Body []byte
}

var fullArticles []FullArticle
var c chan *FullArticle = make(chan *FullArticle)

func Fetch(rssUrl string, numberArticles int) {

	go func() {
		for {
			pfa := <-c
			fa := *pfa
			fullArticles = append(fullArticles, fa)
			pf("done fetching %v \n", fa.Url[27:])
		}
	}()

	bts, err := pbfetch.UrlGetter(rssUrl, nil, false)
	pblog.Fatal(err)

	bts = bytes.Replace(bts, []byte("content:encoded>"), []byte("content-encoded>S"), -1)

	var rssDoc RSS
	err = xml.Unmarshal(bts, &rssDoc)
	pblog.LogE(err)

	bdmp := pbstrings.IndentedDumpBytes(rssDoc)
	pf("RSS resp size: %v\n%s\n", len(bdmp), bdmp[:util.Min(5*excerptLen, len(bdmp)-1)])
	bytes2File("outp_rss.xml", bdmp)

	items := rssDoc.Items
	for i := 0; i < len(items.ItemList); i++ {
		lpItem := items.ItemList[i]

		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		pblog.LogE(err)
		pf("%2v: %v - %v\n", i, t.Format("2.1. 15:04:05"), lpItem.Link)

		go func(argURL string) {
			bs, err := pbfetch.UrlGetter(argURL, nil, false)
			if err != nil {
				log.Fatal(err)
			}
			fa := FullArticle{argURL, bs}
			c <- &fa
		}(lpItem.Link)

		if i+1 >= numberArticles {
			break
		}
	}

	time.Sleep(4 * time.Second)
	pf("\n\n\n")

	for idx, a := range fullArticles {
		orig, numbered := fetchFileName(a.Url, idx+len(testDocs))
		bytes2File(orig, a.Body)
		bytes2File(numbered, a.Body)
	}

}
