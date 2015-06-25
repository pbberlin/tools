package charting

import (
	"fmt"
	"image"
	"image/color"

	"appengine"

	"image/png"

	// Package image/png ... imported for its initialization side-effect
	// => image.Decode understands PNG formatted images.
	_ "image/jpeg"
	_ "image/png"
	// uncomment to allow understanding of GIF images...
	// _ "image/gif"

	"bytes"
	"net/http"
	"os"

	//"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/net/http/loghttp"
)

// An example demonstrating decoding JPEG img + examining its pixels.
func imageAnalyze(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	// Decode the JPEG data.
	// If reading from file, create a reader with
	// reader, err := os.Open("testdata/video-001.q50.420.jpeg")
	// if err != nil {  c.Errorf(err)  }
	// defer reader.Close()

	img, whichFormat := conv.Base64_str_to_img(conv.Img_jpeg_base64)
	c.Infof("retrieved img from base64: format %v - type %T\n", whichFormat, img)

	bounds := img.Bounds()

	// Calculate a 16-bin histogram for m's red, green, blue and alpha components.
	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	var histogram [16][4]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 12 reduces this to the range [0, 15].
			histogram[r>>12][0]++
			histogram[g>>12][1]++
			histogram[b>>12][2]++
			histogram[a>>12][3]++
		}
	}

	b1 := new(bytes.Buffer)
	s1 := fmt.Sprintf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	b1.WriteString(s1)

	for i, x := range histogram {
		s1 := fmt.Sprintf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
		b1.WriteString(s1)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(b1.Bytes())
}

func drawLinesOverGrid(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	// prepare a line color
	lineCol := color.RGBA{}
	lineCol.R, lineCol.G, lineCol.B, lineCol.A = 255, 244, 22, 0
	c.Infof("brush color %#v \n", lineCol)

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

	img, whichFormat, err := image.Decode(f)
	loghttp.E(w, r, err, false, "only jpeg and png are 'activated' ")
	c.Infof("serving img format %v %T\n", whichFormat, img)

	switch imgXFull := img.(type) {

	default:
		loghttp.E(w, r, false, true, "convertibility into internal color format image.RGBA required")

	case *image.RGBA:

		drawLine := FuncDrawLiner(lineCol, imgXFull)

		xb, yb := 40, 440
		P0 := image.Point{xb + 0, yb - 0}
		drawLine(P0, lineCol, imgXFull)

		for i := 0; i < 1; i++ {

			P1 := image.Point{xb + 80, yb - 80}
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 160
			P1.Y = yb - 160
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 240
			P1.Y = yb - 240
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 320
			P1.Y = yb - 320
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 400
			P1.Y = yb - 400
			drawLine(P1, lineCol, imgXFull)

			drawLine = FuncDrawLiner(lineCol, imgXFull)
			yb = 440
			P0 = image.Point{xb + 0, yb - 0}
			drawLine(P0, lineCol, imgXFull)

			P1 = image.Point{xb + 80, yb - 40}
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 160
			P1.Y = yb - 90
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 240
			P1.Y = yb - 120
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 320
			P1.Y = yb - 300
			drawLine(P1, lineCol, imgXFull)
			P1.X = xb + 400
			P1.Y = yb - 310
			drawLine(P1, lineCol, imgXFull)

		}

		SaveImageToDatastore(w, r, imgXFull, "chart1")
		img := GetImageFromDatastore(w, r, "chart1")

		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)

		// end case

	}
}

func init() {
	http.HandleFunc("/image/analyze", util_appengine.Adapter(imageAnalyze))
	http.HandleFunc("/image/draw-lines-example", util_appengine.Adapter(drawLinesOverGrid))
}
