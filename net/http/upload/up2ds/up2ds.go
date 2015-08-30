package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pbberlin/tools/net/http/upload"
)

var urlUp []string

func init() {
	urlUp = []string{
		"http://localhost:8085" + upload.UrlUploadReceive,
		"https://libertarian-islands.appspot.com" + upload.UrlUploadReceive,
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "up2ds"
	app.Usage = "Upload to datastore - local zip files are http posted to remote appengine handler."
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "file, f",
			Value:  "public.zip",
			Usage:  "filename to upload",
			EnvVar: "ENV_VAR_UP2DS_FILENAME_1,ENV_VAR_UP2DS_FILENAME_2",
		},
		cli.StringFlag{
			Name:   "live, l",
			Value:  "",
			Usage:  "'live' to upload to live, default is test",
			EnvVar: "ENV_VAR_UP2DS_UPLOAD_LIVE",
		},
		cli.StringFlag{
			Name:  "mount, m",
			Value: "mnt02",
			Usage: "the mount name for the dsfs",
		},
	}

	app.Action = func(c *cli.Context) {

		args := c.Args()
		if len(args) > 0 {
			fmt.Printf("  %v args given:\n", len(args))
			for k, v := range args {
				fmt.Printf("%4v %-v\n", k, v)
			}
		} else {
			fmt.Printf("  no args given\n")
		}

		fn := c.String("file")
		if fn == "" {
			fmt.Printf("  cannot upload empty file\n")
			return
		} else {
			fmt.Printf("  uploading... %v\n", fn)
		}

		url := urlUp[0]
		dest := c.String("live")
		if dest == "live" {
			url = urlUp[1]
			fmt.Printf("  uploading to live %v\n", fn)
		}

		mount := "mnt02"
		if c.String("mount") != "" {
			mount = c.String("mount")
		}
		fmt.Printf("  mount name %q\n", mount)

		uploadFile(fn, mount, url)

	}

	app.Run(os.Args)
}

func uploadFile(filename, mountName, effUpUrl string) {

	log.SetFlags(log.Lshortfile)

	curdir, _ := os.Getwd()
	filePath := filepath.Join(curdir, filename)

	extraParams := map[string]string{
		"getparam1": "val1",
		"mountname": mountName,
	}

	request, err := upload.CreateFilePostRequest(effUpUrl, "filefield", filePath, extraParams)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	fmt.Printf("    sending %v\n", filePath)
	fmt.Printf("     ... to %v\n", effUpUrl)
	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	} else {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			fmt.Println("\tHdr: ", k, v)
		}
		fmt.Printf("    status: %v\n", resp.StatusCode)

		bod := string(bts)
		bods := strings.Split(bod, "<span class='body'></span>")
		if len(bods) == 3 {
			bod = bods[1]
		}

		fmt.Printf("    body:   %s\n", bod)
	}

}
