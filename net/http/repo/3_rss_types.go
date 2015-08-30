package repo

import (
	"encoding/xml"
	"fmt"
)

var pf = fmt.Printf
var pfRestore = fmt.Printf

var spf = fmt.Sprintf
var wpf = fmt.Fprintf

// RSS directory
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   Items    `xml:"channel"`
}

// List of elements of an RSS directory
type Items struct {
	XMLName  xml.Name `xml:"channel"`
	ItemList []Item   `xml:"item"`
}

// Element of an RSS directory
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
