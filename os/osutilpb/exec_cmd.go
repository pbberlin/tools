package osutilpb

import (
	"log"
	"os/exec"
)

// http://stackoverflow.com/questions/10385551/get-exit-code-go
func ExecCmdWithExitCode(name string, args ...string) (success bool) {

	success = true

	command := exec.Command(name, args...)

	if err := command.Start(); err != nil {
		log.Fatalf("command.Start: %v", err)
	}

	if err := command.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Printf("exit err: %v  - %v %v", exiterr, name, args)
			// The program has exited with an exit code != 0
			success = false
		} else {
			log.Fatalf("command.Wait: %v, success dubious", err)
		}
	}

	return
}
