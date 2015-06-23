package parse2

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"path/filepath"
	"time"

	urlX "net/url"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/util"
)

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

	resBytes, err := pbfetch.UrlGetter(rssUrl, nil, false)
	if err != nil {
		log.Fatal(err)
	}

	resBytes = bytes.Replace(resBytes, []byte("content:encoded>"), []byte("content-encoded>S"), -1)

	var rssDoc RSS
	err = xml.Unmarshal(resBytes, &rssDoc)
	if err != nil {
		pf("%v\n", err)
	}

	bsd := pbstrings.IndentedDumpBytes(rssDoc)
	pf("RSS resp size: %v\n%s\n", len(bsd), bsd[:util.Min(5*excerptLen, len(bsd)-1)])

	items := rssDoc.Items
	for i := 0; i < len(items.ItemList); i++ {
		lpItem := items.ItemList[i]
		pf("%v: %v - %v\n", i, lpItem.Published[5:22], lpItem.Link)

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

	for _, a := range fullArticles {
		fn := fileName(a.Url)
		bytes2File(fn, a.Body)
	}

}

func fileName(url string) string {

	u, err := urlX.Parse(url)
	if err != nil {
		panic(fmt.Errorf("url unparseable: %v", err))
	}
	s := u.RequestURI()
	fn := filepath.Base(s)
	pf("fn: %v\n", fn)
	return fn

}
