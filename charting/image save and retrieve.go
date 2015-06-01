package charting

import (
	"appengine"
	"fmt"
	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/util_err"
	"image"

	"net/http"

	"bytes"
	"io"
	"strings"
)

func SaveImageToDatastore(w http.ResponseWriter, r *http.Request, i *image.RGBA, key string) string {

	s := conv.Rgba_img_to_base64_str(i)
	internalType := fmt.Sprintf("%T", i)
	//buffBytes, _     := StringToVByte(s)  // instead of []byte(s)
	key_combi, err := dsu.BufPut(w, r, dsu.WrapBlob{Name: key, VByte: []byte(s), S: internalType}, key)
	util_err.Err_http(w, r, err, false)

	return key_combi

}

func GetImageFromDatastore(w http.ResponseWriter, r *http.Request, key string) *image.RGBA {

	c := appengine.NewContext(r)

	dsObj, err := dsu.BufGet(w, r, "dsu.WrapBlob__"+key)
	util_err.Err_http(w, r, err, false)

	s := string(dsObj.VByte)

	img, whichFormat := conv.Base64_str_to_img(s)
	c.Infof("retrieved img from base64: format %v - subtype %T\n", whichFormat, img)

	i, ok := img.(*image.RGBA)
	util_err.Err_http(w, r, ok, false, "saved image needs to be reconstructible into a format png of subtype *image.RGBA")

	return i
}

// unused
// probably  efficient enough just to call
// var bEnc []byte = []byte(sEnc)
func StringToVByte(s string) (*bytes.Buffer, *bytes.Buffer) {

	bMsg := new(bytes.Buffer)

	bDst := new(bytes.Buffer)

	const chunksize = 20
	lb := make([]byte, chunksize) // loop buffer
	rdr := strings.NewReader(s)
	for {
		n1, err := rdr.Read(lb)
		if err == io.EOF {
			break
		}
		util_err.Err_log(err)
		if n1 < 1 {
			break
		}

		independentCopy := make([]byte, n1)
		copy(independentCopy, lb)
		n2, err := bDst.Write(independentCopy)
		util_err.Err_log(err)

		bMsg.WriteString(fmt.Sprintln("reading", n1, "bytes - writing", n2, "bytes: \n"))
		bMsg.WriteString(fmt.Sprint(" --", string(independentCopy), "--<br>\n"))
	}

	return bDst, bMsg

}

// based on bytes.Buffer and Writing into it
//   it's probably easier to just call s := string(m)
/*
func VByteToString( m []byte)( *bytes.Buffer, *bytes.Buffer){

	bRet := new(bytes.Buffer)
	bMsg := new(bytes.Buffer)

	n,err := bRet.Write( m )
	bMsg.WriteString( fmt.Sprintln("reading" , n, "bytes\n"))
	util_err.Err_log(err)

	return bRet,bMsg
}
*/
