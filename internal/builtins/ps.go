package builtins

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
)

// builtinPs выводит список запущенных процессов (PID и команду).
func builtinPs() error {
	procs, err := process.Processes()
	if err != nil {
		return err
	}

	fmt.Printf("%6s %s\n", "PID", "CMD")
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}
		fmt.Printf("%6d %s\n", p.Pid, name)
	}
	return nil
}
