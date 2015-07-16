// +build http
// go test -tags=http

package fstest

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/pbberlin/tools/os/fsi/httpfs"
)

func TestHttp(t *testing.T) {

	testRoot := "c:\\temp"
	if runtime.GOOS != "windows" {
		testRoot = "/tmp"
	}

	bb := new(bytes.Buffer)
	msg := ""
	_ = bb

	os.Chdir(testRoot)
	pwd, _ := os.Getwd()
	if pwd == testRoot {
		// os.RemoveAll(pwd)
	}

	Fss, c := initFileSystems()
	defer c.Close()

	// Fss = Fss[0:1]
	// Fss = Fss[1:2]
	// Fss = Fss[2:3]

	portInc := 2
	for _, fs := range Fss {

		wpf(os.Stdout, "-----created fs %v %v-----\n", fs.Name(), fs.String())

		bb, msg = CreateSys(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		httpFSys := &httpfs.HttpFs{SourceFs: fs}

		mux := http.NewServeMux()

		//
		fileserver1 := http.FileServer(httpFSys.Dir("./"))
		mux.Handle("/", fileserver1)

		//
		fileserver2 := http.FileServer(http.Dir("./css/"))
		mux.Handle("/css/", http.StripPrefix("/css/", fileserver2))

		go func(arg int) {
			addr := fmt.Sprintf("localhost:400%v", arg)
			log.Printf("serving %v\n", addr)
			log.Fatal(http.ListenAndServe(addr, mux))
		}(portInc)

		portInc++

	}

	time.Sleep(time.Second * 111)

}
