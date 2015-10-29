// Package util_appengine reveals if requests come from appengine or plain http servers;
// and if the gae development server is running.
package util_appengine

/*
	separated from common utils, because non-app engine projects
	can not use the util package otherwise
*/

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

import "os"
import "strings"

// http://regex101.com/

// IsLocalEnviron tells us, if we are on the
// local development server, or on the google app engine cloud maschine
func IsLocalEnviron() bool {

	return appengine.IsDevAppServer()

	s := os.TempDir()
	s = strings.ToLower(s)
	if s[0:2] == "c:" || s[0:2] == "d:" {
		// we are on windoofs - we are NOT on GAE
		return true
	}
	return false

}

//
func SafelyExtractGaeCtxError(r *http.Request) (context.Context, error) {
	if r == nil {
		return nil, fmt.Errorf("Request is not appengine - request is nil")
	}
	c := checkPanicking(r)
	if c != nil {
		return c, nil
	} else {
		return nil, fmt.Errorf("Request is not appengine")
	}
}

// Same as SafelyExtractGaeCtxError(), but without an error
func SafelyExtractGaeContext(r *http.Request) context.Context {
	if r == nil {
		return nil
	}
	c := checkPanicking(r)
	return c
}

func checkPanicking(r *http.Request) context.Context {
	defer func() {
		recover()
	}()
	c := appengine.NewContext(r)
	return c
}
