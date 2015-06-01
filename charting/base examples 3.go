package charting

import (
	"appengine"

	"image"

	"image/jpeg"

	// Package image/png ... imported for its initialization side-effect
	// => image.Decode understands PNG formatted images.
	_ "image/jpeg"
	_ "image/png"
	// uncomment to allow understanding of GIF images...
	// _ "image/gif"

	"net/http"
	"os"

	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"
)

func imgServingExample3(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	p := r.FormValue("p")
	if p == "" {
		p = "static/chartbg_400x960__480x1040__12x10.png"
	}
	if p == "" {
		p = "static/pberg1.png"
	}
	// try p=static/unmodifiable_format.jpg

	// prepare a cutout rect
	var p1, p2 image.Point
	p1.X, p1.Y = 10, 60
	p2.X, p2.Y = 400, 255
	var rec image.Rectangle = image.Rectangle{Min: p1, Max: p2}

	f, err := os.Open(p)
	util_err.Err_http(w, r, err, false)
	defer f.Close()

	img, whichFormat, err := image.Decode(f)
	util_err.Err_http(w, r, err, false, "only jpeg and png are 'activated' ")
	c.Infof("serving format %v %T\n", whichFormat, img)

	switch t := img.(type) {

	default:
		util_err.Err_http(w, r, false, false, "internal color formats image.YCbCr and image.RGBA are understood")

	case *image.RGBA, *image.YCbCr:
		imgXFull, ok := t.(*image.RGBA)
		util_err.Err_http(w, r, ok, false, "image.YCbCr can not be typed to image.RGBA - this will panic")

		imgXCutout, ok := imgXFull.SubImage(rec).(*image.RGBA)
		util_err.Err_http(w, r, ok, false, "cutout operation failed")

		// we serve it as JPEG
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, imgXCutout, &jpeg.Options{Quality: jpeg.DefaultQuality})

	}

}

func init() {
	http.HandleFunc("/img-serve-example-3", util_appengine.Adapter(imgServingExample3))
}
