package tests

import (
	"log"
	"testing"

	"appengine/aetest"
)

func UncommentThis_TestFunction(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		log.Printf("%v\n", err)
		t.Fatal(err)
	}
	defer c.Close()

	// Now you can run the code and tests requiring
	// the appengine.Context
}
