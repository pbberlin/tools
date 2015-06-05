package parse

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pbberlin/tools/util"
)

func init() {

}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   Items    `xml:"channel"`
}
type Items struct {
	XMLName  xml.Name `xml:"channel"`
	ItemList []Item   `xml:"item"`
}
type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	GUID        string    `xml:"guid"`
	Description string    `xml:"description"`
	Category    string    `xml:"category"`
	Published   string    `xml:"pubDate"`
	Enc         Enclosure `xml:"enclosure"`
	Content     string    `xml:"content-encoded"`
}

type Enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
	Len  int    `xml:"length,attr"`
}

type FullArticle struct {
	URL  string
	Body *[]byte
}

var fullArticles []FullArticle
var c chan *FullArticle = make(chan *FullArticle)

func Fetch(amount int) {

	go func() {
		for {
			pfa := <-c
			fa := *pfa
			fullArticles = append(fullArticles, fa)
			pf("done fetching %v \n", fa.URL[27:])
		}
	}()

	// cx := appengine.NewContext(r)
	// cl := urlfetch.Client(cx)
	cl := http.DefaultClient

	resp, err := cl.Get("http://www.handelsblatt.com/contentexport/feed/schlagzeilen")
	if err != nil {
		pf("%v\n", err)
	}
	bcntent, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		pf("%v\n", err)
	}

	bcntent = bytes.Replace(bcntent, []byte("content:encoded>"), []byte("content-encoded>S"), -1)

	// scntent := string(bcntent)
	// pf("size: %v \n%v\n", len(scntent), util.Ellipsoider(scntent, 1450))

	var rssDoc RSS
	err = xml.Unmarshal(bcntent, &rssDoc)
	if err != nil {
		pf("%v\n", err)
	}

	ps := util.IndentedDump(rssDoc)
	s := *ps
	pf("- %v - \n%v\n", len(s), s[:util.Min(1600, len(s)-1)])

	items := rssDoc.Items
	for i := 0; i < len(items.ItemList); i++ {
		lpItem := items.ItemList[i]
		pf("%v: %v - %v\n", i, lpItem.Published[5:22], lpItem.Link)

		go func(argURL string) {
			cl := http.DefaultClient
			resp, err := cl.Get(argURL)
			if err != nil {
				pf(" full art %v %v\n", argURL, err)
			}
			bcntent, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				pf(" full art %v %v\n", argURL, err)
			}
			fa := FullArticle{}
			fa.URL = argURL
			fa.Body = &bcntent
			c <- &fa
		}(lpItem.Link)

		if i+1 >= amount {
			break
		}
	}

	time.Sleep(4 * time.Second)
	pf("\n\n\n")
	for i := 0; i < len(fullArticles); i++ {
		lpFa := fullArticles[i]
		bBody := *fullArticles[i].Body
		// pf("%v: %v\n\n", lpFa.URL[27:], util.Ellipsoider(string(bBody), 200))

		fileName := lpFa.URL
		fileName = strings.Replace(fileName, "https://", "", -1)
		fileName = strings.Replace(fileName, "http://", "", -1)
		pf("%v\n", fileName)
		fileName = fileName[strings.Index(fileName, "/")+1:]
		fileName = strings.Replace(fileName, "/", "--", 1)
		pf("%v\n", fileName)
		nextSlash := strings.Index(fileName, "/")
		if nextSlash > 0 {
			fileName = fileName[:strings.Index(fileName, "/")]
			fileName += ".html"
		}
		pf("%v\n", fileName)

		f, err := os.Create(fileName)
		if err != nil {
			pf(" file open %v %v\n", fileName, err)
		}
		defer f.Close()
		n2, err := f.Write(bBody)
		pf("wrote %d bytes - err |%v| \n", n2, err)

		err = ioutil.WriteFile("f1.html", bBody, os.ModePerm)
		if err != nil {
			pf("can not write file: %v", err)
		}
		pf("file written")

	}

}
