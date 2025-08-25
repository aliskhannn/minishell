package shell

import (
	"fmt"
	"strconv"
	"syscall"
)

// builtinKill sends a SIGTERM signal to the process with the given PID.
// It validates arguments, parses the PID, and delegates the actual signal
// sending to syscallKill.
func builtinKill(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("kill: missed PID")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid PID: %v", err)
	}

	if err = syscallKill(pid); err != nil {
		return err
	}

	return nil
}

// syscallKill sends a SIGTERM signal to the process with the given PID.
// It uses the low-level syscall.Kill function.
func syscallKill(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
