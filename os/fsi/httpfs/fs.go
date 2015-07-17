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
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/stringspb"
)

type HttpFs struct {
	SourceFs fsi.FileSystem
}

func (h HttpFs) Dir(s string) *httpDir {
	log.Printf("httpfs (base)dir %v", s)
	return &httpDir{basePath: s, fs: h}
}

func (h HttpFs) Name() string {
	return fmt.Sprintf("httpfs over %v %v", h.SourceFs.Name(), h.SourceFs.String())
}

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

	if strings.HasSuffix(name, "favicon.ico") {
		return nil, os.ErrNotExist
	}

	f, err := h.SourceFs.Open(name)
	if err == nil {

		// Gather som info
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		tp := "F"
		if stat.IsDir() {
			tp = "D"
		}
		fn := fmt.Sprintf("%v %v", f.Name(), tp)

		// report info
		log.Printf("httpfs open      %-22v fnd %-22v %v", name, fn, h.Name())

		// return fo as http.File
		if httpfile, ok := f.(http.File); ok {
			return httpfile, nil
		}
	}

	// otherwise: error logging
	log.Printf("httpfs open      %-22v     %-22v %v", name, "", h.Name())
	log.Printf("             err %-22v", stringspb.Ellipsoider(err.Error(), 24))

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
