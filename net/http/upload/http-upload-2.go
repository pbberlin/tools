package upload

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// CreateFilePostRequest2 prepares a multipart
// file upload POST request
func CreateFilePostRequest2(url, fileParamName, filePath string,
	extraParams map[string]string) (*http.Request, error) {

	// try opening the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// first 512 bytes are used to evaluate mime type
	first512 := make([]byte, 512)
	file.Read(first512)
	file.Seek(0, 0)

	// prepare request body buffer
	body := &bytes.Buffer{}

	// create a multipart file header with the given param name and file
	wrtr := multipart.NewWriter(body)

	// append extra params
	for param, val := range extraParams {
		err = wrtr.WriteField(param, val)
		if err != nil {
			return nil, err
		}
	}

	h := make(textproto.MIMEHeader)
	cd := fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		escapeQuotes(fileParamName), escapeQuotes(filepath.Base(filePath)))
	h.Set("Content-Disposition", cd)
	ct := http.DetectContentType(first512)
	log.Println("detected", ct)
	h.Set("Content-Type", ct)

	formfilePart, err := wrtr.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// write the file to the file header
	num, err := io.Copy(formfilePart, file)
	log.Println("copied", num, "bytes")

	// the writer must be closed to finalize the file entry
	err = wrtr.Close()
	if err != nil {
		return nil, err
	}

	// finally create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// set type & boundary
	req.Header.Set("Content-Type", wrtr.FormDataContentType())

	return req, err

}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
