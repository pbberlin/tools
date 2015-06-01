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

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf("blobstore: "+format, args...)
}



// ParseUpload parses the synthetic POST request that your app gets from
// App Engine after a user's successful upload of blobs. Given the request,
// ParseUpload returns a map of the blobs received (keyed by HTML form
// element name) and other non-blob POST parameters.
func OfficialParseUpload(req *http.Request) (blobs map[string][]*BlobInfo, other url.Values, err error) {
	_, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, nil, err
	}
	boundary := params["boundary"]
	if boundary == "" {
		return nil, nil, errorf("did not find MIME multipart boundary")
	}

	blobs = make(map[string][]*BlobInfo)
	other = make(url.Values)

	mreader := multipart.NewReader(io.MultiReader(req.Body, strings.NewReader("\r\n\r\n")), boundary)
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
		formKey := params["name"]

		ctype, params, err = mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return nil, nil, err
		}
		bi.BlobKey = appengine.BlobKey(params["blob-key"])
		if ctype != "message/external-body" || bi.BlobKey == "" {
			if formKey != "" {
				slurp, serr := ioutil.ReadAll(part)
				if serr != nil {
					return nil, nil, errorf("error reading %q MIME part", formKey)
				}
				other[formKey] = append(other[formKey], string(slurp))
			}
			continue
		}

		// App Engine sends a MIME header as the body of each MIME part.
		tp := textproto.NewReader(bufio.NewReader(part))
		header, mimeerr := tp.ReadMIMEHeader()
		if mimeerr != nil {
			return nil, nil, mimeerr
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
			return nil, nil, errorf("expected to find an X-AppEngine-Upload-Creation header")
		}
		bi.CreationTime, err = time.Parse("2006-01-02 15:04:05.000000", createDate)
		if err != nil {
			return nil, nil, errorf("error parsing X-AppEngine-Upload-Creation: %s", err)
		}

		if hdr := header.Get("Content-MD5"); hdr != "" {
			md5, err := base64.URLEncoding.DecodeString(hdr)
			if err != nil {
				return nil, nil, errorf("bad Content-MD5 %q: %v", hdr, err)
			}
			bi.MD5 = string(md5)
		}

		// If the GCS object name was provided, record it.
		bi.ObjectName = header.Get("X-AppEngine-Cloud-Storage-Object")

		blobs[formKey] = append(blobs[formKey], bi)
	}
	return
}




