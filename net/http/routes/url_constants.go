// Package routes contains hostname logic and path constants,
// that would otherwise cause circular dependencies.
package routes

import "appengine"

const URLParamKey = "url-x"

const ProxifyURI = "/prox"
const DedupURI = "/dedup"
const FetchSimilarURI = "/fetch/similar"

const FormRedirector = "/blob2/form-redirector"

var appID, devServerPort = "", ""
var appHost = ""

var devAdminPort = ""

func InitAppHost(ID, port, adminPort string) {
	appID = ID
	devServerPort = port
	devAdminPort = adminPort

	if appengine.IsDevAppServer() {
		appHost = "localhost" + ":" + devServerPort
	} else {
		appHost = appID + ".appspot.com"
	}

}

func AppHostDev() string {
	return "localhost" + ":" + devServerPort
}
func AppHostLive() string {
	return appID + ".appspot.com"
}

func AppHost() string {
	return appHost
}
func AppID() string {
	return appID
}

func DevAdminPort() string {
	return devAdminPort
}
