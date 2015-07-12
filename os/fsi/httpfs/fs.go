// Package httpfs wraps any other fsi filesystem
// so that it works with http.FileServer;
// it adds as basePath to inner filesystem.
package httpfs

// Copyright © 2014 Steve Francia <spf@spf13.com>.
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
	"net/http"
	"os"

	"github.com/pbberlin/tools/os/fsi"
)

type HttpFs struct {
	SourceFs fsi.FileSystem
}

func (h HttpFs) Dir(s string) *httpDir {
	return &httpDir{basePath: s, fs: h}
}

func (h HttpFs) Name() string { return "h HttpFs" }

func (h HttpFs) Create(name string) (fsi.File, error) {
	return h.SourceFs.Create(name)
}

func (h HttpFs) Mkdir(name string, perm os.FileMode) error {
	return h.SourceFs.Mkdir(name, perm)
}

func (h HttpFs) MkdirAll(path string, perm os.FileMode) error {
	return h.SourceFs.MkdirAll(path, perm)
}

func (h HttpFs) Open(name string) (http.File, error) {
	f, err := h.SourceFs.Open(name)
	if err == nil {
		if httpfile, ok := f.(http.File); ok {
			return httpfile, nil
		}
	}
	return nil, err
}

func (h HttpFs) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return h.SourceFs.OpenFile(name, flag, perm)
}

func (h HttpFs) Remove(name string) error {
	return h.SourceFs.Remove(name)
}

func (h HttpFs) RemoveAll(path string) error {
	return h.SourceFs.RemoveAll(path)
}

func (h HttpFs) Rename(oldname, newname string) error {
	return h.SourceFs.Rename(oldname, newname)
}

func (h HttpFs) Stat(name string) (os.FileInfo, error) {
	return h.SourceFs.Stat(name)
}
