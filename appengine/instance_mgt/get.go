// Package instance_mgt computes instance info; views are in instance_info.
package instance_mgt

import (
	"net/http"

	"appengine"
	"appengine/module"

	"errors"
	"strings"
	"time"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/stringspb"
)

var ii *Instance = new(Instance) // initialization check via ii.LastUpdated.IsZero()

func GetStatic() *Instance {
	return ii
}

func Get(r *http.Request) *Instance {
	c := util_appengine.SafelyExtractGaeContext(r)
	return GetByContext(c)
}

// Todo: When c==nil we are in a non-appengine environment.
// We still want to return at least ii.PureHostname
func GetByContext(c appengine.Context) *Instance {

	tstart := time.Now()

	if !ii.LastUpdated.IsZero() {

		age := tstart.Sub(ii.LastUpdated)

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

			if err != nil {
				eStr := err.Error()
				eCmp1, eCmp2, eCmp3 := "API error", "INVALID_VERSION)", "Could not find the given version"
				if strings.Contains(eStr, eCmp1) && strings.Contains(eStr, eCmp2) && strings.Contains(eStr, eCmp3) {
					c.Infof("get num instances works only live and without autoscale; %v", err)
				} else {
					c.Errorf("get num instances error; %v", err)
				}
			}

		}

	}

	// in auto scaling, google reports "zero" - which can not be true
	// we assume at least 1
	if ii.NumInstances == 0 {
		ii.NumInstances = 1
	}

	// http://[0-2].1.default.libertarian-islands.appspot.com/instance-info

	ii.Hostname, err = appengine.ModuleHostname(c, ii.ModuleName, ii.VersionMajor, "")
	if err != nil {
		c.Errorf("ModuleHostname1: %v", err)
	}

	ii.PureHostname = appengine.DefaultVersionHostname(c)

	if !appengine.IsDevAppServer() {
		ii.HostnameInst0, err = appengine.ModuleHostname(c, ii.ModuleName, ii.VersionMajor, "0")
		if err != nil && (err.Error() == autoScalingErr1 || err.Error() == autoScalingErr2) {
			c.Infof("inst 0: " + autoScalingErrMsg)
			err = nil
		}
		if err != nil {
			c.Errorf("ModuleHostname2: %v", err)
		}

		ii.HostnameInst1, err = appengine.ModuleHostname(c, ii.ModuleName, ii.VersionMajor, "1")
		if err != nil && (err.Error() == autoScalingErr1 || err.Error() == autoScalingErr2) {
			c.Infof("inst 1: " + autoScalingErrMsg)
			err = nil
		}
		if err != nil {
			c.Errorf("ModuleHostname3: %v", err)
		}

		ii.HostnameMod02, err = appengine.ModuleHostname(c, "mod02", "", "")
		if err != nil {
			c.Infof("ModuleHostname4: %v", err)
		}

	}

	ii.LastUpdated = time.Now()

	c.Infof("collectInfo() completed, %v  - %v - %v - %v - %v, took %v",
		stringspb.Ellipsoider(ii.InstanceID, 4), ii.VersionMajor, ii.ModuleName,
		ii.Hostname, ii.PureHostname, time.Now().Sub(tstart))

	return ii
}
