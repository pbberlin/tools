package main

import (
	"fmt"
	"go/build"

	"golang.org/x/tools/refactor/rename"
)

func main() {
	ct := build.Default
	var err error

	err = rename.Move(
		&ct, "github.com/pbberlin/tools/htmlpb/", "github.com/pbberlin/tools/pbhtml/", "git mv {{.Src}} {{.Dst}}")
	// &ct, "github.com/pbberlin/tools/pbhtml/", "github.com/pbberlin/tools/htmlpb/", "git mv {{.Src}} {{.Dst}}")
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = rename.Move(
		&ct, "github.com/pbberlin/tools/stringspb/", "github.com/pbberlin/tools/pbstrings/", "git mv {{.Src}} {{.Dst}}")
	// &ct, "github.com/pbberlin/tools/pbhtml/", "github.com/pbberlin/tools/htmlpb/", "git mv {{.Src}} {{.Dst}}")
	if err != nil {
		fmt.Printf("%v", err)
	}
}
