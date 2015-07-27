package upload

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// CreateFilePostRequest prepares a multipart
// file upload POST request
func CreateFilePostRequest(url, fileParamName, filePath string,
	extraParams map[string]string) (*http.Request, error) {

	// try opening the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

	formFilePart, err := wrtr.CreateFormFile(fileParamName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// write the file to the file header
	_, err = io.Copy(formFilePart, file)

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
