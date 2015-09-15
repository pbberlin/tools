package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/net/http/tplx"

	tt "html/template"
)

const form = `
	<style> .ib { display:inline-block; }</style>



	<form>
		<div style='margin:8px;'>
			<span class='ib' style='width:40px'>Count </span>
			<input id='inp1' name="cnt"                      size="3"  value="2"><br/>

			<span class='ib' style='width:40px'>URL </span>
			<input id='inp2' name="{{.fieldname}}"           size="120"  value="{{.val}}"><br/>
			
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
			//$( "#inp2" ).text( "focusout fired: " + focus + "x" );
			$( "#inp2" ).val(  $.trim( $(this).text() )   );
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

		fullURL := fmt.Sprintf("https://%s%s?%s=%s&cnt=%s", r.Host, UriFetchSimilar, routes.URLParamKey, rURL, r.FormValue("cnt"))
		lg("lo - sending to URL:    %v\n", fullURL)

		fo := fetch.Options{}
		fo.URL = fullURL
		bts, inf, err := fetch.UrlGetter(r, fo)
		_ = inf
		lg(err)
		if err != nil {
			return
		}

		if len(bts) == 0 {
			lg("empty bts")
			return
		}

		var mp map[string][]byte
		err = json.Unmarshal(bts, &mp)
		lg(err)
		if err != nil {
			lg("%s", bts)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, ok := mp["msg"]; ok {
			w.Write(mp["msg"])
		}

		for k, v := range mp {
			if k != "msg" {
				wpf(w, "<br><br>%s:\n", k)
				if true {
					wpf(w, "len %v", len(v))
				} else {
					wpf(w, "%s", html.EscapeString(string(v)))
				}
			}
		}

	}

}
