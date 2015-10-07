// Implements  coinbase.com integration
package coinbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"appengine"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"
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

/*

requestPay is unused.

We use payment buttons instead.
https://developers.coinbase.com/docs/merchants/payment-buttons


Oauth is also not used.
Look here for a preconfigured app with oauth:
	https://www.coinbase.com/oauth/applications/560fbcaca4221973720002c7

*/
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

https://developers.coinbase.com/docs/merchants/callbacks


id				Order number used to uniquely identify an order on Coinbase
completed_at	ISO 8601 timestamp when the order completed
status			[completed, mispaid, expired]
event			[completed, mispayment]. If mispayment => check key mispayment_id. Distinction from status ...
total_btc		Total amount of the order in ‘satoshi’ (1 BTC = 100,000,000 Satoshi). Note the use of the word ‘cents’ in the callback really means satoshi in this context. The btc amount will be calculated at the current exchange rate at the time the order is placed (current to within 15 minutes).
total_native	Units of local currency. 1 unit = 100 cents. Equal to the price from creating the button.
total_payout	Units of local currency deposited using instant payout.
custom			Custom parameter from data-custom attribute of button. Usually an Order, User, or Product ID
receive_address	Bitcoin address associated with this order. This is where the payment was sent.
button			Button details. ID matches the data-code parameter in your embedded HTML code.
transaction		Hash and number of confirmations of underlying transaction.
				Number of confirmations typically zero at the time of the first callback.
customer		Customer information from order form. Can include email xor shipping address.
refund_address	Experimental parameter that is subject to change.


*/
func confirmPay(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		// loghttp.Pf(w, r, b.String())
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

	if len(bts) < 1 {
		lg("lo empty post body")
		w.WriteHeader(http.StatusOK)
		b = new(bytes.Buffer)
		return
	}

	var mp map[string]interface{}
	err = json.Unmarshal(bts, &mp)
	lg(err)

	mpPayout := submap(mp, "payout", lg)
	if len(mpPayout) > 0 {
		lg("lo " + stringspb.IndentedDump(mpPayout))
	}
	mpAddress := submap(mp, "address", lg)
	if len(mpAddress) > 0 {
		lg("lo " + stringspb.IndentedDump(mpAddress))
	}

	var cents, BTC float64
	var status string

	mpOrder := submap(mp, "order", lg)
	if len(mpOrder) < 1 {
		w.WriteHeader(http.StatusLengthRequired)
		lg("mpOrder not present %v", status)
		return
	} else {
		lg("lo " + stringspb.IndentedDump(mpOrder))

		mpBTC := submap(mpOrder, "total_btc", lg)
		// lg("lo " + stringspb.IndentedDump(mpBTC))

		if icents, ok := mpBTC["cents"]; ok {
			cents, ok = icents.(float64)
			if !ok {
				lg(" mpBTC[cents] is of unexpected type %T ", mpBTC["cents"])
			}
			BTC = cents / (1000 * 1000 * 100)

		} else {
			lg(" mpBTC[cents] not present")
		}
		lg("received %18.2f satoshi, %2.9v BTC ", cents, BTC)

		if _, ok := mpOrder["status"]; ok {
			status, ok = mpOrder["status"].(string)
			if !ok {
				lg(" mpOrder[status] is of unexpected type %T ", mpOrder["status"])
			}
		}

		lg("status    %v  ", status)
		lg("custom   %#v  ", mpOrder["custom"])
		lg("customer %#v - mostly empty", mpOrder["customer"])

		var values url.Values
		if _, ok := mpOrder["custom"]; ok {
			var err error
			values, err = url.ParseQuery(mpOrder["custom"].(string))
			lg(err)
			if err != nil {
				w.WriteHeader(http.StatusLengthRequired)
				lg("unsatisfactory query in custom string %v", mpOrder["custom"])
				return
			}
		} else {
			w.WriteHeader(http.StatusLengthRequired)
			lg("custom string not present")
			return
		}

		//  save
		if status == "completed" {
			blob := dsu.WrapBlob{
				VByte: stringspb.IndentedDumpBytes(mpOrder),
			}
			blob.Name = values.Get("uID")
			blob.S = values.Get("productID")
			blob.Desc = status
			blob.F = BTC

			blob.VVByte, _ = conv.String_to_VVByte(string(blob.VByte)) // just to make it readable

			newKey, err := dsu.BufPut(appengine.NewContext(r), blob, blob.Name)
			lg("key is %v", newKey)
			lg(err)

			retrieveAgain, err := dsu.BufGet(appengine.NewContext(r), "dsu.WrapBlob__"+blob.Name)
			lg(err)
			lg("retrieved %v %v %v", retrieveAgain.Name, retrieveAgain.Desc, retrieveAgain.F)

		} else {
			w.WriteHeader(http.StatusLengthRequired)
			lg("unsatisfactory status %v", status)
			return
		}

	}
	w.WriteHeader(http.StatusOK)
	b = new(bytes.Buffer)

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
		// lg("mp[%v] not present", key)
	}

	return mp
}
