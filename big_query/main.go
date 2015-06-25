package big_query

import (
	"net/http"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/tpl_html"
)

func ViewHTML(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	b1, ml := disLegend(w, r)
	_ = b1

	add, tplExec := tpl_html.FuncTplBuilder(w, r)

	add("n_html_title", "The Battle of Computer Languages", "")

	add("n_cont_0", tpl_html.PrefixLff+"chart_body", map[string]map[string]string{"legend": ml})
	add("tpl_legend", tpl_html.PrefixLff+"chart_body_embed01", "")

	add("n_cont_1", `<a 
			target='openhub'
			href='https://www.openhub.net/languages/compare?measure=loc_changed&percent=true&l0=-1&l1=golang&l2=php&l3=python&l4=ruby&l5=-1&commit=Update' 
			>Here is a good comparison</a>`, "")

	tplExec(w, r)

}

func init() {
	http.HandleFunc("/big-query/html", util_appengine.Adapter(ViewHTML))
	http.HandleFunc("/big-query/test-gob-codec", testGobDecodeEncode)
}
