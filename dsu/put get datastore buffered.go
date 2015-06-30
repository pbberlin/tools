package dsu

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"

	"fmt"

	"strings"

	"github.com/pbberlin/tools/appengine/instance_mgt"

	"io/ioutil"
	"net/http"

	"github.com/alexcesaro/mail/quotedprintable"
)

// BufPut - buffered put - saves its contents to memory, to memcache and the datastore
//   Todo: Don't buffer in local memory.
func BufPut(c appengine.Context, wb WrapBlob, skey string) (mkk string, errClosure error) {

	t := fmt.Sprintf("%T", wb)
	mkk = t + "__" + skey // kombi key

	errClosure = datastore.RunInTransaction(c,
		func(c appengine.Context) error {
			dskey1 := datastore.NewKey(c, t, skey, 0, nil)
			_, err := datastore.Put(c, dskey1, &wb)
			McacheSet(c, mkk, wb)
			memoryInstanceStore[mkk] = &wb
			multiCastInstanceCacheChange(c, mkk)
			c.Infof("saved to ds and memcache and instance RAM - combikey is %v", mkk)
			return err
		}, nil)
	c.Errorf("%v", errClosure)

	return
}

// BufGet - buffered get - fetches value from memory, or from memcache or from the datastore
//   todo: Paramter to reach "through" to datastore - without the buffer layers
func BufGet(c appengine.Context, mkk string) (WrapBlob, error) {

	wb1 := new(WrapBlob)

	// first check instance memory
	wb1, ok := memoryInstanceStore[mkk]
	if ok {
		c.Infof("received %q from static instance memory", mkk)
		//util_err.StackTrace(6)
		return *wb1, nil
	}

	// secondly check memcache
	ok = McacheGet(c, mkk, wb1)
	if ok && wb1 != nil && wb1.Name != "" {
		// we could replenish memcache TTL here - instead we do that below
		c.Infof("retrieved from memcache - combi_key %v", mkk)
		memoryInstanceStore[mkk] = wb1
		return *wb1, nil
	}

	// third: retrieve from datastore
	var wb2 = WrapBlob{}
	vk := strings.Split(mkk, "__")
	if len(vk) != 2 {
		return WrapBlob{}, fmt.Errorf("key must have one '__' delimiter; %q, size %v", mkk, len(vk))
	}
	t := vk[0]
	skey := vk[1]

	key := datastore.NewKey(c, t, skey, 0, nil)
	err := datastore.Get(c, key, &wb2)
	c.Errorf("%v", err)
	// missing entity and a present entity will both work.
	if err != nil && err != datastore.ErrNoSuchEntity {
		return wb2, err
	}
	McacheSet(c, mkk, wb2)
	memoryInstanceStore[mkk] = &wb2

	c.Infof("retrieved from ds - re-inserted into memcache + instance RAM - combi_key %v", mkk)

	return wb2, nil
}

func multiCastInstanceCacheChange(c appengine.Context, mkk string) {

	ii := instance_mgt.Get(c, map[string]interface{}{})

	/*
		making a get request to all instances
		submitting the key and the sender instance id
	*/
	for i := 0; i < ii.NumInstances; i++ {

		// http://[inst0-2].[v2].default.libertarian-islands.appspot.com/instance-info

		// note that Hostname already has VersionMajor and Module name as prefixed subdomains
		url := fmt.Sprintf("https://%v.%v/_ah/invalidate-instance-cache?mkk=%v&senderInstanceId=%v",
			i,
			ii.Hostname,
			mkk, ii.InstanceID)
		_ = url
		c.Infof(" url\n%v", url)

		//response, err := http.Get(url)  // not available in gae
		// instead:
		client := urlfetch.Client(c)
		response, err := client.Get(url)

		// with task queues - things would look similar:
		// t     := taskqueue.NewPOSTTask("/_ah/namespaced-counters/queue-pop", m)
		// but we could not enforce each instance getting
		// one and exactly one message

		// xmpp chat messages have the same disadvantage
		// the handler is the same for all instances
		//    /_ah/xmpp/message/chat/

		if err != nil {
			c.Infof("  could not launch get request; %v", err)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				c.Infof("  could not read response; %v", err)
			}
			c.Infof("%s\n", string(contents))
		}

	}

}

func invalidate(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	ii := instance_mgt.Get(c, map[string]interface{}{})

	mkk := r.FormValue("mkk")
	sii := r.FormValue("senderInstanceId")

	c.Infof(" %s  ---------  %s\n", sii, mkk)

	w.WriteHeader(http.StatusOK)

	if ii.InstanceID == sii {
		w.Write([]byte("Its ME " + mkk + "\n"))
		w.Write([]byte(sii))
	} else {
		w.Write([]byte("got it " + mkk + "\n"))
		w.Write([]byte(sii + "\n"))
		w.Write([]byte(ii.InstanceID))
	}

}

func showDsuObject(w http.ResponseWriter, r *http.Request) {

	var mkk string

	if r.FormValue("mkk") != "" {
		mkk = r.FormValue("mkk")

		c := appengine.NewContext(r)
		b, err := BufGet(c, mkk)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		boundary := b.S

		//dsu.WrapBlob__email-2015-05-02-16-47-52
		{

			// buf, msg := conv.VVByte_to_string(b.VVByte)
			// str := buf.String()

			sBytes := b.VByte
			str := string(sBytes)
			str = splitX(str, boundary)

			w.Write([]byte(fmt.Sprintf("raw1: %v<br>\n\n ", str)))

			reader := quotedprintable.NewDecoder(strings.NewReader(str))
			bufDec, err := ioutil.ReadAll(reader)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("err1 %v <br>\n", err)))
			} else {
				str2 := string(bufDec)
				w.Write([]byte(fmt.Sprintf("<hr>quoted printable<br>:\n %v\n\n ", str2)))
			}

		}

	} else {

		const frmMkk = `
		<div style='margin:8px;'>
			<form method="post" >
				MKK:  
				<input name="mkk"    value="dsu.WrapBlob__ chart_data_test_1"  size="80"  ><br/>
				<input type="submit" value="Fetch"            accesskey='f'></div>
			</form>
		</div>
		`

		w.Write([]byte(frmMkk))
	}

}

func splitX(s, sep string) string {

	ss0 := strings.Split(s, sep)

	if len(ss0) > 1 {
		ss1 := strings.Split(ss0[1], "<html>")
		if len(ss1) > 1 {
			return "<html>" + ss1[1]
		} else {
			return ss0[1]
		}
	}

	return s

}

func init() {
	http.HandleFunc("/dsu/show", showDsuObject)
	http.HandleFunc("/_ah/invalidate-instance-cache", invalidate)
}
