/*

	Remember, there is a *quota* on getIntance...()  requests.

	Make instance info available in Memory.


	Have it refresh upon init and
	 every xxx minutes.

	There is a strange "version suffix".
	We call it VersionMinor and chop it off.



*/
package instance_mgt

import (
	"net/http"

	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/pbberlin/tools/net/http/loghttp"

	"appengine"
)

type Instance struct {
	InstanceID    string
	VersionFull   string // v2.243253253
	VersionMajor  string // v2
	VersionMinor  string //    243253253
	ModuleName    string
	NumInstances  int
	Hostname      string
	HostnameInst0 string
	HostnameInst1 string
	HostnameMod02 string
	LastUpdated   time.Time
}

//var info = new(Instance)
var m1 map[string]*Instance = map[string]*Instance{"dummy": new(Instance)}
var ci chan *Instance = make(chan *Instance)

//  &proto.RequiredNotSetError{field:"{Unknown}"}
const autoScalingErr1 = `proto: required field "{Unknown}" not set`
const autoScalingErr2 = `API error 3 (modules: INVALID_INSTANCES): The specified instance does not exist for this module/version.`
const autoScalingErrMsg = `NotSupportedWithautoScalingError`

func (i *Instance) String() string {

	if i == nil {
		return "<empty pointer to instance info>"
	}

	b1 := new(bytes.Buffer)
	b1.WriteString(fmt.Sprintf("Module %q Version %q %q \n", i.ModuleName,
		i.VersionMajor, i.VersionMinor))
	b1.WriteString(fmt.Sprintf("Number of Instances %v \n", i.NumInstances))
	b1.WriteString(fmt.Sprintf("Instance Id %q  \n", i.InstanceID))
	b1.WriteString("\n")

	b1.WriteString(fmt.Sprintf("Hostname Module 'default' Inst 'x'  %q \n",
		i.Hostname))
	b1.WriteString(fmt.Sprintf("Hostname Module '%v' Inst '0'  %q \n",
		i.ModuleName, i.HostnameInst0))
	b1.WriteString(fmt.Sprintf("Hostname Module '%v' Inst '1'  %q \n",
		i.ModuleName, i.HostnameInst1))
	b1.WriteString(fmt.Sprintf("Hostname Module 'mod02'   Inst 'x'  %q \n", i.HostnameMod02))

	b1.WriteString("\n")

	return b1.String()
}

func collectInfo(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	Get(w, r, m)
}

func onStart(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("instance started by appengine")

	// func() {
	// 	time.Sleep(200 * time.Millisecond)
	// 	collectInfo(w, r, map[string]interface{}{})
	// }()

}

func onStop(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("instance stopped by appengine")
}

func init() {

	// InstanceId := appengine.InstanceID() // does not during init, only after a few seconds

	http.HandleFunc("/instance-info/view", loghttp.Adapter(view))
	http.HandleFunc("/instance-info/collect", loghttp.Adapter(collectInfo))

	http.HandleFunc("/_ah/start", onStart)
	http.HandleFunc("/_ah/stop", onStop)

	go func() {
		for i := 0; i < 22222; i++ {
			// s := fmt.Sprintf("%v - %#v\n", util.TimeMarker(), info)
			_, ok := m1["info"]
			if ok {
				s := fmt.Sprintf("%v", m1["info"].Hostname)
				log.Print(s)
			}
			time.Sleep(12000 * time.Millisecond)
		}

	}()

}
