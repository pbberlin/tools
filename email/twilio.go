package email

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/net/http/loghttp"

	"appengine"
	"appengine/urlfetch"
	// "github.com/subosito/twilio"
)

const (
	accountSid      = "ACcd078ef6d32ed27fbd3ace940b145594"
	authToken       = "6fd78d20a95a43b5362bb634da9d9224"
	phoneTo         = "+491729258218"
	phoneFrom       = "+491729258218"
	twilioURLPrefix = "https://api.twilio.com/2010-04-01/Accounts/"
	twilioFullURL   = twilioURLPrefix + accountSid + "/Calls.json"
	callbackUrl     = "http://libertarian-islands.appspot.com/twiml"
)

type TwiML struct {
	XMLName xml.Name `xml:"Response"`
	Say     string   `xml:",omitempty"`
	Play    string   `xml:",omitempty"`
}

var twiml1 TwiML = TwiML{Say: "Hello World!"}
var twiml2 TwiML = TwiML{Play: "http://demo.rickyrobinett.com/huh.mp3"}

func twilioResponse(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	x, err := xml.MarshalIndent(twiml1, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func call(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	// params trial: https://www.twilio.com/user/account/developer-tools/api-explorer/call-create
	// params docu:  https://www.twilio.com/docs/api/rest/making-calls
	v := url.Values{}
	v.Set("To", phoneTo)
	v.Set("From", phoneFrom)
	v.Set("Url", callbackUrl)
	v.Set("Method", "GET")
	v.Set("FallbackMethod", "GET")
	v.Set("StatusCallbackMethod", "GET")
	// v.Set("SendDigits", "32168")
	v.Set("Timeout", "4")
	v.Set("Record", "false")
	rb := *strings.NewReader(v.Encode())

	// Create Client
	// client := &http.Client{}  // local, not on appengine

	c := appengine.NewContext(r)
	// following appengine method, derived from big-query is non working:
	// config := oauth2_google.NewAppEngineConfig(c, []string{twilioURLPrefix})
	// client := &http.Client{Transport: config.NewTransport()}
	clientClassic := urlfetch.Client(c)
	//clientTwilio := twilio.NewClient(accountSid, authToken, clientClassic)  // not needed

	req, _ := http.NewRequest("POST", twilioFullURL, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := clientClassic.Do(req)
	loghttp.E(w, r, err, false, "something wrong with the http client")

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			fmt.Println(data["sid"])
		}
		fmt.Fprintf(w, "%#v", resp)
	} else {
		loghttp.E(w, r, err, false, "twilio response not ok", resp.Status)
	}
}

func init() {
	http.HandleFunc("/send-sms", loghttp.Adapter(sendSMS))
	http.HandleFunc("/call", loghttp.Adapter(call))
	http.HandleFunc("/twiml", loghttp.Adapter(twilioResponse))
}

// unused
// only here, we require the twilio package
//   github.com/subosito/twilio"
func sendSMS(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	/*
		// Initialize twilio client
		appEngineWay := true

		var clientTwilio *twilio.Client
		if appEngineWay {
			c := appengine.NewContext(c) // r is a *http.Request
			clientClassic := urlfetch.Client(a)
			clientTwilio = twilio.NewClient(accountSid, authToken, clientClassic)
		} else {
			clientTwilio = twilio.NewClient(accountSid, authToken, nil) // non appengine
		}

		var twMsg *twilio.Message
		var twRsp *twilio.Response
		var err error
		params := twilio.MessageParams{Body: "Hello Go!", MediaUrl: []string{}, StatusCallback: "", ApplicationSid: ""}
		// Send SMS
		twMsg, twRsp, err = clientTwilio.Messages.Send(phoneFrom, phoneTo, params)
		loghttp.E(w, r, err, false, "SMS Message send failed")

		fmt.Fprintf(w, "%+v  %+v\n", twMsg, twRsp)
	*/

}
