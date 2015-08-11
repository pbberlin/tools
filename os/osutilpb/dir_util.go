package osutilpb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func DirOfExecutable() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	return dir
}
