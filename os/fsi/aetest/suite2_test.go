// +build suite2
// go test -tags=suite2

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

import (
	"bytes"
	"io"
	"os"
	"syscall"
	"testing"
)

func TestTruncate(t *testing.T) {

	for _, fs := range Fss {
		f := newFile("TestTruncate", fs, t)
		defer fs.Remove(f.Name())
		defer f.Close()

		checkSize(t, f, 0)
		f.Write([]byte("hello, world\n"))
		checkSize(t, f, 13)
		f.Truncate(10)
		checkSize(t, f, 10)
		f.Truncate(1024)
		checkSize(t, f, 1024)
		f.Truncate(0)
		checkSize(t, f, 0)
		_, err := f.Write([]byte("surprise!"))
		if err == nil {
			checkSize(t, f, 13+9) // wrote at offset past where hello, world was.
		}
	}
}

func TestSeek(t *testing.T) {

	for _, fs := range Fss {
		f := newFile("TestSeek", fs, t)
		defer fs.Remove(f.Name())
		defer f.Close()

		const data = "hello, world\n"
		io.WriteString(f, data)

		type test struct {
			in     int64
			whence int
			out    int64
		}
		var tests = []test{
			{0, 1, int64(len(data))},
			{0, 0, 0},
			{5, 0, 5},
			{0, 2, int64(len(data))},
			{0, 0, 0},
			{-1, 2, int64(len(data)) - 1},
			{1 << 33, 0, 1 << 33},
			{1 << 33, 2, 1<<33 + int64(len(data))},
		}
		for i, tt := range tests {
			off, err := f.Seek(tt.in, tt.whence)
			if off != tt.out || err != nil {
				if e, ok := err.(*os.PathError); ok && e.Err == syscall.EINVAL && tt.out > 1<<32 {
					// Reiserfs rejects the big seeks.
					// http://code.google.com/p/go/issues/detail?id=91
					break
				}
				t.Errorf("#%d: Seek(%v, %v) = %v, %v want %v, nil", i, tt.in, tt.whence, off, err, tt.out)
			}
		}
	}
}

func TestReadAt(t *testing.T) {

	for _, fs := range Fss {
		f := newFile("TestReadAt", fs, t)
		defer fs.Remove(f.Name())
		defer f.Close()

		const data = "hello, world\n"
		io.WriteString(f, data)

		b := make([]byte, 5)
		n, err := f.ReadAt(b, 7)
		if err != nil || n != len(b) {
			t.Fatalf("ReadAt 7: %d, %v", n, err)
		}
		if string(b) != "world" {
			t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
		}
	}
}

func TestWriteAt(t *testing.T) {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered in TestWriteAt", r)
	// 	}
	// }()

	for _, fs := range Fss {

		f := newFile("TestWriteAt", fs, t)
		defer fs.Remove(f.Name())
		defer f.Close()

		const data = "hello, world\n"
		io.WriteString(f, data)

		n, err := f.WriteAt([]byte("WORLD"), 7)
		if err != nil || n != 5 {
			t.Fatalf("WriteAt 7: %d, %v", n, err)
		}

		err = f.Sync()
		if err != nil {
			t.Fatalf("Saving file %v: %v", f.Name(), err)
		}
		// log.Printf(" Saved file %v", f.Name())

		f2, err := fs.Open(f.Name())
		if err != nil {
			t.Fatalf("Reopening file %v: %v", f.Name(), err)
		}
		defer f2.Close()

		// log.Printf("so far 4 %v", fs.Name())

		buf := new(bytes.Buffer)
		buf.ReadFrom(f2)
		b := buf.Bytes()

		if err != nil {
			t.Fatalf("%v: ReadFile %s: %v", fs.Name(), f.Name(), err)
		}
		if string(b) != "hello, WORLD\n" {
			t.Fatalf("after write: have %q want %q", string(b), "hello, WORLD\n")
		}

	}
}

//func TestReaddirnames(t *testing.T) {
//for _, fs := range Fss {
//testReaddirnames(fs, ".", dot, t)
////testReaddirnames(sysdir.name, fs, sysdir.files, t)
//}
//}

//func TestReaddir(t *testing.T) {
//for _, fs := range Fss {
//testReaddir(fs, ".", dot, t)
////testReaddir(sysdir.name, fs, sysdir.files, t)
//}
//}
