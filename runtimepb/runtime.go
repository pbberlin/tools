// Package uruntime gives a stracktrace and traces memory allocactions.
package runtimepb

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

func StackTrace(max int) {
	for i := 1; i <= max; i++ {
		// _, file, line, _ := runtime.Caller(i)
		line, file := LineFileXUp(i)
		log.Printf("        %s:%d ", file, line)
	}
}

func LineFileXUp(levelsUp int) (int, string) {
	_, file, line, _ := runtime.Caller(levelsUp + 1) // plus one for myself-func

	path1 := filepath.Dir(file)
	dir1 := filepath.Base(path1)

	path2 := filepath.Dir(path1)
	dir2 := filepath.Base(path2)

	file = filepath.Join(dir2, dir1, filepath.Base(file))
	return line, file
}

/*
	bookkeeping of memory allocation

*/
// alloc returns Memory
// use like this:
//		before := alloc()
//		after  := alloc()
func alloc() uint64 {
	var stats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&stats)
	// return stats.Alloc - uint64(unsafe.Sizeof(hs[0]))*uint64(cap(hs))
	return stats.Alloc
}

/*
	AllocLogger() takes care of negative memory changes.
	Without special treatment uint - uint
	may produce overflowed results.

	Otherwise we simplify the checkpointing.
	Simply call as follows

	fLogger, fDumper := runtimepb.AllocLogger()  // init and first checkpoint

	fLogger() // checkpoint
	...
	fLogger()
	...
	fDumper() // checkpoint and printing of results
*/
func AllocLogger() (func(), func() string) {
	logPoints := make([]uint64, 0, 10)
	logPoints = append(logPoints, alloc())
	fLogger := func() {
		logPoints = append(logPoints, alloc())
	}
	fDumper := func() string {
		logPoints = append(logPoints, alloc())
		msg := ""
		for k, v := range logPoints {
			if k == 0 {
				continue
			}
			var diff1 int64
			if v >= logPoints[k-1] {
				diff1 = int64(v - logPoints[k-1])
			} else {
				diff1 = int64(logPoints[k-1]-v) * -1
			}
			msg += fmt.Sprintf("P%02v: %6v | ", k, diff1)
			if k%10 == 0 {
				msg += "\n"
			}
		}
		if len(msg) > 1 {
			msg = msg[:len(msg)-2]
		}
		fmt.Printf("%s\n", msg)
		return msg
	}
	return fLogger, fDumper
}
