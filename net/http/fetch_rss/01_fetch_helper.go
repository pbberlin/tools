package fetch_rss

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/url"
	"path"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/stringspb"
)

//
func rssXMLFile(rssUrl string) (rssDoc RSS, rssUrlObj *url.URL) {

	var bts []byte
	var err error

	bts, rssUrlObj, err = fetch.UrlGetter(rssUrl, nil, false)
	logif.F(err)

	bts = bytes.Replace(bts, []byte("content:encoded>"), []byte("content-encoded>S"), -1) // hack

	err = xml.Unmarshal(bts, &rssDoc)
	logif.E(err)

	// save it
	bdmp := stringspb.IndentedDumpBytes(rssDoc)
	err = fs.WriteFile(path.Join(docRoot, rssUrlObj.Host, "outp_rss.xml"), bdmp, 0755)
	logif.E(err)
	pf("RSS resp size, outp_rss.xml, : %v\n", len(bdmp))

	return
}

func condenseTrailingDir(uri string, n int) (ret string) {

	switch n {
	case 1:
		return uri
	case 2:
		base1 := path.Base(uri)
		rdir1 := path.Dir(uri) // rightest Dir

		base2 := path.Base(rdir1)

		rdir2 := path.Dir(rdir1)

		ret = path.Join(rdir2, base2+"-"+base1)
	default:
		log.Fatalf("not implemented n > 2")
	}

	return

}
