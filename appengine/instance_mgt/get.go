package instance_mgt

import (
	"net/http"

	"appengine"
	"appengine/module"

	"errors"
	"strings"
	"time"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"
)

var ii = new(Instance)

func Get(w http.ResponseWriter, r *http.Request, m map[string]interface{}) *Instance {

	c := appengine.NewContext(r)

	startFunc := time.Now()

	if !ii.LastUpdated.IsZero() {

		age := startFunc.Sub(ii.LastUpdated)

		if age < 200*time.Millisecond {
			c.Infof("instance info update too recently: %v, skipping.\n", age)
			return ii
		}

		if age < 1*time.Hour {
			if len(ii.Hostname) > 2 {
				return ii
			}

		}

		c.Infof("instance info update too old: %v, recomputing.\n", age)
	}

	ii.ModuleName = appengine.ModuleName(c)
	ii.InstanceID = appengine.InstanceID()
	ii.VersionFull = appengine.VersionID(c)

	majorMinor := strings.Split(ii.VersionFull, ".")
	if len(majorMinor) != 2 {
		panic("we need a version string of format X.Y")
	}

	ii.VersionMajor = majorMinor[0]
	ii.VersionMinor = majorMinor[1]

	var err = errors.New("dummy creation error message")

	ii.NumInstances, err = module.NumInstances(c, ii.ModuleName, ii.VersionFull)
	if err != nil {
		// this never works with version full
		// we do not log this - but try version major
		err = nil

		if !util_appengine.IsLocalEnviron() {
			ii.NumInstances, err = module.NumInstances(c, ii.ModuleName, ii.VersionMajor)
			loghttp.E(w, r, err, true, "get num instances works only live and without autoscale")
		}

	}

	// in auto scaling, google reports "zero" - which can not be true
	// we assume at least 1
	if ii.NumInstances == 0 {
		ii.NumInstances = 1
	}

	// http://[0-2].1.default.libertarian-islands.appspot.com/instance-info

	ii.Hostname, err = appengine.ModuleHostname(c, ii.ModuleName,
		ii.VersionMajor, "")
	loghttp.E(w, r, err, false)

	if !util_appengine.IsLocalEnviron() {
		ii.HostnameInst0, err = appengine.ModuleHostname(c, ii.ModuleName,
			ii.VersionMajor, "0")
		if err != nil && (err.Error() == autoScalingErr1 || err.Error() == autoScalingErr2) {
			c.Infof("inst 0: " + autoScalingErrMsg)
			err = nil
		}
		loghttp.E(w, r, err, true)

		ii.HostnameInst1, err = appengine.ModuleHostname(c, ii.ModuleName,
			ii.VersionMajor, "1")
		if err != nil && (err.Error() == autoScalingErr1 || err.Error() == autoScalingErr2) {
			c.Infof("inst 1: " + autoScalingErrMsg)
			err = nil
		}
		loghttp.E(w, r, err, true)

		ii.HostnameMod02, err = appengine.ModuleHostname(c, "mod02", "", "")
		loghttp.E(w, r, err, true)

	}

	ii.LastUpdated = time.Now()

	c.Infof("collectInfo() completed, %v.%v.%v.%v, took %v",
		stringspb.Ellipsoider(ii.InstanceID, 4), ii.VersionMajor, ii.ModuleName,
		ii.Hostname,
		ii.LastUpdated.Sub(startFunc))

	return ii
}
