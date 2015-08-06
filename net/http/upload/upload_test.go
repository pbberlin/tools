package upload

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	curdir, _ := os.Getwd()
	filePath := curdir + "/test.zip"

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Pete",
		"description": "A zip file - containing dirs and files",
	}

	urlUp := "https://google.com/upload"
	urlUp = "http://localhost:8085/blob2/zipupload"
	// urlUp = "https://libertarian-islands.appspot.com/blob2/zipupload"

	request, err := CreateFilePostRequest2(
		urlUp, "filefield", filePath, extraParams)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	log.Printf("Sending req ...")
	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	} else {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		for k, v := range resp.Header {
			fmt.Println("\tHdr: ", k, v)
		}
		fmt.Printf("%s", bts)
	}
}
