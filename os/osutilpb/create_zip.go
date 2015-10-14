package osutilpb

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"log"
)

func CreateZipFile(files []string, archiveName string) {

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	for _, fn := range files {
		f, err := w.Create(fn)
		if err != nil {
			log.Fatal(err)
		}
		bts, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(bts)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(archiveName, buf.Bytes(), 0777)
	if err != nil {
		log.Fatal(err)
	}

}
