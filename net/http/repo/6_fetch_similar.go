package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"appengine"

	"github.com/pbberlin/tools/appengine/instance_mgt"
	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/stringspb"

	"html"
	tt "html/template"
)

// FetchSimilar is an extended version of Fetch
// It is uses a DirTree of crawled *links*, not actual files.
// As it moves up the DOM, it crawls every document for additional links.
// It first moves up to find similar URLs on the same depth
//                        /\
//          /\           /  \
//    /\   /  \         /    \
// It then moves up the ladder again - to accept higher URLs
//                        /\
//          /\
//    /\
func FetchSimilar(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	wpf(b, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Find similar HTML URLs"}))
	defer wpf(b, tplx.Foot)

	wpf(b, "<pre>")
	defer wpf(b, "</pre>")

	fs1 := GetFS(appengine.NewContext(r))

	err := r.ParseForm()
	lg(err)

	surl := r.FormValue(routes.URLParamKey)
	ourl, err := fetch.URLFromString(surl)
	lg(err)
	if err != nil {
		return
	}
	if ourl.Host == "" {
		lg("host is empty (%v)", surl)
		return
	}

	srcDepth := strings.Count(ourl.Path, "/")

	cmd := FetchCommand{}
	cmd.Host = ourl.Host
	cmd.SearchPrefix = ourl.Path
	cmd = addDefaults(w, r, cmd)

	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, EndPoint: true}
	fnDigest := path.Join(docRoot, cmd.Host, "digest2.json")
	loadDigest(w, r, lg, fs1, fnDigest, dirTree) // previous
	lg("dirtree 400 chars is %v end of dirtree\n", stringspb.ToLen(dirTree.String(), 400))

	err = fetchCrawlSave(w, r, lg, dirTree, fs1, path.Join(cmd.Host, ourl.Path))
	lg(err)
	if err != nil {
		return
	}

	var treePath string
	treePath = "/news"
	treePath = "/blogs"
	treePath = "/blogs/freeexchange"
	treePath = "/news/europe"
	treePath = path.Dir(ourl.Path)

	opt := LevelWiseDeeperOptions{}
	opt.Rump = treePath
	opt.ExcludeDir = "/news/americas"
	opt.ExcludeDir = "/blogs/buttonwood"
	opt.ExcludeDir = "/something-impossible"
	opt.MinDepthDiff = 1
	opt.MaxDepthDiff = 1
	opt.CondenseTrailingDirs = cmd.CondenseTrailingDirs
	opt.MaxNumber = 2266
	opt.MaxNumber = cmd.DesiredNumber + 1  // one more for "self"
	opt.MaxNumber = cmd.DesiredNumber + 10 // collect more, 'cause we filter out those too old later

	var subtree *DirTree
	links := []FullArticle{}

MarkOuter:
	for j := 0; j < srcDepth; j++ {
		treePath = path.Dir(ourl.Path)
	MarkInner:
		for i := 1; i < 22; i++ {

			subtree, treePath = DiveToDeepestMatch(dirTree, treePath)

			lg("\nLooking from height %v to level %v  - %v", srcDepth-i, srcDepth-j, treePath)

			err = fetchCrawlSave(w, r, lg, dirTree, fs1, path.Join(cmd.Host, treePath))
			lg(err)
			if err != nil {
				return
			}

			if subtree == nil {
				lg("\n#%v treePath %q ; subtree is nil", i, treePath)
			} else {
				// lg("\n#%v treePath %q ; subtree exists", i, treePath)

				opt.Rump = treePath
				opt.MinDepthDiff = i - j
				opt.MaxDepthDiff = i - j
				lvlLinks := LevelWiseDeeper(nil, nil, subtree, opt)
				links = append(links, lvlLinks...)
				for _, art := range lvlLinks {
					lg("#%v     fnd %v", i, stringspb.ToLen(art.Url, 100))
				}

				if len(links) >= opt.MaxNumber {
					lg("found enough")
					break MarkOuter
				}

				pathPrev := treePath
				treePath = path.Dir(treePath)
				// lg("#%v  bef %v - aft %v", i, pathPrev, treePath)

				if pathPrev == "." && treePath == "." ||
					pathPrev == "/" && treePath == "/" ||
					pathPrev == "" && treePath == "." {
					lg("break to innner")
					break MarkInner
				}
			}

		}
	}

	lg("\nNow reading/fetching actual similar files - not just the links")
	//
	tried := 0
	selecteds := []FullArticle{}

	for i, art := range links {

		tried = i + 1

		if art.Url == ourl.Path {
			lg("skipping self")
			continue
		}

		useExisting := false

		semanticUri := condenseTrailingDir(art.Url, cmd.CondenseTrailingDirs)
		p := path.Join(docRoot, cmd.Host, semanticUri)
		lg("reading  %v", p)
		f, err := fs1.Open(p)
		if err == nil {
			defer f.Close()
			fi, err := f.Stat()
			if err == nil {
				if fi.ModTime().After(time.Now().Add(-10 * time.Hour)) {
					lg(" using file")
					art.Mod = fi.ModTime()
					bts, err := ioutil.ReadAll(f)
					lg(err)
					art.Body = bts
					selecteds = append(selecteds, art)
					useExisting = true
				}

			}
		}

		if !useExisting {
			surl := path.Join(cmd.Host, art.Url)
			lg("fetching %v", surl)
			bts, inf, err := fetch.UrlGetter(r, fetch.Options{URL: surl, RedirectHandling: 1})
			lg(err)

			if inf.Mod.IsZero() {
				inf.Mod = time.Now().Add(-75 * time.Minute)
			}

			lg("saving   %v", p)
			dir := path.Dir(p)
			err = fs1.MkdirAll(dir, 0755)
			lg(err)
			err = fs1.Chtimes(dir, time.Now(), time.Now())
			lg(err)
			err = fs1.WriteFile(p, bts, 0644)
			lg(err)
			err = fs1.Chtimes(p, inf.Mod, inf.Mod)
			lg(err)

			if inf.Mod.After(time.Now().Add(-10 * time.Hour)) {
				lg(" using fetched")
				art.Mod = inf.Mod
				art.Body = bts
				selecteds = append(selecteds, art)
			}

		}

		if len(selecteds) > 3 {
			break
		}

		if tried > 4 {
			break
		}

	}

	lg("tried %v to find %v new similars", tried, len(selecteds))

	mp := map[string][]byte{}
	mp["msg"] = b.Bytes()
	mp["semanticSelf"] = []byte(condenseTrailingDir(ourl.Path, cmd.CondenseTrailingDirs))

	for i, v := range selecteds {
		mp["url__"+spf("%02v", i)] = []byte(v.Url)
		mp["mod__"+spf("%02v", i)] = []byte(v.Mod.Format(http.TimeFormat))
		mp["bod__"+spf("%02v", i)] = v.Body
	}

	//
	smp, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		lg(b, "marshalling mp to []byte failed\n")
		return
	}

	r.Header.Set("X-Custom-Header-Counter", "nocounter")
	w.Header().Set("Content-Type", "application/json")
	w.Write(smp)

	b.Reset()             // this keeps the  buf pointer intact; outgoing defers are still heeded
	b = new(bytes.Buffer) // creates a *new* buf pointer; outgoing defers write into the *old* buf

	return

}

