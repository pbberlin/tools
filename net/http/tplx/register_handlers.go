package tplx

import (
	"bytes"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
)

func InitHandlers() {
	http.HandleFunc("/tpl/demo1", loghttp.Adapter(templatesCompileDemo))
	http.HandleFunc("/tpl/demo2", loghttp.Adapter(templatesDemo2))
	http.HandleFunc("/tpl/reset", loghttp.Adapter(TemplateFromHugoReset))
}

// BackendUIRendered returns a userinterface rendered to HTML
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Tpl Demo 1", "/tpl/demo1", "")
	htmlfrag.Wb(b1, "Tpl Demo 2", "/tpl/demo1", "")
	htmlfrag.Wb(b1, "Tpl Reset", "/tpl/reset", "")
	return b1
}
