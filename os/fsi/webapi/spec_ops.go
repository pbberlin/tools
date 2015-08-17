package webapi

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"

	tt "html/template"

	"appengine"
	"appengine/memcache"
)

var str = `
	<style>
		span.b {
			display:inline-block;
			width: 150px;
			align:middle;
		}
	</style>
	<form 
		action='{{.Url}}?getparam1=val1' 
		method='post'
		style='line-height:32px;padding:8px;'
	>

		<span class='b' >MountName </span> 
		<input name='mountname' type='text' value='mntftch' length=5 /><br>

		<span class='b' >Path Prefix </span> 
		<input name='pathprefix' type='text' value='/dir1/dir2' length=25 /><br>

		<span class='b' > </span> 
		<input type='submit' value='submit1' accesskey='s' /><br>

	</form>

	`

var tplBase = tt.Must(tt.New("tplName01").Parse(str))

func deleteSubtree(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	err := r.ParseForm()
	lge(err)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Delete Subtree for curr FS"}))
	defer wpf(w, tplx.Foot)

	if r.Method == "POST" {
		wpf(w, "<pre>\n")
		defer wpf(w, "\n</pre>")

		mountPoint := dsfs.MountPointLast()
		if len(r.FormValue("mountname")) > 0 {
			mountPoint = r.FormValue("mountname")
		}
		lg("mount point is %v", mountPoint)

		pathPrefix := "impossible-value"
		if len(r.FormValue("pathprefix")) > 0 {
			pathPrefix = r.FormValue("pathprefix")
		}
		lg("pathprefix is %v", pathPrefix)

		fs := getFS(appengine.NewContext(r), mountPoint)
		lg("created fs %v-%v ", fs.Name(), fs.String())

		lg("removing %v/*... ", pathPrefix)
		err := fs.RemoveAll(pathPrefix)
		lge(err)

		errMc := memcache.Flush(appengine.NewContext(r))
		lge(errMc)

		if err == nil && errMc == nil {
			lg("success")
		}

	} else {
		tData := map[string]string{"Url": UriDeleteSubtree}
		err := tplBase.ExecuteTemplate(w, "tplName01", tData)
		lge(err)

	}

}

func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Delete all filesystem data"}))
	defer wpf(w, tplx.Foot)

	confirm := r.FormValue("confirm")
	if confirm != "yes" {
		wpf(w, "All dsfs contents are deletes. All memfs contents are deleted<br>\n")
		wpf(w, "Put a get param into the URL ?confirm - and set it to 'yes'<br>\n")
		return
	}

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	//
	//
	fs := dsfs.New(dsfs.AeContext(appengine.NewContext(r)))
	wpf(w, "dsfs:\n")
	msg, err := fs.DeleteAll()
	if err != nil {
		wpf(w, "err during delete %v\n", err)
	}
	wpf(w, msg)

	memMapFileSys = memfs.New()
	wpf(w, "\n")
	wpf(w, "memMapFs new")

	// cleanup must be manual
	osFileSys = osfs.New()

}
