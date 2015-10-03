// Implements  coinbase.com integration
package coinbase

/*
https://developers.coinbase.com/docs/merchants/payment-buttons






app tec-news:
	https://www.coinbase.com/oauth/applications/560fbcaca4221973720002c7

	Developer access token
	https://api.coinbase.com/v1/users/self/
	?access_token=84a201e56a55185a785e707f952f8f5456f0ac6ec8f93465f12443ac0377a2ca


	https://www.coinbase.com/oauth/authorize
	?client_id=ca67d4b8b0fcfda6de7801cc225e03702d98a0e5a6e76d3d7478a9f08367fd3d
	&redirect_uri=https%3A%2F%2Ftec-news.appspot.com%2Fcoinbase-integr%2Fconfirm
	&response_type=code&scope=wallet%3Auser%3Aread



*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"appengine"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/stringspb"
)

const uriRequestPayment = "/coinbase-integr/request"
const uriConfirmPayment = "/coinbase-integr/confirm"

const coinbaseHost = "www.coinbase.com"

const walletAddress = "1E37asSURuvPDjjvPGSwAgDMnNgZDJdMDY" // for the entire account

// look into exclude.go
const (
	XX_apiKey    = "----------------------"
	XX_apiSecret = "------------------------" // salt for SHA256 signing
)

var wpf = fmt.Fprintf

// InitHandlers is called from outside,
// and makes the EndPoints available.
func InitHandlers() {
	http.HandleFunc(uriRequestPayment, loghttp.Adapter(requestPay))
	http.HandleFunc(uriConfirmPayment, loghttp.Adapter(confirmPay))
}

// BackendUIRendered returns a userinterface rendered to HTML
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Coinbase integration", uriRequestPayment, "request payment")
	htmlfrag.Wb(b1, "Confirm", uriConfirmPayment, "")

	return b1
}

func requestPay(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,
	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	protoc := "https://"
	if appengine.IsDevAppServer() {
		protoc = "http://"
	}

	host := appengine.DefaultVersionHostname(appengine.NewContext(r))
	if appengine.IsDevAppServer() {
		host = "not-localhost"
	}

	confirmURL := fmt.Sprintf("%v%v%v", protoc, host, uriConfirmPayment)
	confirmURL = url.QueryEscape(confirmURL)

	addrURL := fmt.Sprintf("https://%v/api/receive?method=create&address=%v&callback=%v&customsecret=49&api_code=%v",
		coinbaseHost, walletAddress, confirmURL, apiKey)

	req, err := http.NewRequest("GET", addrURL, nil)
	lg(err)
	if err != nil {
		return
	}
	bts, inf, err := fetch.UrlGetter(r, fetch.Options{Req: req})
	bts = bytes.Replace(bts, []byte(`","`), []byte(`", "`), -1)

	if err != nil {
		lg(err)
		lg(inf.Msg)
		return
	}

	lg("response body 1:\n")
	lg("%s\n", string(bts))

	lg("response body 2:\n")
	var data1 map[string]interface{}
	err = json.Unmarshal(bts, &data1)
	lg(err)
	lg(stringspb.IndentedDumpBytes(data1))
	// lg("%#v", data1)

	inputAddress, ok := data1["input_address"].(string)
	if !ok {
		lg("input address could not be casted to string; is type %T", data1["input_address"])
		return
	}
	feePercent, ok := data1["fee_percent"].(float64)
	if !ok {
		lg("fee percent could not be casted to float64; is type %T", data1["fee_percent"])
		return
	}

	lg("Input Adress will be %q; fee percent will be %4.2v", inputAddress, feePercent)

}

/*



*/
func confirmPay(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	/*


	   http://abc.de/ef?input_transaction_hash=46178baf7de078954b5aebb71c12120b33d998faac1c165af195eae90f19b25c&shared=false&address=18tpXf8WWuhJP95JbDASbZvavmZJbrydut&destination_address=18tpXf8WWuhJP95JbDASbZvavmZJbrydut&input_address=1ZTnjSdknZvur9Gc73gvB8XBTWL7nV1m6&test=true&anonymous=false&confirmations=0&value=82493362&transaction_hash=46178baf7de078954b5aebb71c12120b33d998faac1c165af195eae90f19b25c
	*/

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,
	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	htmlfrag.SetNocacheHeaders(w)

	wpf(b, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Payment confirmation"}))
	defer wpf(b, tplx.Foot)

	wpf(b, "<pre>")
	defer wpf(b, "</pre>")

	err := r.ParseForm()
	lg(err)

	custSecret := ""
	if r.FormValue("customsecret") != "" {
		custSecret = r.FormValue("customsecret")
	}
	lg("custom secret is %q", custSecret)

	val := ""
	if r.FormValue("value") != "" {
		val = r.FormValue("value")
	}
	lg("value is %q", val)

	w.WriteHeader(http.StatusOK)

}
