// Package memfs offers a filesystem in memory.

// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright 2013 tsuru authors. All rights reserved.
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
package memfs

import (
	"os"
	"sync"

	"github.com/pbberlin/tools/os/fsi"
)

var mux = &sync.Mutex{}

func init() {

	// forcing our implementations
	// to comply with our interfaces

	f := InMemoryFile{}
	ifa := fsi.File(&f)
	_ = ifa

	fi := InMemoryFileInfo{}
	ifi := os.FileInfo(&fi)
	_ = ifi

	fs := MemMapFs{}
	ifs := fsi.FileSystem(&fs)
	_ = ifs

}
