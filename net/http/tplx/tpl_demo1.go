package tplx

import (
	"html"
	tt "html/template"
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
)

// preconceive all embedded templates:
const T0 = `
	T0	<br>
		<span style='color:#aaf; line-height:200%;display:inline-block; margin:8px; margin-left:120px; border: 1px solid #aaf'>
				{{template "T1" .}}
		</span>
		<hr>
`

const T1 = `
{{define "T1"}}
	T1 <br>
		--{{if .key1}}{{.key1}}{{else}}no dyn data{{end}}--<br>
		<span style='min-width:200px;color:#faa; display:inline-block; margin:8px; margin-left:120px; border: 1px solid #faa'>
			{{template "T2" .key2 }}
		</span>
	
{{end}}
`

const iterOver = `{{ $mapOrArray := . }} 
{{range $index, $element := $mapOrArray }}
   <li><strong>$index</strong>: $element </li>
{{end}}`

const treatFirstIterDifferent = `{{if $index}},x{{end}}`

func templatesCompileDemo(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	w.Header().Set("Content-Type", "text/html")

	funcMap := tt.FuncMap{
		"unescape": html.UnescapeString,
		"escape":   html.EscapeString,
	}

	var t_base *tt.Template
	var err error = nil

	// creating T0 - naming it - adding func map
	t_base = tt.Must(tt.New("str_T0_outmost").Funcs(funcMap).Parse(T0))
	loghttp.E(w, r, err, false)

	// adding the definition of T1 - introducing reference to T2 - undefined yet
	t_base, err = t_base.Parse(T1) // definitions must appear at top level - but not at the start
	loghttp.E(w, r, err, false)

	// create two clones
	// now both containing T0 and T1
	tc_1, err := t_base.Clone()
	loghttp.E(w, r, err, false)
	tc_2, err := t_base.Clone()
	loghttp.E(w, r, err, false)

	// adding different T2 definitions
	s_if := "{{if .}}{{.}}{{else}}no dyn data{{end}}"
	tc_1, err = tc_1.Parse("{{define `T2`}}T2-A  <br>--" + s_if + "--  {{end}}")
	loghttp.E(w, r, err, false)
	tc_2, err = tc_2.Parse("{{define `T2`}}T2-B  <br>--" + s_if + "--  {{end}}")
	loghttp.E(w, r, err, false)

	// writing both clones to the response writer
	err = tc_1.ExecuteTemplate(w, "str_T0_outmost", nil)
	loghttp.E(w, r, err, false)

	// second clone is written with dynamic data on two levels
	dyndata := map[string]string{"key1": "dyn val 1", "key2": "dyn val 2"}
	err = tc_2.ExecuteTemplate(w, "str_T0_outmost", dyndata)
	loghttp.E(w, r, err, false)

	// Note: it is important to pass the DOT
	//		 {{template "T1" .}}
	//		 {{template "T2" .key2 }}
	//						 ^
	// otherwise "dyndata" can not be accessed by the inner templates...

	// leaving T2 undefined => error
	tc_3, err := t_base.Clone()
	loghttp.E(w, r, err, false)
	err = tc_3.ExecuteTemplate(w, "str_T0_outmost", dyndata)
	// NOT logging the error:
	// loghttp.E(w, r, err, false)

}