const form = `
	<style> .ib { display:inline-block; }</style>



	<form>
		<div style='margin:8px;'>
			<span class='ib' style='width:40px'>URL </span>
			<input id='inp1' name="{{.fieldname}}"           size="120"  value="{{.val}}"><br/>
			
			<span class='ib' style='width:40px' ></span> 
			<span class='ib' tabindex='11'> 
			www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html
			</span>
			<br/>

			<span class='ib' style='width:40px' ></span> 
			<span class='ib' tabindex='11'> 
			www.economist.com/news/britain/21663648-hard-times-hard-hats-making-britain-make-things-again-proving-difficult  
			</span>
			<br/>

			<span class='ib' style='width:40px' ></span> 
			<span class='ib' tabindex='12'> 
			www.economist.com/news/americas/21661804-gender-equality-good-economic-growth-girl-power  
			</span>
			<br/>

			<span class='ib' style='width:40px'> </span>
			<input type="submit" value="Get similar (shit+alt+f)" accesskey='f'>
		</div>
	</form>

	<script src="http://ajax.googleapis.com/ajax/libs/jquery/1/jquery.min.js" 
			type="text/javascript"></script>


	<script>
		var focus = 0,
		blur = 0;
		//focusout
		$( "span" ).focusin(function() {
			focus++;
			//$( "#inp1" ).text( "focusout fired: " + focus + "x" );
			$( "#inp1" ).val(  $.trim( $(this).text() )   );
			console.log("fired")
		});
	</script>	



	`

func fetchSimForm(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	// on live server => always use https
	if r.URL.Scheme != "https" && !util_appengine.IsLocalEnviron() {
		r.URL.Scheme = "https"
		r.URL.Host = r.Host
		lg("lo - redirect %v", r.URL.String())
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}

	err := r.ParseForm()
	lg(err)

	rURL := ""
	if r.FormValue(routes.URLParamKey) != "" {
		rURL = r.FormValue(routes.URLParamKey)
	}
	if len(rURL) == 0 {

		wpf(b, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Find similar HTML URLs"}))
		defer wpf(b, tplx.Foot)

		tm := map[string]string{
			"val":       "www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html",
			"fieldname": routes.URLParamKey,
		}
		tplForm := tt.Must(tt.New("tplName01").Parse(form))
		tplForm.Execute(b, tm)

	} else {

		ii := instance_mgt.Get(r)
		fullURL := fmt.Sprintf("https://%s%s?%s=%s", ii.PureHostname, UriFetchSimilar, routes.URLParamKey, rURL)
		lg("lo - sending to URL:    %v\n", fullURL)

		fo := fetch.Options{}
		fo.URL = fullURL
		bts, inf, err := fetch.UrlGetter(r, fo)
		if err != nil {
			lg("Requesting %v failed", inf)
			return
		}

		if len(bts) == 0 {
			lg("empty bts")
		} else {
			var mp map[string][]byte
			err = json.Unmarshal(bts, &mp)
			lg(err)
			if err != nil {
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if _, ok := mp["msg"]; ok {
				w.Write(mp["msg"])
			}

			for k, v := range mp {
				if k != "msg" {
					lg("<br><br>%s:<br>\n", k)
					lg("%s", html.EscapeString(string(v)))
					// lg("%s\n", bts)
				}
			}

		}

	}

}
