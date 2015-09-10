package fileserver

import (
	"bytes"
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/stringspb"
)

func GetDirContents(hostWithPrefix, dir string) ([]string, []string, *bytes.Buffer, error) {

	lg, lge := loghttp.Logger(nil, nil)
	_ = lg

	var b = new(bytes.Buffer)

	dirs := []string{}
	fils := []string{}

	// build url
	urlSubDirs, err := url.Parse(path.Join(hostWithPrefix, dir))
	lge(err)
	if err != nil {
		return dirs, fils, b, err
	}
	sd := urlSubDirs.String()
	sd = common.Directorify(sd)
	wpf(b, "requ subdirs from  %v", sd)

	// make req
	bsubdirs, effU, err := fetch.UrlGetter(nil, fetch.Options{URL: sd})
	lge(err)
	if err != nil {
		return dirs, fils, b, err
	}
	wpf(b, "got %s - %v", bsubdirs, effU)

	// parse json
	mpSubDir := []map[string]string{}
	err = json.Unmarshal(bsubdirs, &mpSubDir)
	lge(err)
	if err != nil {
		// lg("%s", bsubdirs)
		return dirs, fils, b, err
	}
	wpf(b, "json of subdir is %s", stringspb.IndentedDump(mpSubDir))

	for _, v := range mpSubDir {

		if dir, ok := v["path"]; ok {
			if strings.HasSuffix(dir, "/") {
				dirs = append(dirs, dir)
			} else {
				fils = append(fils, dir)
			}
		}

		if smod, ok := v["mod"]; ok {
			t, err := time.Parse(time.RFC1123Z, smod)
			lge(err)
			wpf(b, "age %-6.2v", time.Now().Sub(t).Hours())
		}

	}
	return dirs, fils, b, nil

}
