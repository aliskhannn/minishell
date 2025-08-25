package shell

import "os"

// IsBuiltin checks whether the given command name corresponds to a shell builtin.
// Currently supported builtins: "cd", "pwd", "echo", "kill", "ps".
func IsBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "pwd", "echo", "kill", "ps":
		return true
	default:
		return false
	}
}

// ExpandEnv expands environment variables in the input string.
func ExpandEnv(s string) string {
	return os.Expand(s, func(key string) string {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}

		return ""
	})
}
