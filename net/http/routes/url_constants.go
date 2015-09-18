// Package routes contains route path constants,
// that would otherwise cause circular dependencies.
package routes

import "appengine"

const URLParamKey = "url-x"

const ProxifyURI = "/prox"
const WeedOutURI = "/weedout"
const FetchSimilarURI = "/fetch/similar"

const FormRedirector = "/blob2/form-redirector"

var AppHost01 = "localhost:8085"

func init() {
	if !appengine.IsDevAppServer() {
		AppHost01 = "libertarian-islands.appspot.com"
	}
}
