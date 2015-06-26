// Package tokenize is a try in splitting a html file
// into tokens, prior to building a dom.
package tokenize

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pbberlin/tools/sort/sortmap"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/util"

	"golang.org/x/net/html"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var pf func(format string, a ...interface{}) (int, error) = fmt.Printf

func Tokenize() {

	extension := ".html"
	directory := ""

	ss := util.GetFilesByExtension(directory, extension, false)
	pss := stringspb.IndentedDump(ss)
	pf("%v \n\n", *pss)

	if len(ss) < 1 {
		pf("did not find any files with %q\n", extension)
		return
	}

	ss = ss[0:1]

	for i := 0; i < len(ss); i++ {
		sb, err := ioutil.ReadFile(ss[i])
		if err != nil {
			pf("%v \n", err)
		}

		r := bytes.NewReader(sb)
		b, err := cleanseHtml(r)
		if err != nil {
			pf("%v \n", err)
		}

		util.WriteBytesToFilename("xx_"+ss[i], b)

		//
		pf("\n\n")
		r = bytes.NewReader(b.Bytes())
		decomposeHtml(r)

	}

}

func cleanseHtml(r io.Reader) (*bytes.Buffer, error) {

	skip := map[string]string{
		"script":   "skip",
		"noscript": "skip",
		"link":     "skip",
		"meta":     "skip",
		"iframe":   "skip",
	}

	b := new(bytes.Buffer)

	d := html.NewTokenizer(r)
	cntrErr := 0
	cntrTkn := 0
	fuckOff := false
	for {
		tokenType := d.Next()
		cntrTkn++

		if tokenType == html.ErrorToken {
			cntrErr++
			if cntrErr > 5 {
				return b, errors.New(spf("error loop at pos %v", cntrTkn))
			}
			continue
		}

		token := d.Token()
		s2 := strings.TrimSpace(string(token.Data))
		attr := getAttr(token.Attr)

		cntrErr = 0
		switch tokenType {
		case html.StartTagToken:
			if _, ok := skip[s2]; ok {
				fuckOff = true
			} else {
				s2 = "\n<" + s2 + attr + ">"
			}
		case html.EndTagToken: // </tag>
			if _, ok := skip[s2]; ok {
				fuckOff = false
				s2 = ""
			} else {
				// s2 = "</" + s2 + ">"
				s2 = "\n</" + s2 + ">\n"
			}
		case html.SelfClosingTagToken:
			if _, ok := skip[s2]; ok {
				s2 = ""
			} else {
				s2 = "\n<" + s2 + attr + "/>\n"
			}
		case html.DoctypeToken:
			s2 = "<!DOCTYPE " + s2 + `><meta content="text/html; charset=utf-8" http-equiv="Content-Type"/>`

		case html.TextToken:
			// nothing
		case html.CommentToken:
			s2 = ""
		default:
			// nothing
		}

		if !fuckOff {
			b.WriteString(s2)
		} else {
			if s2 != "" {
				s2 = strings.Replace(s2, "\n", "", -1)
				s2 = stringspb.Ellipsoider(s2, 30)
				pf("skipped %v \n", s2)

			}
		}
	}
	return b, nil

}

// src http://golang-examples.tumblr.com/page/2
func decomposeHtml(r io.Reader) {

	// type Token struct {
	//     Type     TokenType
	//     DataAtom atom.Atom
	//     Data     string
	//     Attr     []Attribute
	// }
	// type Attribute struct {
	//     Namespace, Key, Val string
	// }

	skip := map[string]string{
		"meta":       "skip",
		"html":       "skip",
		"head":       "skip",
		"title":      "skip",
		"body":       "skip",
		"link":       "skip",
		"script":     "skip",
		"noscript":   "skip",
		"----------": "skip",
		"iframe":     "skip",
		"nav":        "skip",
		"form":       "skip",
	}
	histogram := map[string]interface{}{}

	d := html.NewTokenizer(r)
	cntrErr := 0
	cntrTkn := 0
	for {
		tokenType := d.Next()
		cntrTkn++

		if tokenType == html.ErrorToken {
			pf("#%v err ", cntrTkn)
			cntrErr++
			if cntrErr > 5 {
				break
			}
			continue
		}

		token := d.Token()
		cntrErr = 0
		s1 := strings.TrimSpace(spf(" %#v", token))
		s2 := strings.TrimSpace(string(token.Data))
		s3 := string(token.DataAtom)
		_, _, _ = s1, s2, s3

		switch tokenType {
		case html.StartTagToken, html.SelfClosingTagToken:
			if _, ok := skip[s2]; !ok {
				pf("\n%v ", s2)
				if _, ok := histogram[s2]; !ok {
					histogram[s2] = 1
				} else {
					val := histogram[s2].(int)
					histogram[s2] = val + 1
				}
			}
		case html.TextToken:
			if s2 != "" && len(s2) > 1 && !strings.HasPrefix(s2, `//`) {
				s2 = strings.Replace(s2, "\n", "", -1)
				pf("\t%v", stringspb.Ellipsoider(s2, 22))
			}
		case html.EndTagToken: // </tag>
			// pf("/%v ", s2)
		case html.CommentToken:
			// pf("comment ")
		case html.DoctypeToken:

		default:
			pf("default case %v\n", s1)
		}
	}

	hSort := sortmap.StringKeysToSortedArray(histogram)

	pf("\n\n")
	for _, v := range hSort {
		pf("%10s %4v\n", v, histogram[v])
	}

}

// type Attribute struct {
//     Namespace, Key, Val string
// }
func getAttr(attributes []html.Attribute) string {
	ret := ""
	for i := 0; i < len(attributes); i++ {
		attr := attributes[i]
		if attr.Key == "href" || attr.Key == "src" || attr.Key == "alt" {
			if attr.Key == "src" && strings.HasPrefix(attr.Val, "/") {
				attr.Val = "http://handelsblatt.com" + attr.Val
			}
			ret += spf(" %v=%q ", attr.Key, attr.Val)
		}
	}
	return ret
}

func getAttrVal(attributes []html.Attribute, key string) string {
	for i := 0; i < len(attributes); i++ {
		attr := attributes[i]
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
