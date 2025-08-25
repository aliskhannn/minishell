package shell

type Shell struct {
}

func New() *Shell {
	return &Shell{}
}

// ExecuteLine parses a single line of shell input and executes it.
// It first converts the line into a Pipeline (commands, pipes, conditionals)
// and then runs the pipeline.
func (s *Shell) ExecuteLine(line string) error {
	// Parse the line into a Pipeline structure.
	p, err := Parse(line)
	if err != nil {
		return err
	}

	// Execute the parsed pipeline.
	return s.Run(p)
}
