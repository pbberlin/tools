package conv


import (
	"encoding/base64"
	"fmt"
	"io"
	"image"
	// "image/jpeg"
	"image/png"
	"strings"
	"bytes"
	"github.com/pbberlin/tools/util_err"
	

)




func Base64_str_to_img( base64_img string )(img image.Image, format string){
	
	pos := base64HeaderPosition(base64_img)	
	
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader( base64_img[pos:] ))
	img, format , err := image.Decode(reader)
	util_err.Err_log(err)
	return
}



// convert image to string, prepend mime header
//		inspiration_1 http://stackoverflow.com/questions/22945486/golang-converting-image-image-to-byte
//		inspiration_2 https://github.com/polds/imgbase64/blob/master/images.go
//
// mime type is always image/png
// because we only accept *image.RGBA and use image/png Encoder
func Rgba_img_to_base64_str( img *image.RGBA) string {

	// img to bytes
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img) ; util_err.Err_log(err)
	imgBytes := buf.Bytes()

	// binary bytes to base64 bytes
	e64 := base64.StdEncoding
	maxEncLen := e64.EncodedLen(  len(imgBytes)  )
	imgEnc    := make( []byte, maxEncLen )
	e64.Encode(imgEnc, imgBytes)

	// base64 bytes to string
	mimeType := "image/png" 
	return fmt.Sprintf("data:%s;base64,%s", mimeType, imgEnc)
	
}

func f____________________________(){}


// we want to find exact position of comma in
//     "data:image/png;base64,"
//   (if the header is present at all)
// string s could be large but according to 
//   https://groups.google.com/forum/#!topic/golang-nuts/AdO_d4E_x6k
// the runtime makes it a pointer on its own
func base64HeaderPosition(s string)(pos int){

	test_for_header := "data:image/"
	headerLen := len(test_for_header)

	if s[:headerLen] == test_for_header {  // header is present
		pos = strings.Index(s,",")				// get comma pos
		pos++
	}

	return
	
}


// we want to extract the 'image/png' from 
// 	"data:image/png;base64,..."
//func MimeFromBase64(b *bytes.Buffer)(mime string){
func MimeFromBase64(b io.Reader)(mime string){

	// to avoid huge string allocation
	//   we read the first 100 bytes
	b1 := make([]byte,100)
	_,err := b.Read(b1)	
	util_err.Err_panic(err)
	
	s := string(b1)

	pos := base64HeaderPosition(s)	
	
	if pos > 0 {
		tmp1 := s[:pos]		
		tmp2 := strings.Split(tmp1,";")
		tmp3 := strings.Split(tmp2[0],":")
		mime = tmp3[1]
	}

	return
	
}