package shell

import (
	"fmt"
)

func RunCommand(c *Command) error {
	if c == nil || c.Name == "" {
		return nil
	}

	switch c.Name {
	case "cd":
		return builtinCD(c.Args)
	case "pwd":
		return buildinPWD()
	case "echo":
		return builtinEcho(c.Args, c.Output)
	case "ps":
		return builtinPs()
	case "kill":
		return builtinKill(c.Args)
	default:
		return fmt.Errorf("unknown builtin %q", c.Name)
	}
}
