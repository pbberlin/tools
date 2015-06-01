package namespaces_taskqueues

import (
	"io"
	"net/http"
	//"net/url"

	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
	"github.com/pbberlin/tools/util_err"
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

func readBothNamespaces(w http.ResponseWriter, r *http.Request) {

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
	util_err.Err_log(err)

	{
		c, err = appengine.Namespace(c, altNamespace)
		util_err.Err_log(err)
		c2, err = agnosticReadReset(c, reset)
		util_err.Err_log(err)
	}

	s1 = fmt.Sprintf("%v", c1)
	s2 = fmt.Sprintf("%v", c2)

	io.WriteString(w, "|"+s1+"|    |"+s2+"|")
	if reset {
		io.WriteString(w, "     and reset")
	}

}

func incrementBothNamespaces(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	err := agnosticIncrement(c)
	util_err.Err_log(err)

	{
		c, err := appengine.Namespace(c, altNamespace)
		util_err.Err_log(err)
		err = agnosticIncrement(c)
		util_err.Err_log(err)
	}

	s := `counters updates für ns=''  and ns='ns01'.` + "\n"
	io.WriteString(w, s)
	readBothNamespaces(w, r)

}

func queuePush(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	m := map[string][]string{"counter_name": []string{nscStringKey}}
	t := taskqueue.NewPOSTTask("/_ah/namespaced-counters/queue-pop", m)

	taskqueue.Add(c, t, "")

	c, err := appengine.Namespace(c, altNamespace)
	util_err.Err_log(err)
	taskqueue.Add(c, t, "")

	io.WriteString(w, "tasks enqueued\n")

	io.WriteString(w, "\ncounter values now: \n")
	readBothNamespaces(w, r)

	io.WriteString(w, "\n\n...sleeping... \n")
	time.Sleep(time.Duration(400) * time.Millisecond)
	readBothNamespaces(w, r)

}

func queuePop(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := agnosticIncrement(c)
	c.Infof("qp")
	util_err.Err_log(err)
}

func init() {
	http.HandleFunc("/namespaced-counters/increment", incrementBothNamespaces)
	http.HandleFunc("/namespaced-counters/read", readBothNamespaces)

	http.HandleFunc("/_ah/namespaced-counters/queue-pop", queuePop)
	http.HandleFunc("/namespaced-counters/queue-push", queuePush)

}
