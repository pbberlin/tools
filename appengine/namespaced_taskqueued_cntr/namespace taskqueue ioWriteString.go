package namespaced_taskqueued_cntr

import (
	"io"
	"net/http"
	//"net/url"

	"fmt"
	"time"

	"github.com/pbberlin/tools/net/http/loghttp"

	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
)

type Counter struct {
	Count int64
}

const nscStringKey = "singularKey"
const nscKind = "NamespaceCounter" // namespace counter kind

const altNamespace = "ns01"

func agnosticIncrement(c appengine.Context) error {

	key := datastore.NewKey(c, nscKind, nscStringKey, 0, nil)

	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		var ctr Counter
		err := datastore.Get(c, key, &ctr)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		ctr.Count++
		_, err = datastore.Put(c, key, &ctr)
		c.Infof("+1")
		return err
	}, nil)
}

func agnosticReadReset(c appengine.Context, doReset bool) (int64, error) {

	key := datastore.NewKey(c, nscKind, nscStringKey, 0, nil)

	var ctrRd Counter
	err := datastore.Get(c, key, &ctrRd)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}

	if doReset {
		var ctrSt Counter
		ctrSt.Count = -1
		_, err = datastore.Put(c, key, &ctrSt)
	}

	return ctrRd.Count, err
}

func readBothNamespaces(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = b

	var c1, c2 int64
	var s1, s2 string
	var err error
	var reset bool = false

	p := r.FormValue("reset")
	if p != "" {
		reset = true
	}

	c := appengine.NewContext(r)
	c1, err = agnosticReadReset(c, reset)
	lg(err)

	{
		c, err = appengine.Namespace(c, altNamespace)
		lg(err)
		c2, err = agnosticReadReset(c, reset)
		lg(err)
	}

	s1 = fmt.Sprintf("%v", c1)
	s2 = fmt.Sprintf("%v", c2)

	io.WriteString(w, "|"+s1+"|    |"+s2+"|")
	if reset {
		io.WriteString(w, "     and reset")
	}

}

func incrementBothNamespaces(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = b

	c := appengine.NewContext(r)
	err := agnosticIncrement(c)
	lg(err)

	{
		c, err := appengine.Namespace(c, altNamespace)
		lg(err)
		err = agnosticIncrement(c)
		lg(err)
	}

	s := `counters updates für ns=''  and ns='ns01'.` + "\n"
	io.WriteString(w, s)
	readBothNamespaces(w, r, m)

}

func queuePush(w http.ResponseWriter, r *http.Request, mx map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = b

	c := appengine.NewContext(r)

	m := map[string][]string{"counter_name": []string{nscStringKey}}
	t := taskqueue.NewPOSTTask("/_ah/namespaced-counters/queue-pop", m)

	taskqueue.Add(c, t, "")

	c, err := appengine.Namespace(c, altNamespace)
	lg(err)
	taskqueue.Add(c, t, "")

	io.WriteString(w, "tasks enqueued\n")

	io.WriteString(w, "\ncounter values now: \n")
	readBothNamespaces(w, r, mx)

	io.WriteString(w, "\n\n...sleeping... \n")
	time.Sleep(time.Duration(400) * time.Millisecond)
	readBothNamespaces(w, r, mx)

}

func queuePop(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = b

	c := appengine.NewContext(r)
	err := agnosticIncrement(c)
	c.Infof("qp")
	lg(err)
}

func init() {
	http.HandleFunc("/namespaced-counters/increment", loghttp.Adapter(incrementBothNamespaces))
	http.HandleFunc("/namespaced-counters/read", loghttp.Adapter(readBothNamespaces))

	http.HandleFunc("/_ah/namespaced-counters/queue-pop", loghttp.Adapter(queuePop))
	http.HandleFunc("/namespaced-counters/queue-push", loghttp.Adapter(queuePush))

}
