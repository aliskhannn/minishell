package shell

import (
	"github.com/aliskhannn/minishell/internal/builtins"
	"strings"
)

type Shell struct {
}

func New() *Shell {
	return &Shell{}
}

func (s *Shell) ExecuteLine(line string) error {
	argv := strings.Fields(line)

	return builtins.Run(argv)
}
