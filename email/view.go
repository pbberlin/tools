package email

import (
	"net/http"

	"bytes"

	"appengine"

	// "github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_err"

	"fmt"
	"strings"
)

const keyLatest = "latestEmail"
const sepHeaderContent = "\r\n\r\n"

var sp func(format string, a ...interface{}) string = fmt.Sprintf

const reservoirSize = 100

type resEntry struct {
	when        string
	fn          string // filename
	contentType string
	b64Img      *string
}

var reservoirRevolver int = 0
var Images []resEntry = make([]resEntry, reservoirSize)

func view(w http.ResponseWriter, r *http.Request) {
	parseFurther(w, r, false)
}

//
func parseFurther(w http.ResponseWriter, r *http.Request, saveImages bool) {

	c := appengine.NewContext(r)

	b := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		w.Write(b.Bytes())
	}()

	// Get the item from the memcache
	wb1 := new(dsu.WrapBlob)
	ok := dsu.McacheGet(c, keyLatest, wb1)
	util_err.Err_http(w, r, ok, true)

	if ok {
		b.WriteString(sp("name %v\n", wb1.Name))
		b.WriteString(sp("S (boundary): %q\n", wb1.S))

		// dumps the entire body
		// b.WriteString(sp("B: %v\n", string(wb1.VByte)))

		// instead we split it by multipart mime
		vb := bytes.Split(wb1.VByte, []byte("--"+wb1.S))
		for i, v := range vb {
			h := ""  // header
			fn := "" // filename
			s := string(v)
			s = strings.Trim(s, "\r \n")
			ctype := ""

			b.WriteString(sp("\n___________mime boundary index %v___________\n", i))
			if strings.HasPrefix(s, "Content-Type: image/png;") ||
				strings.HasPrefix(s, "Content-Type: image/jpeg;") {

				if start := strings.Index(s, sepHeaderContent); start > 0 {
					h = s[:start]
					vh := strings.Split(h, "\r\n")
					for _, v := range vh {
						v := strings.TrimSpace(v)
						// b.WriteString("\t\t" + v + "\n")
						if strings.HasPrefix(v, "name=") {
							vv := strings.Split(v, "=")
							fn = pbstrings.LowerCasedUnderscored(vv[1])
						}
					}
					s = s[start+len(sepHeaderContent):]
					if posSemicol := strings.Index(h, ";"); posSemicol > 0 {
						ctype = h[0:posSemicol]
					}
				}
			}

			if ctype == "" {
				b.WriteString("unparseable: " + pbstrings.Ellipsoider(s, 400))
			} else {
				b.WriteString(sp("\n\tctype=%v\n\t------------", ctype))
				if fn != "" {
					b.WriteString(sp("\n\tfilename=%v\n\t------------", fn))
				}
				if saveImages {
					rE := resEntry{}
					rE.when = util.TimeMarker()
					rE.contentType = ctype
					rE.fn = fn
					rE.b64Img = &s
					Images[reservoirRevolver%reservoirSize] = rE
					reservoirRevolver++
					c.Infof("Put image into reservoir %v %v", fn, ctype)
				}
			}

		}

	}

}

func viewImages(w http.ResponseWriter, r *http.Request) {

	b := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		w.Write(b.Bytes())
	}()

	for idx, imgB64 := range Images {
		if imgB64.when == "" {
			b.WriteString(sp("image %v empty<br>", idx))
			continue
		}
		b.WriteString(sp("image #%v from %v - %v - %v<br>", idx, imgB64.when, imgB64.contentType, imgB64.fn))
		b.WriteString("<img width=200px src=\"")
		b.WriteString(`data:image/png;base64,`)
		b.WriteString(*imgB64.b64Img)
		b.WriteString("\"><br> ")

	}

}

func init() {
	http.HandleFunc("/email-view", view)
	http.HandleFunc("/email-images", viewImages)
}
