/*  Package instance_mgt tries to send and receive data to sibling instances.


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

	"google.golang.org/appengine"

	"bytes"
	"fmt"
	"log"
	"time"

	aelog "google.golang.org/appengine/log"
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

	PureHostname string // without version, module
}

//var info = new(Instance)
// var ci chan *Instance = make(chan *Instance)

//  &proto.RequiredNotSetError{field:"{Unknown}"}
const autoScalingErr1 = `proto: required field "{Unknown}" not set`
const autoScalingErr2 = `API error 3 (modules: INVALID_INSTANCES): The specified instance does not exist for this module/version.`
const autoScalingErrMsg = `NotSupportedWithautoScalingError`

func (i *Instance) String() string {

	if i == nil {
		return "<empty pointer to instance info>"
	}

	b1 := new(bytes.Buffer)
	b1.WriteString(fmt.Sprintf("Module %q Version %q %q \n", i.ModuleName, i.VersionMajor, i.VersionMinor))
	b1.WriteString(fmt.Sprintf("Number of Instances %v \n", i.NumInstances))
	b1.WriteString(fmt.Sprintf("Instance Id %q  \n", i.InstanceID))
	b1.WriteString("\n")

	b1.WriteString(fmt.Sprintf("Hostname Pure                       %q \n", i.PureHostname))
	b1.WriteString(fmt.Sprintf("Hostname Module 'default' Inst 'x'  %q \n", i.Hostname))
	b1.WriteString(fmt.Sprintf("Hostname Module '%v' Inst '0'  %q \n", i.ModuleName, i.HostnameInst0))
	b1.WriteString(fmt.Sprintf("Hostname Module '%v' Inst '1'  %q \n", i.ModuleName, i.HostnameInst1))
	b1.WriteString(fmt.Sprintf("Hostname Module 'mod02'   Inst 'x'  %q \n", i.HostnameMod02))

	b1.WriteString("\n")

	return b1.String()
}

func onStart(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aelog.Infof(c, "instance started by appengine")

	// func() {
	// 	time.Sleep(200 * time.Millisecond)
	// 	collectInfo(w, r, map[string]interface{}{})
	// }()

}

func onStop(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aelog.Infof(c, "instance stopped by appengine")
}

func init() {

	// InstanceId := appengine.InstanceID() // does not during init, only after a few seconds
	http.HandleFunc("/_ah/start", onStart)
	http.HandleFunc("/_ah/stop", onStop)

	go func() {
		for {
			if ii.LastUpdated.IsZero() {
				time.Sleep(36 * time.Second)
				continue
			}

			s := fmt.Sprintf("hostname is %v, pure hostname is %v (%v)", ii.Hostname, ii.PureHostname, time.Now().Sub(ii.LastUpdated))
			log.Println(s)
			time.Sleep(36 * time.Second)
		}

	}()

}
