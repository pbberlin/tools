package transposablematrix

import (
	"fmt"
	"os"
	"strings"

	"github.com/pbberlin/tools/uruntime"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf // plain and unchanged

var epf func(format string, a ...interface{}) error = fmt.Errorf

var pfOld func(format string, a ...interface{}) (int, error) = fmt.Printf // no more, old style

var pf = func(format string, a ...interface{}) (int, error) { // new - a wrapper

	s := spf(format, a...)

	if len(appStageLogs) == 0 || len(appStageLogs) < currStage-1 {
		fmt.Printf("Premature pf() or appStageLogs too small (%v,%v) %s\n", len(appStageLogs), currStage, s)
		uruntime.StackTrace(3)
		ExitWithLogDump()
	}

	buf := appStageLogs[currStage]

	if strings.HasSuffix(buf, "\n\n") {
		buf = buf[0 : len(buf)-1]
	}

	if strings.HasPrefix(s, "sameline") && strings.HasSuffix(buf, "\n") {
		buf = buf[0 : len(buf)-1]
		s = s[len("sameline"):]
	}

	buf += s
	appStageLogs[currStage] = buf

	// return fmt.Printf("wrapped: %q\n", spf(format, a...))
	return 0, nil
}

// Finally, the tools to disable all printing
// within a function and callees:
func pfDevNull(format string, a ...interface{}) (int, error) {
	return 0, nil // sucking void
}
func intermedPf(f func(format string, a ...interface{}) (int, error)) (g func(format string, a ...interface{}) (int, error)) {
	g = f         // return previous pf for temporal storage
	f = pfDevNull //
	return
}

func ExampleUsage() {
	pfTmp := intermedPf(pf)
	defer func() { pf = pfTmp }()
}

func ExitWithLogDump() {
	for i := 0; i < len(appStageLogs); i++ {
		fmt.Println("Stage", i, ":")
		fmt.Println(appStageLogs[i])
	}
	os.Exit(1)
}

func checkPrint(e error) {
	if e != nil {
		pf("%v\n", e)
	}
}
