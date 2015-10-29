// Implements  coinbase.com integration
package coinbase

/*

requestPay is unused.

We use payment buttons instead.
See https://developers.coinbase.com/docs/merchants/payment-buttons



Oauth is also not used.
Look here for a preconfigured app with oauth:
	https://www.coinbase.com/oauth/applications/560fbcaca4221973720002c7

*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"
	"google.golang.org/appengine"
)

const uriRequestPayment = "/coinbase-integr/request" // unused
const uriConfirmPayment = "/coinbase-integr/confirm"
const uriRedirectSuccess = "/coinbase-integr/redir-success1"

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
	http.HandleFunc(uriRedirectSuccess, loghttp.Adapter(paymentSuccess))
}

// BackendUIRendered returns a userinterface rendered to HTML
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Coinbase integration", uriRequestPayment, "request payment")
	htmlfrag.Wb(b1, "Confirm", uriConfirmPayment, "")
	return b1
}

const BtnTestFormat = `
					<a class="coinbase-button"
						data-code="0025d69ea925b48ba2b7adeb2a911ca2"
						data-custom="productID=%v&uID=%v"
						data-env="sandbox"
						href="https://sandbox.coinbase.com/checkouts/0025d69ea925b48ba2b7adeb2a911ca2"
					>Pay With Bitcoin</a>
					<script src="https://sandbox.coinbase.com/assets/button.js" type="text/javascript"></script>
					`

const BtnLiveFormat = `
					<a class="coinbase-button" 
						data-code="aa4e03abbc5e2f5321d27df32756a932" 
						data-custom="productID=%v&uID=%v" 
						href="https://www.coinbase.com/checkouts/aa4e03abbc5e2f5321d27df32756a932" 
					>Pay With Bitcoin</a>
					<script src="https://www.coinbase.com/assets/button.js" type="text/javascript"></script>

				`

//
//
// requestPay is unused
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

	// Response body contains the suggested bitcoin address for payment.
	// And the minimum recommended fee percentage
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
			if mpOrder["custom"] == "123456789" {
				lg("test request recognized")
				values = url.Values{}
				values.Add("uID", "testUser123")
				values.Add("productID", "/member/somearticle")
			} else {
				var err error
				values, err = url.ParseQuery(mpOrder["custom"].(string))
				lg(err)
				if err != nil {
					w.WriteHeader(http.StatusLengthRequired)
					lg("unsatisfactory query in custom string %v", mpOrder["custom"])
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusLengthRequired)
			lg("custom string not present")
			return
		}

		//  save
		if status == "completed" {
			lg("status 'completed'")
			blob := dsu.WrapBlob{
				VByte: stringspb.IndentedDumpBytes(mpOrder),
			}
			blob.Name = values.Get("uID")
			blob.Category = "invoice"
			blob.S = values.Get("productID")
			blob.Desc = status
			blob.F = BTC
			blob.I = int(time.Now().Unix())

			// blob.VVByte, _ = conv.String_to_VVByte(string(blob.VByte)) // just to make it readable

			newKey, err := dsu.BufPut(appengine.NewContext(r), blob, blob.Name+blob.S)
			lg("key is %v", newKey)
			lg(err)

			retrieveAgain, err := dsu.BufGet(appengine.NewContext(r), "dsu.WrapBlob__"+blob.Name+blob.S)
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

/*
https://tec-news.appspot.com/coinbase-integr/redir-success1?
order[button][description]=When and how Bitcoin decline will start.&
order[button][id]=0025d69ea925b48ba2b7adeb2a911ca2&
order[button][name]=Bitcoin Analysis&
order[button][repeat]=&
order[button][resource_path]=/v2/checkouts/4f1e5ecc-c8fc-56fc-926c-15a7eebd8314&
order[button][subscription]=false&
order[button][type]=buy_now&
order[button][uuid]=4f1e5ecc-c8fc-56fc-926c-15a7eebd8314&
order[created_at]=2015-10-26 08:03:17 -0700&
order[custom]=productID=/member/tec-news/crypto-experts-neglect-one-vital-aspect&uID=14952300052240127534&
order[event]=&
order[id]=GAB5VN36&
order[metadata]=&
order[receive_address]=myL84ofiymQpzzmJ7Foc9F2wQ4GMuSuQ3f&
order[refund_address]=mwaz3wxMbnZrBZUSZpVHr51xjQ6Swx756b&
order[resource_path]=/v2/orders/9bbf6fde-530a-53a4-bf94-d54fc3f43d40&
order[status]=completed&
order[total_btc][cents]=5600.0&
order[total_btc][currency_iso]=BTC&
order[total_native][cents]=50.0&
order[total_native][currency_iso]=EUR&
order[total_payout][cents]=0.0&
order[total_payout][currency_iso]=USD&
order[transaction][confirmations]=0&
order[transaction][hash]=ada26d75ff1e16b4febf539433d5260441171560c57adfff2ac968be37108112&
order[transaction][id]=562e40dede472f26be000018&
order[uuid]=9bbf6fde-530a-53a4-bf94-d54fc3f43d40
*/
func paymentSuccess(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	custom := r.Form.Get("order[custom]")
	// w.Write([]byte("custom=" + custom + "<br>\n"))

	values, err := url.ParseQuery(custom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productID := values.Get("productID")
	uID := values.Get("uID")

	if productID != "" {
		http.Redirect(w, r, productID+"?redirected-from=paymentsucc", http.StatusFound)
		return
	}

	w.Write([]byte("productID=" + productID + " uID=" + uID + "<br>\n"))

}
