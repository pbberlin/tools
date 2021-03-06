// Package move_package scripts refactorings of package locations and names;
// sadly refactor/rename still has some limitations.
package main

import (
	"fmt"
	"go/build"

	"golang.org/x/tools/refactor/rename"
)

func main() {
	ct := build.Default
	var err error

	const pref = "github.com/pbberlin/tools/"

	tasks := [][]string{
		[]string{"htmlpb/", "pbhtml/"},
		[]string{"stringspb/", "stringspb/"},
		[]string{"ancestored_.../", "dsu/ancestored_.../"},
	}

	for idx, task := range tasks {
		fmt.Printf("%v: %v\n", idx, task)
		if len(task) != 2 {
			panic("need src and dst entry")
		}
		err = rename.Move(&ct, pref+task[0], pref+task[1], "git mv {{.Src}} {{.Dst}}")
		if err != nil {
			fmt.Printf("\t%v\n", err)
		}

	}

}
