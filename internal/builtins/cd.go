package builtins

import (
	"fmt"
	"os"
	"strings"
)

var ErrTooManyArguments = fmt.Errorf("too many arguments")

// BuiltinCD implements the "cd" command for changing directories.
// Returns an error if more than one argument is provided or if the directory change fails.
func builtinCD(args []string) error {
	if len(args) > 1 {
		return ErrTooManyArguments
	}

	var path string
	if len(args) == 0 {
		path = "" // no argument provided, default to home directory
	} else {
		path = args[0]
	}

	// If the path is "-", change to the previous directory.
	if path == "-" {
		return chdirToPrevious()
	}

	return chdir(path)
}

// chdir changes the current working directory to the specified path.
func chdir(path string) error {
	// If the path is empty or just "~", change to the home directory.
	if path == "" || path == "~" {
		return chdirToHome()
	}

	// Get the current working directory to set OLDPWD.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cd: cannot get current directory: %w", err)
	}

	// Set OLDPWD to the current working directory before changing.
	_ = os.Setenv("OLDPWD", cwd)

	// If the path starts with "~/", replace it with the user's home directory.
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir() // get the user's home directory
		if err != nil {
			return fmt.Errorf("cd: cannot get user home: %s", err)
		}

		path = strings.Replace(path, "~", home, 1) // replace "~" with home directory
	}

	// Change to the specified directory.
	if err := changeDir(path); err != nil {
		return fmt.Errorf("cd: %w", err)
	}

	return nil
}

// chdirToPrevious changes the current working directory to the previous directory stored in OLDPWD.
func chdirToPrevious() error {
	// Get the previous directory from the OLDPWD environment variable.
	prevDir, ok := os.LookupEnv("OLDPWD")
	if !ok {
		return fmt.Errorf("cd: OLDPWD not set")
	}

	// Get the current working directory to set OLDPWD.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cd: cannot get current directory: %w", err)
	}

	// Set OLDPWD to the current working directory before changing.
	_ = os.Setenv("OLDPWD", cwd)

	// Change to the previous directory.
	if err := changeDir(prevDir); err != nil {
		return fmt.Errorf("cd: %w", err)
	}

	// Print the previous directory to stdout like the shell does.
	fmt.Println(prevDir)

	return nil
}

// chdirToHome changes the current working directory to the user's home directory.
func chdirToHome() error {
	// Get the current working directory to set OLDPWD.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cd: cannot get current directory: %w", err)
	}

	// Set OLDPWD to the current working directory before changing.
	_ = os.Setenv("OLDPWD", cwd)

	// Get the user's home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cd: cannot get user home: %s", err)
	}

	// Change to the home directory.
	if err := changeDir(home); err != nil {
		return fmt.Errorf("cd: %w", err)
	}

	return nil
}

// changeDir attempts to change the current working directory to the specified path.
func changeDir(path string) error {
	if err := os.Chdir(path); err != nil {
		msg := strings.Replace(err.Error(), "chdir ", "", 1) // Remove "chdir" prefix if present
		return fmt.Errorf("%s", msg)
	}

	return nil
}
