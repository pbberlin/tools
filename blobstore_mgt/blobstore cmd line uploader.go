package blobstore_mgt

// http://stackoverflow.com/questions/20144066/is-it-possible-to-store-arbitrary-data-in-gae-golang-blobstore

import (
	"bytes"
	"net/http"
	"mime/multipart"

	"appengine"
	"appengine/blobstore"
	"appengine/urlfetch"
)

const SampleData = `foo,bar,spam,eggs`


func blobAutoUploadTest(w http.ResponseWriter, r *http.Request) {


	c := appengine.NewContext(r)

	// First you need to create the upload URL:
	upl_url, err := blobstore.UploadURL(c, "/blob/auto-upload", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%s", err)
		return
	}


	// Now you can prepare a form that you will submit to that URL.
	var b1 bytes.Buffer
	fw := multipart.NewWriter(&b1)
	// Do not change the form field, it must be "file"!
	// You are free to change the filename though, it will be stored 
	// in the BlobInfo.
	file, err := fw.CreateFormFile("file", "filename-for-blobinfo-example.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%s", err)
		return
	}
	if _, err = file.Write([]byte(SampleData)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%s", err)
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	fw.Close()




	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", upl_url.String(), &b1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%s", err)
		return
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", fw.FormDataContentType())



	// Now submit the request.
	client := urlfetch.Client(c)
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%s", err)
		return
	}
	c.Infof("Autouploader Status - %v", res.Status)

	// Check the response status, it should be whatever you return in the `/upload` handler.
	if res.StatusCode != http.StatusCreated   &&   res.StatusCode != http.StatusOK{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("bad status: %s", res.Status)
		return
	}


	// Everything went fine.
	w.WriteHeader(res.StatusCode)
}


func blobAutoUpload(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// Here we just checked that the upload went through as expected.
	if _, _, err := blobstore.ParseUpload(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%s", err)
		return
	}
	// Everything seems fine. Signal the other handler using the status code.
	w.WriteHeader(http.StatusCreated)
}	


func init() {
	http.HandleFunc("/blob/auto-upload-test", blobAutoUploadTest)
	http.HandleFunc("/blob/auto-upload", blobAutoUpload)
}