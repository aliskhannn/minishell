package builtins

import (
	"fmt"
	"strconv"
	"syscall"
)

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

func syscallKill(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
