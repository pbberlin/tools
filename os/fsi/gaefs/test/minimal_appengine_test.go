package test

import (
	"log"
	"testing"

	"appengine/aetest"
)

func UncommentThis_TestMyFunction(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		log.Printf("%v\n", err)
		t.Fatal(err)
	}
	defer c.Close()

	// Run code and tests requiring the appengine.Context using c.
}
