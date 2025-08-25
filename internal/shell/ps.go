package shell

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

// builtinPs lists all currently running processes, similar to the "ps" command.
// It prints a simple table with PID and process name to standard output.
func builtinPs() error {
	// Retrieve all processes on the system using the gopsutil library.
	procs, err := process.Processes()
	if err != nil {
		return err
	}

	// Print table header with column names.
	fmt.Printf("%6s %s\n", "PID", "CMD")

	for _, p := range procs {
		// Get the name of the process.
		name, err := p.Name()
		if err != nil {
			continue
		}

		// Print the process PID and name
		// %6d ensures PID is right-aligned in 6-character width for neat formatting.
		fmt.Printf("%6d %s\n", p.Pid, name)
	}
	return nil
}
