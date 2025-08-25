package shell

import (
	"fmt"
	"os"
)

// buildinPWD prints the current working directory to standard output,
// similar to the "pwd" command in Unix shells.
func buildinPWD() error {
	// Retrieve the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: cannot get current directory: %w", err)
	}

	// Print the current working directory.
	fmt.Println(dir)

	return nil
}
