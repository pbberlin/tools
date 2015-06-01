package charting

/*


import (
	resizer "github.com/nfnt/resize"
	"image"
	"image/png"
	"net/http"
	"strconv"
	"strings"
)

func resize(picture image.Image, width uint, height uint) image.Image {
	return resizer.Resize(width, height, picture, resizer.Bilinear)
}

func crop(picture image.Image, width uint, height uint) image.Image {
	pictureBounds := picture.Bounds()
	pictureCenter := pictureBounds.Max.Sub(pictureBounds.Min)

	croppedArea := image.Rectangle{
		pictureCenter.Sub(image.Point{int(width / 2), int(height / 2)}),
		pictureCenter.Add(image.Point{int(width / 2), int(height / 2)}),
	}

	return picture.(*image.Paletted).SubImage(croppedArea)
}

func parseSize(size string) (uint, uint) {
	segments := strings.Split(size, "x")

	width, err := strconv.Atoi(segments[0])
	if err != nil {
		width = 0
	}
	height, err := strconv.Atoi(segments[1])
	if err != nil {
		height = 0
	}

	return uint(width), uint(height)
}

func imageResize(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	// Get the picture from the given url
	if query["image"] == nil {
		return
	}

	imageUrl := query["image"][0]
	response, err := http.Get(imageUrl)
	if err != nil {
		return
	}

	picture, _, err := image.Decode(response.Body)
	if err != nil {
		return
	}

	// Decide what to do with the picture
	var operation string

	if query["type"] == nil {
		operation = "resize"
	} else {
		operation = query["type"][0]
	}

	// Get the size of the picture
	var width, height uint

	if query["size"] != nil {
		width, height = parseSize(query["size"][0])
	} else {
		width, height = 0, 0
	}

	// Modify the picture
	if operation == "resize" {
		picture = resize(picture, width, height)
	} else if operation == "crop" {
		picture = crop(picture, height, width)
	}

	// Send the picture as a PNG
	png.Encode(writer, picture)
}

func init() {
	http.HandleFunc("/image-resize", imageResize)
}


*/