package charting

import (
	"appengine"

	"image"
	"image/color"
	"image/jpeg"

	"net/http"
)

func imgServingExample2(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// prepare a rectangle
	var p1, p2 image.Point
	p1.X, p1.Y = 10, 10
	p2.X, p2.Y = 400, 255
	var rec image.Rectangle = image.Rectangle{Min: p1, Max: p2}

	// prepare a line color
	lineCol := color.RGBA{}
	lineCol.R, lineCol.G, lineCol.B = 255, 44, 22
	//lineCol.A = 0
	c.Infof("brush color %#v \n", lineCol)

	// create empty memory image - all pixels are 0 => black
	imgRGBA := image.NewRGBA(rec)

	for i := 20; i < 140; i++ {
		lineCol.A = uint8(i)
		imgRGBA.Set(i, i, lineCol)
		imgRGBA.Set(i+1, i, lineCol)
		imgRGBA.Set(i+2, i, lineCol)
		imgRGBA.Set(i+3, i, lineCol)
		imgRGBA.Set(i+3, i+1, lineCol)
	}
	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, imgRGBA, &jpeg.Options{Quality: jpeg.DefaultQuality})

}

func init() {
	http.HandleFunc("/img-serve-example-2", imgServingExample2)
}
