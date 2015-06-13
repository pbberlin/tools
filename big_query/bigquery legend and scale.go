package big_query

// https://godoc.org/code.google.com/p/google-api-go-client/bigquery/v2
// https://developers.google.com/bigquery/bigquery-api-quickstart
import (
	"bytes"
	"fmt"

	"net/http"
	//"appengine"
	"github.com/pbberlin/tools/colors"
	htmlpb "github.com/pbberlin/tools/pbhtml"
	"github.com/pbberlin/tools/util_appengine"
)

var p func(a ...interface{}) string = fmt.Sprint
var f func(format string, a ...interface{}) string = fmt.Sprintf

func legendAsHTML(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	b1 := new(bytes.Buffer)
	defer func() {
		w.Header().Set("Content-Type", "text/html")
		w.Write(b1.Bytes())
	}()

	b1, _ = disLegend(w, r)

}

// display Legend
func disLegend(w http.ResponseWriter, r *http.Request) (b1 *bytes.Buffer, m map[string]string) {

	//c := appengine.NewContext(r)

	b1 = new(bytes.Buffer)
	m = make(map[string]string)

	cd1 := GetChartDataFromDatastore(w, r, "chart_data_01")
	cd := *cd1

	span := htmlpb.GetSpanner()

	widthLabel := 80
	widthColorBox := 120
	widthDiv := widthLabel + widthColorBox + 2*4

	b1.WriteString(f("<div style='width:%dpx;margin:4px; padding: 4px; line-height:140%%; background-color:#eee;'>", widthDiv))

	for langIndex, lang := range cd.VLangs {

		gci := langIndex % len(colors.GraphColors) // graph color index
		// %x is the hex format, %2.2x makes padding zeros
		col := f("%2.2x%2.2x%2.2x", colors.GraphColors[gci][0],
			colors.GraphColors[gci][1],
			colors.GraphColors[gci][2])

		b1.WriteString(span(lang, widthLabel))

		block := f("<div style='display:inline-block;width:%dpx;height:5px; background-color:%s;'></div>", widthColorBox, col)
		b1.WriteString(span(block, widthColorBox))

		b1.WriteString("<br>")
		m[lang] = col

	}
	b1.WriteString("</div>")

	return
}

func init() {
	http.HandleFunc("/big-query/legend", util_appengine.Adapter(legendAsHTML))
}
