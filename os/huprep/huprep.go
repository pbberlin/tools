package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/pbberlin/tools/net/http/upload"
	"github.com/pbberlin/tools/os/osutilpb"
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
	var data1 map[string]string
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

	if ok := osutilpb.ExecCmdWithExitCode("cp", "config.toml", "../statify.toml"); !ok {
		return
	}
	log.Printf("copied to statify.toml...\n")

	// hugo    --destination="cnt_statified"  --config="statify.toml"  --disableRSS=true  --disableSitemap=true
	os.Chdir("./..")
	curdir, _ := os.Getwd()
	log.Printf("curdir now %v\n", curdir)

	// rm -rf  ./cnt_statified/*
	if ok := osutilpb.ExecCmdWithExitCode("rm", "-rf", "./cnt_statified/*"); !ok {
		return
	}
	log.Printf("cnt_statify cleaned up...\n")

	if ok := osutilpb.ExecCmdWithExitCode("hugo",
		`--destination=cnt_statified`,
		`--config=statify.toml`,
		`--disableRSS=true`,
		`--disableSitemap=true`); !ok {
		return
	}
	log.Printf("statified...\n")

	paths := []string{}

	//
	//
	// now the search and replace
	type WalkFunc func(path string, info os.FileInfo, err error) error
	fc := func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			paths = append(paths, path)
		}

		if ext := filepath.Ext(path); ext == ".html" {

			bts, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("could not read file: %v\n", path, err)
				return err
			}
			bef := len(bts)
			bts = bytes.Replace(bts, []byte("http://localhost:1313"), []byte("/"), -1)
			bts = bytes.Replace(bts, []byte("&copy;2015 Peter Buchmann"), []byte("&copy;2015"), -1)
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

	filepath.Walk("./cnt_statified", fc)
	log.Printf("replacements finished...\n")

	log.Printf("=========fine-new-statification=======\n")

	//
	os.Chdir("./cnt_statified")
	curdir, _ = os.Getwd()
	log.Printf("curdir now %v\n", curdir)

	for idx, path := range paths {
		pathAfter := strings.TrimPrefix(path, "cnt_statified")
		if strings.HasPrefix(pathAfter, "/") || strings.HasPrefix(pathAfter, "\\") {
			pathAfter = pathAfter[1:]
		}
		paths[idx] = pathAfter
		// log.Printf("%-54v %-54v", path, pathAfter)
	}

	osutilpb.CreateZipFile(paths, "cnt_statified.zip")
	log.Printf("=========zip-file-created=============\n")

}
