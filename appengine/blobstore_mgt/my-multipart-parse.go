package blobstore_mgt

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"net/url"
	"strconv"
	"strings"
	"time"

	"appengine"
)



// ParseUpload parses the synthetic POST request that your app gets from
// App Engine after a user's successful upload of blobs. 
func MyParseUpload(req *http.Request) (blobs map[string][]*BlobInfo, 
	other url.Values, err error) {
	
	_, params, errCT := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if errCT != nil {
		return nil, nil, errCT
	}

	boundary := params["boundary"]
	if boundary == "" {
		return nil, nil, fmt.Errorf("did not find MIME multipart boundary")
	}

	blobs = make(map[string][]*BlobInfo)
	other = make(url.Values)

	mreader := multipart.NewReader(io.MultiReader(req.Body, strings.NewReader("\r\n\r\n")), boundary)
	cntr := 0
	for {
		part, perr := mreader.NextPart()
		if perr == io.EOF {
			break
		}
		if perr != nil {
			return nil, nil, errorf("error reading next mime part with boundary %q (len=%d): %v",
				boundary, len(boundary), perr)
		}

		bi := &BlobInfo{}
		ctype, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
		if err != nil {
			return nil, nil, err
		}
		bi.Filename = params["filename"]
		pFile := params["post_field_file"] // WRONG
		pFile = params["name"]

		ctype, params, err = mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return nil, nil, err
		}
		

		bi.BlobKey = appengine.BlobKey(params["blob-key"])
		if ctype != "message/external-body" || bi.BlobKey == "" {
			if pFile != "" {
				slurp, serr := ioutil.ReadAll(part)
				if serr != nil {
					return nil, nil, errorf("error reading %q MIME part", pFile)
				}
				other[pFile] = append(other[pFile], string(slurp))
			}
			continue
		}


		// App Engine sends a MIME header as the body of each MIME part.
		tp := textproto.NewReader(bufio.NewReader(part))


		header, mimeerr := tp.ReadMIMEHeader()
		if mimeerr != nil {
			s := mimeerr.Error()
			if strings.HasPrefix(s, "malformed MIME header") {
				return nil, nil, fmt.Errorf("'malformed'  %q", mimeerr)
			} else {
				return nil, nil, fmt.Errorf("error reading again %q", mimeerr)
			}
		}

		bi.Size, err = strconv.ParseInt(header.Get("Content-Length"), 10, 64)
		if err != nil {
			return nil, nil, err
		}
		bi.ContentType = header.Get("Content-Type")

		// Parse the time from the MIME header like:
		// X-AppEngine-Upload-Creation: 2011-03-15 21:38:34.712136
		createDate := header.Get("X-AppEngine-Upload-Creation")
		if createDate == "" {
			return nil, nil, fmt.Errorf("expected to find an X-AppEngine-Upload-Creation header")
		}
		bi.CreationTime, err = time.Parse("2006-01-02 15:04:05.000000", createDate)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing X-AppEngine-Upload-Creation: %s", err)
		}

		if hdr := header.Get("Content-MD5"); hdr != "" {
			md5, err := base64.URLEncoding.DecodeString(hdr)
			if err != nil {
				return nil, nil, fmt.Errorf("bad Content-MD5 %q: %v", hdr, err)
			}
			bi.MD5 = string(md5)
		}


		if pFile != "" {
			// slurp, serr := ioutil.ReadAll(part)
			// if serr != nil {
			// 	return nil, nil, fmt.Errorf("error reading all from %q", pFile)
			// }
			// other[pFile] = append(other[pFile], string(slurp))
		}

		blobs[pFile] = append(blobs[pFile], bi)
		cntr++
	}
	return
}




