// +build suite3
// go test -tags=suite3

package aetest

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

import "testing"

func TestRename(t *testing.T) {
	for _, fs := range Fss {
		from, to := testDir+"/renamefrom", testDir+"/renameto"
		fs.Remove(to)              // Just in case.
		fs.MkdirAll(testDir, 0777) // Just in case.
		file, err := fs.Create(from)
		if err != nil {
			t.Fatalf("open %q failed: %v", to, err)
		}
		if err = file.Close(); err != nil {
			t.Errorf("close %q failed: %v", to, err)
		}
		err = fs.Rename(from, to)
		if err != nil {
			t.Fatalf("rename %q, %q failed: %v", to, from, err)
		}
		defer fs.Remove(to)
		_, err = fs.Stat(to)
		if err != nil {
			t.Errorf("stat %q failed: %v", to, err)
		}
	}
}
