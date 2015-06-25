package charting

import (
	"github.com/pbberlin/tools/net/http/loghttp"

	"appengine"

	"image"

	"image/jpeg"

	// Package image/png ... imported for its initialization side-effect
	// => image.Decode understands PNG formatted images.
	_ "image/jpeg"
	_ "image/png"
	// uncomment to allow understanding of GIF images...
	// _ "image/gif"

	"io"
	"net/http"
	"os"
)

func imgServingExample1(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	p := r.FormValue("p")
	if p == "" {
		p = "static/chartbg_400x960__480x1040__12x10.png"
	}
	if p == "" {
		p = "static/pberg1.png"
	}

	f, err := os.Open(p)
	loghttp.E(w, r, err, false)
	defer f.Close()

	mode := r.FormValue("mode")
	if mode == "" {
		mode = "direct"
	}

	if mode == "direct" {

		// file reader directly to http writer

		w.Header().Set("Content-Type", "image/png")
		c.Infof("serving directly %v \n", p)
		io.Copy(w, f)

	} else {

		// file to memory image - memory image to http writer

		img, whichFormat, err := image.Decode(f)
		loghttp.E(w, r, err, false)
		c.Infof("serving as memory image %v - format %v - type %T\n", p, whichFormat, img)

		// initially a png - we now encode it to jpg
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	}

}

func init() {

	http.HandleFunc("/img-serve-example-1", imgServingExample1)
}
