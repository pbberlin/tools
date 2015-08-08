package upload

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"

	tt "html/template"
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
		enctype='multipart/form-data'
		style='line-height:32px;padding:8px;'
	>
		<span class='b' >Datei</span> 
		<input 
			type='file'  
			name='filefield'
			acceptxx='application/zip, image/gif, image/jpeg' 
			accept='application/zip' 
			accesskey='w' 
		/><br>

		<span class='b' >MountName </span> 
		<input name='mountname' type='text' value='mnt01' length=5 /><br>


		<span class='b' > </span> 
		<input type='submit' value='submit1' accesskey='s' /><br>

	</form>

	`
var tplBase = tt.Must(tt.New("tplName01").Parse(str))

func sendUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, _ := loghttp.Logger(w, r)
	// c := appengine.NewContext(r)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Post an Upload"}))
	defer wpf(w, tplx.Foot)

	tData := map[string]string{"Url": UrlUploadReceive}
	err := tplBase.ExecuteTemplate(w, "tplName01", tData)
	if err != nil {
		lg("tpl did not compile: %v", err)
	}

}
