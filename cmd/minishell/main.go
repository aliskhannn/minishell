package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

	"github.com/aliskhannn/minishell/internal/shell"
)

func main() {
	// Ignore the standard Ctrl+C so that the shell does not terminate.
	signal.Ignore(syscall.SIGINT)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	sh := shell.New()
	reader := bufio.NewReader(os.Stdin)

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// This is for shell-like behavior.
	// When you press Ctrl+C in shell,
	// it prints something like username@host:cwd$ ^C on each line.
	// So, minishell does so.
	go func() {
		for range sigCh {
			// Build and print the shell prompt (username@host:cwd$)
			prompt := makePrompt(u, host)
			fmt.Println()
			fmt.Print(prompt)
		}
	}()

	for {
		// Build and print the shell prompt (username@host:cwd$)
		prompt := makePrompt(u, host)
		fmt.Print(prompt)

		// Read one line from stdin. Handles Ctrl+D (EOF) and errors internally.
		line, err := readLine(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Ctrl+D was pressed at an empty prompt: exit gracefully.
				fmt.Println("\nexiting shell...")
				return
			}

			// Other input error: print it, but keep shell running
			_, _ = fmt.Fprintln(os.Stderr, "shell:", err)
			continue
		}

		// Execute the parsed command line
		if err := sh.ExecuteLine(line); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					if status.Signaled() && status.Signal() == syscall.SIGINT {
						// Ignore Ctrl+C.
						fmt.Println()
						continue
					}
				}
			}

			_, _ = fmt.Fprintln(os.Stderr, "shell:", err)
		}
	}
}

func makePrompt(u *user.User, host string) string {
	// Get current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Replace absolute home directory path with '~' (like in bash).
	if strings.HasPrefix(dir, u.HomeDir) {
		dir = strings.Replace(dir, u.HomeDir, "~", 1)
	}

	// Use ANSI escape codes for coloring (bold green for username and path).
	boldGreen := "\033[1;32m"
	reset := "\033[0m"

	// Example: "user@host:~/Projects$ "
	return fmt.Sprintf("%s%s@%s%s:%s%s%s$ ",
		boldGreen, strings.ToLower(u.Name), host, reset,
		boldGreen, dir, reset,
	)
}

func readLine(reader *bufio.Reader) (string, error) {
	// Read a line from standard input.
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Ctrl+D: return EOF so that main() can decide to exit.
			return "", io.EOF
		}

		return "", fmt.Errorf("error reading input: %w", err)
	}

	// Trim spaces and trailing newline before returning the command.
	return strings.TrimSpace(line), nil
}
