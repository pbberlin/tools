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
	"io/ioutil"
	"net/http"
	"net/url"

	"appengine"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
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

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,
	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	htmlfrag.SetNocacheHeaders(w)

	//____________________________________________________________________

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		lg("cannot read resp body: %v", err)
		return
	}
	defer r.Body.Close()

	// lg("bytes are -%s-", stringspb.ToLen(string(bts), 20))

	var mp map[string]interface{}
	err = json.Unmarshal(bts, &mp)
	lg(err)

	mpOrder := submap(mp, "order", lg)
	lg("lo " + stringspb.IndentedDump(mpOrder))

	mpBTC := submap(mpOrder, "total_btc", lg)
	lg("lo " + stringspb.IndentedDump(mpBTC))

	// if branchTemp, ok := mp["order"]; ok {
	// 	var okConv bool
	// 	mp, okConv = branchTemp.(map[string]interface{})
	// 	if !okConv {
	// 		lg(" mp[order] of type %T ", branchTemp)
	// 	}

	// } else {
	// 	lg("mp[order] not present")
	// }

	w.WriteHeader(http.StatusOK)

}

func submap(mpArg map[string]interface{}, key string, lg loghttp.FuncBufUniv) map[string]interface{} {

	var mp map[string]interface{}

	if branchTemp, ok := mpArg[key]; ok {
		var okConv bool
		mp, okConv = branchTemp.(map[string]interface{})
		if !okConv {
			lg(" mp[%v] of type %T ", key, branchTemp)
		}

	} else {
		lg("mp[%v] not present", key)
	}

	return mp
}
