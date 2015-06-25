package charting

import (
	"image"
	"image/png"
	"strings"

	"appengine"

	// Package image/png ... imported for its initialization side-effect
	// => image.Decode understands PNG formatted images.
	_ "image/jpeg"
	_ "image/png"
	// uncomment to allow understanding of GIF images...
	// _ "image/gif"

	"io"
	"net/http"
	"os"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/net/http/loghttp"
)

// output static file as base64 string
//   image must exist as file in /static
func imagefileAsBase64(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	p := r.FormValue("p")
	if p == "" {
		p = "static/chartbg_400x960__480x1040__12x10.png"
	}

	f, err := os.Open(p)
	loghttp.E(w, r, err, false)
	defer f.Close()

	img, whichFormat, err := image.Decode(f)
	loghttp.E(w, r, err, false)
	c.Infof("format %v - subtype %T\n", whichFormat, img)

	imgRGBA, ok := img.(*image.RGBA)
	loghttp.E(w, r, ok, true, "source image was not convertible to image.RGBA - gifs or jpeg?")

	// => header with mime prefix always prepended
	//   and its always image/PNG
	str_b64 := conv.Rgba_img_to_base64_str(imgRGBA)

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<p>Image embedded in HTML as Base64:</p><img width=200px src=\"")
	io.WriteString(w, str_b64)
	io.WriteString(w, "\"> ")

}

// output variable as base64 string
//  we have only two such variables as constants
func imagevariAsBase64(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	p := r.FormValue("p")
	if p == "" {
		p = "1"
	}

	var str_src string
	if p == "1" {
		str_src = conv.Img_jpeg_base64
	}
	if p == "2" {
		str_src = conv.Img_rgba_base64
	}
	if p == "3" {
		str_src = conv.Img_rgba_base64_old
	}

	// here we could check integrity...
	//img,whichFormat := conv.Base64_str_to_img(str_src)
	//c.Infof( "retrieved img from base64: format %v - subtype %T\n" , whichFormat, img )

	rdr := strings.NewReader(str_src)
	explicitMime := ""
	existingMime := conv.MimeFromBase64(rdr)
	if existingMime == "" {
		explicitMime = "data:image/jpeg;base64,"
	}

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<p>Image embedded in HTML as Base64:</p><img width=200px src=\"")
	io.WriteString(w, explicitMime)
	io.WriteString(w, str_src)
	io.WriteString(w, "\"> ")

}

func datastoreAsBase64(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	p := r.FormValue("p")
	if p == "" {
		p = "chart1"
	}

	dsObj, _ := dsu.BufGet(w, r, "dsu.WrapBlob__"+p)

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<p>Image embedded in HTML as Base64:</p><img width=200px src=\"")
	io.WriteString(w, string(dsObj.VByte))
	io.WriteString(w, "\"> ")

}

func imageFromDatastore(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	p := r.FormValue("p")
	if p == "" {
		p = "chart1"
	}

	img := GetImageFromDatastore(w, r, p)

	// regardless of previous image type
	//   we encode as png
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)

}

func init() {

	http.HandleFunc("/image/base64-from-file", loghttp.Adapter(imagefileAsBase64))
	http.HandleFunc("/image/base64-from-var", loghttp.Adapter(imagevariAsBase64))
	http.HandleFunc("/image/base64-from-datastore", loghttp.Adapter(datastoreAsBase64))

	http.HandleFunc("/image/img-from-datastore", loghttp.Adapter(imageFromDatastore))

}
