package fetch_rss

import (
	"encoding/xml"
	"fmt"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf

var spf func(format string, a ...interface{}) string = fmt.Sprintf

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
