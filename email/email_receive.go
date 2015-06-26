package email

import (
	"net/http"

	go_mail "net/mail"

	ae_mail "appengine/mail"

	"appengine"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"

	"strings"

	"bytes"
)

/*

	Domain is different!!!

	appspotMAIL.com

	not

	appspot.com

	peter@libertarian-islands@appspotmail.com


  https://developers.google.com/appengine/docs/python/mail/receivingmail

 	email-address:	 string@appid.appspotmail.com
 	is routed to
   /_ah/mail/string@appid.appspotmail.com

   peter@libertarian-islands.appspotmail.com
 	is routed to
   /_ah/mail/peter@libertarian-islands.appspotmail.com
*/
func emailReceiveAndStore(w http.ResponseWriter, r *http.Request, mx map[string]interface{}) {

	c := appengine.NewContext(r)
	defer r.Body.Close()

	msg, err := go_mail.ReadMessage(r.Body)
	loghttp.E(w, r, err, false, "could not do ReadMessage")
	if msg == nil {
		c.Warningf("-empty msg- " + r.URL.Path)
		return
	}

	// see http://golang.org/pkg/net/mail/#Message
	b1 := new(bytes.Buffer)
	// for i, m1 := range msg.Header {
	// 	c.Infof("--msg header %q : %v", i, m1)
	// }

	from := msg.Header.Get("from") + "\n"
	b1.WriteString("from: " + from)

	to := msg.Header.Get("to") + "\n"
	b1.WriteString("to: " + to)

	subject := msg.Header.Get("subject") + "\n"
	b1.WriteString("subject: " + subject)

	when, _ := msg.Header.Date()
	swhen := when.Format("2006-01-02 - 15:04 \n")
	b1.WriteString("when: " + swhen)

	ctype := msg.Header.Get("Content-Type")
	c.Infof("content type header: %q", ctype)
	boundary := ""
	// [multipart/mixed; boundary="------------060002090509030608020402"]
	if strings.HasPrefix(ctype, "[multipart/mixed") ||
		strings.HasPrefix(ctype, "multipart/mixed") {
		vT1 := strings.Split(ctype, ";")
		if len(vT1) > 1 {
			c.Infof("substring 1: %q", vT1[1])
			sT11 := vT1[1]
			sT11 = strings.TrimSpace(sT11)
			sT11 = strings.TrimPrefix(sT11, "boundary=")
			sT11 = strings.Trim(sT11, `"`)
			boundary = sT11
			c.Infof("substring 2: %q", boundary)
		}
	}

	b1.WriteString("\n\n")
	b1.ReadFrom(msg.Body)

	dsu.McacheSet(c, keyLatest, dsu.WrapBlob{Name: subject, S: boundary, VByte: b1.Bytes()})
	if strings.HasPrefix(to, "foscam") {
		// send confirmation to sender
		var m map[string]string = nil
		m = make(map[string]string)
		m["sender"] = from
		m["subject"] = "confirmation: " + subject
		emailSend(w, r, m)

		parseFurther(w, r, true)
		call(w, r, mx)
	} else {
		blob := dsu.WrapBlob{Name: subject + "from " + from + "to " + to,
			S: boundary, VByte: b1.Bytes()}
		blob.VVByte, _ = conv.String_to_VVByte(b1.String())
		dsu.BufPut(w, r, blob, "email-"+util.TimeMarker())
	}

}

func emailReceiveSimple(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	defer r.Body.Close()

	/* alternative code from https://developers.google.com/appengine/docs/go/mail/
	not net/mail ,but
	*/
	var b2 bytes.Buffer
	if _, err := b2.ReadFrom(r.Body); err != nil {
		c.Errorf("Error reading body: %v", err)
		return
	}
	c.Infof("\n\nb2: " + b2.String() + "--\n ")

}

func emailSend(w http.ResponseWriter, r *http.Request, m map[string]string) {

	c := appengine.NewContext(r)
	//addr := r.FormValue("email")

	if _, ok := m["subject"]; !ok {
		m["subject"] = "empty subject line"
	}

	email_thread_id := []string{"3223"}

	msg := &ae_mail.Message{
		//Sender:  "Peter Buchmann <peter.buchmann@web.de",
		//		Sender: "peter.buchmann@web.de",
		Sender: "peter.buchmann.68@gmail.com",
		//To:	   []string{addr},
		To: []string{"peter.buchmann@web.de"},

		Subject: m["subject"],
		Body:    "some_body some_body2",
		Headers: go_mail.Header{"References": email_thread_id},
	}
	err := ae_mail.Send(c, msg)
	loghttp.E(w, r, err, false, "could not send the email")

}

func init() {
	http.HandleFunc("/_ah/mail/", loghttp.Adapter(emailReceiveAndStore))
	//http.HandleFunc("/_ah/mail/"  , loghttp.Adapter(emailReceiveSimple))
}
