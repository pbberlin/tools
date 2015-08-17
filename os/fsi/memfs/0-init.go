// Package memfs offers a filesystem in memory.
//
// It was taken from Steve Francia's afero and extended.
// Most importantly, memMapFs.fos (previously memMapFs.fos)
// now holds the full path of each directory.
// Before that, directory names had to be unique.
// Same for InMemoryFile.memDir.
// The type memDirMap was removed; just use builtin map semantics.
// The entire pathing logic was redone.
// There is no os-dependent filepath anymore;
// everything is unix forward slashed.
//
// Creation happens with New(options...)
//
// fileinfo.Name() now returns the basename of the file;
// just as os.FileInfo.Name() does.
//
// All internal usage of Name() had to be rewritten.
//
// The locking approach remains a mystery to me.
// There are multiple locks and mutexes in Remove/Rename/Open and Close.
// I kept them in place.
// Strangely, the InMemoryFile.memDir map is *not* synced at all,
// though I think it should.
// Strangely, Remove() did not unregister with parent.
//
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

const (
	sep = "/" // No support for windows
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

	fs := memMapFs{}
	ifs := fsi.FileSystem(&fs)
	_ = ifs

}
