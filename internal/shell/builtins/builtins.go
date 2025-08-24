package builtins

func Run(argv []string) error {
	if len(argv) == 0 {
		return nil
	}

	switch argv[0] {
	case "cd":
		return builtinCD(argv[1:])
	case "pwd":
		return buildinPWD()
	case "echo":
		return builtinEcho(argv[1:])
	default:
		return nil
	}
}
