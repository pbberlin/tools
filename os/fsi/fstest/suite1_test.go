// +build suite1
// go test -tags=suite1

package fstest

// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright 2009 The Go Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"io"
	"testing"

	"github.com/pbberlin/tools/os/fsi/memfs"
)

//Read with length 0 should not return EOF.
func TestRead0(t *testing.T) {

	Fss, c := initFileSystems()
	defer c.Close()
	for _, fs := range Fss {
		path := testDir + "/" + testName
		if err := fs.MkdirAll(testDir, 0777); err != nil {
			t.Fatal(fs.Name(), "unable to create dir", err)
		}

		f, err := fs.Create(path)
		if err != nil {
			t.Fatal(fs.Name(), "create failed:", err)
		}
		defer f.Close()
		f.WriteString("Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

		b := make([]byte, 0)
		n, err := f.Read(b)
		if n != 0 || err != nil {
			t.Errorf("%v: Read(0) = %d, %v, want 0, nil", fs.Name(), n, err)
		}
		f.Seek(0, 0)
		b = make([]byte, 100)
		n, err = f.Read(b)
		if n <= 0 || err != nil {
			t.Errorf("%v: Read(100) = %d, %v, want >0, nil", fs.Name(), n, err)
		}
	}
}

func TestMemFileRead(t *testing.T) {

	Fss, c := initFileSystems()
	defer c.Close()
	for _, fs := range Fss {

		fsc, ok := memfs.Unwrap(fs)
		if !ok {
			return
		}

		f, err := fsc.Create("testfile")
		if err != nil {
			t.Errorf("MemFileRead - create failed %v -  %v", err, f)
		}

		f.WriteString("abcd")
		f.Seek(0, 0)
		b := make([]byte, 8)
		n, err := f.Read(b)
		if n != 4 {
			t.Errorf("didn't read all bytes: %v %v %v", n, err, b)
		}
		if err != nil {
			t.Errorf("err is not nil: %v %v %v", n, err, b)
		}
		n, err = f.Read(b)
		if n != 0 {
			t.Errorf("read more bytes: %v %v %v", n, err, b)
		}
		if err != io.EOF {
			t.Errorf("error is not EOF: %v %v %v", n, err, b)
		}
	}

}
