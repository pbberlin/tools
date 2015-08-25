package weedout

import (
	"bytes"
	"fmt"

	"github.com/pbberlin/tools/os/osutilpb"
)

func similaritiesToFile(logdir string, frags []TextifiedTree, stage int) {

	// bfrags := stringspb.IndentedDumpBytes(frags)
	b := new(bytes.Buffer)
	for _, v := range frags {
		b.WriteString(fmt.Sprintf("%v %2v ", v.SourceID, v.Lvl))
		b.WriteString(fmt.Sprintf("%-8v             ", v.Outline))
		b.Write(v.Text)
		b.WriteString("\n")
		for _, v1 := range v.Similars {
			b.WriteString(fmt.Sprintf("%v %2v ", v1.SourceID, v1.Lvl))
			b.WriteString(fmt.Sprintf("%-8v    ", string(v1.Outline)))
			b.WriteString(spf("%2v ", v1.AbsLevenshtein))
			b.WriteString(spf("%-5.2v ", v1.RelLevenshtein))
			b.Write(v1.Text)
			b.WriteByte(10)
		}
		b.WriteByte(10)
	}
	osutilpb.Bytes2File(spf("%v/outp_frags_st%v.txt", logdir, stage), b.Bytes())

}
