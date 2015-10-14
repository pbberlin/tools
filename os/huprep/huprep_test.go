package main

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestBla(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	files := []string{"test.zip", "test.jpg"}

	for _, v := range files {

		curdir, _ := os.Getwd()
		filePath := path.Join(curdir, v)

		_ = filePath

	}
}
