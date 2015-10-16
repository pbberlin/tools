// hugoprep is a utility to customize a hugo-config.toml,
// to statify, to replace strings, to zip-compress,
// to upload.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/pbberlin/tools/net/http/upload"
	"github.com/pbberlin/tools/os/osutilpb"
)

var urlUp []string

const dirStat = "cnt_statified"
const dirStatPref = "./cnt_statified"

func init() {
	urlUp = []string{
		"http://localhost:8085" + upload.UrlUploadReceive,
		"https://libertarian-islands.appspot.com" + upload.UrlUploadReceive,
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "hugoprep"
	app.Usage = "Prepare a core hugo installation for multiple sites, with distinct config and directories."
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "name, n",
			Value:  "tec-news",
			Usage:  "The site name with all its JSON settings",
			EnvVar: "ENV_VAR_HUGOPREP_NAME_1,ENV_VAR_HUGOPREP_NAME_2",
		},
		cli.StringFlag{
			Name:  "arg2, a2",
			Value: "arg2default",
			Usage: "some second unused arg",
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

		name := c.String("name")
		if name == "" {
			fmt.Printf("  need the name of a site config\n")
			os.Exit(1)
		} else {
			// fmt.Printf("  preparing config... %v\n", name)
		}

		prepareConfigToml(name, "arg2")

	}

	app.Run(os.Args)
}

func prepareConfigToml(name, arg2 string) {

	log.SetFlags(log.Lshortfile)

	configTpl := template.Must(template.ParseFiles("config.tmpl"))

	bts, err := ioutil.ReadFile(name + ".json")
	if err != nil || len(bts) < 10 {
		log.Printf("did not find or too tiny file %v - %v\n", name+".json", err)
		os.Exit(1)
	}
	var data1 map[string]interface{}
	err = json.Unmarshal(bts, &data1)
	if err != nil {
		log.Printf("unmarshalling failed %v - %v\n", name+".json", err)
		os.Exit(1)
	}

	// lg(stringspb.IndentedDumpBytes(data1))
	if len(data1) < 2 {
		log.Printf("unmarshalling  yielded too fee settings for %v - %v\n", name+".json")
		os.Exit(1)
	}

	log.Printf("config data successfully unmarshalled")

	// sweep previous
	err = os.Remove("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("no previous file config.toml to delete")
		} else {
			log.Printf("deletion failed %v\n", err)
			os.Exit(1)
		}
	} else {
		log.Printf("removed previous %q\n", "config.toml")
	}

	f, err := os.Create("config.toml")
	if err != nil {
		log.Printf("could not open file config.toml: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	err = configTpl.Execute(f, data1)
	if err != nil {
		log.Printf("tmpl execution failed: %v\n", err)
		os.Exit(1)
	}

	log.Printf("=========fine-new-config.toml=======\n")

	// if ok := osutilpb.ExecCmdWithExitCode("cp", "config.toml", "../statify.toml"); !ok {
	// 	return
	// }
	// log.Printf("copied to statify.toml...\n")

	// os.Chdir("./..")
	// curdir, _ := os.Getwd()
	// log.Printf("curdir now %v\n", curdir)

	if ok := osutilpb.ExecCmdWithExitCode("rm", "-rf", dirStatPref+"/*"); !ok {
		return
	}
	log.Printf("cnt_statify cleaned up...\n")

	// hugo    --destination="cnt_statified"  --config="statify.toml"  --disableRSS=true  --disableSitemap=true
	if ok := osutilpb.ExecCmdWithExitCode("hugo",
		`--destination=`+dirStat,
		`--config=statify.toml`,
		`--disableRSS=true`,
		`--disableSitemap=true`); !ok {
		return
	}
	log.Printf("statified...\n")

	paths1 := []string{}

	//
	//
	// now the search and replace
	type WalkFunc func(path string, info os.FileInfo, err error) error
	fc := func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() || true {
			paths1 = append(paths1, path)
		}

		if ext := filepath.Ext(path); ext == ".html" {

			bts, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("could not read file: %v\n", path, err)
				return err
			}
			bef := len(bts)
			bts = bytes.Replace(bts, []byte("http://localhost:1313"), []byte("/"), -1)
			// bts = bytes.Replace(bts, []byte("&copy;2015 Peter Buchmann"), []byte("&copy;2015"), -1)
			aft := len(bts)
			err = ioutil.WriteFile(path, bts, 0777)
			if err != nil {
				log.Printf("could not write file: %v\n", path, err)
				return err
			}

			//
			//
			log.Printf("\t %-55v  %5v  %5v\n", path, bef, aft)
		}
		return nil
	}

	filepath.Walk(dirStatPref, fc)
	log.Printf("replacements finished...\n")

	log.Printf("=========fine-new-statification=======\n")

	//
	err = os.Chdir(dirStatPref)
	if err != nil {
		log.Printf("could not change dir: %v\n", err)
		return
	}
	curdir, _ := os.Getwd()
	log.Printf("curdir now %v\n", curdir)

	paths2 := make([]string, 0, len(paths1))
	for _, path := range paths1 {
		if path == dirStatPref {
			continue
		}
		pathAfter := strings.TrimPrefix(path, dirStat)
		if strings.HasPrefix(pathAfter, "/") || strings.HasPrefix(pathAfter, "\\") {
			pathAfter = pathAfter[1:]
		}
		paths2 = append(paths2, pathAfter)
		// log.Printf("%-54v %-54v", path, pathAfter)
	}

	osutilpb.CreateZipFile(paths2, dirStat+".zip")
	log.Printf("=========zip-file-created=============\n")

	if tgs, ok := data1["targets"]; ok {
		log.Printf("\t found targets ...")
		if tgs1, ok := tgs.(map[string]interface{}); ok {
			log.Printf("\t   ... %v urls - map[string]interface{} \n", len(tgs1))

			for _, v := range tgs1 {
				v1 := v.(string)
				// curl --verbose -F mountname=mnt02 -F filefield=@cnt_statified.zip subdm.appspot.com/filesys/post-upload-receive
				//                -F, --form CONTENT  Specify HTTP multipart POST data

				// We could also use
				// req, err := upload.CreateFilePostRequest(effUpUrl, "filefield", filePath, extraParams)
				// and client.Do(req)

				log.Printf("\t trying upload %v\n", v1)
				if ok := osutilpb.ExecCmdWithExitCode("curl", "--verbose",
					"-F", "mountname=mnt02",
					"-F", "filefield=@"+dirStat+".zip",
					v1); !ok {
					return
				}

				// Prepend a protocol for url.Parse to behave as desired
				if !strings.HasPrefix(v1, "http://") && !strings.HasPrefix(v1, "https://") {
					v1 = "https://" + v1
				}
				url2, err := url.Parse(v1)
				if err != nil {
					log.Printf("could not parse url: %q - %v\n", v1, err)
					return
				}

				url2.Path = "/reset-memfs"
				log.Printf("\t trying reset memfs %v\n", url2.Host+url2.Path)
				if ok := osutilpb.ExecCmdWithExitCode("curl", "--verbose", url2.Host+url2.Path); !ok {
					return
				}

				url2.Path = "/tpl/reset"
				log.Printf("\t trying reset templates %v\n", url2.Host+url2.Path)
				if ok := osutilpb.ExecCmdWithExitCode("curl", "--verbose", url2.Host+url2.Path); !ok {
					return
				}

			}

		}

	}
	log.Printf("====== upload completed ==============\n")
	log.Printf("CTRC+C to exit\n")

	for {
		time.Sleep(100 * time.Millisecond)
		runtime.Gosched()
	}

}
