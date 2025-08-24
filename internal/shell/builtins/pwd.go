package builtins

import (
	"fmt"
	"os"
)

func buildinPWD() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: cannot get current directory: %w", err)
	}

	fmt.Println(dir)

	return nil
}
