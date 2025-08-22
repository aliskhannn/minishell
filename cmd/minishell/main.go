package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/aliskhannn/minishell/internal/shell"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	for {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasPrefix(dir, u.HomeDir) {
			dir = strings.Replace(dir, u.HomeDir, "~", 1)
		}

		boldGreen := "\033[1;32m"
		reset := "\033[0m"

		fmt.Printf("%s%s@%s%s:%s%s%s$ ", boldGreen, strings.ToLower(u.Name), host, reset, boldGreen, dir, reset)

		// Read a line from standard input
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				// End of file reached, exit the loop
				return
			}

			_, _ = fmt.Fprintln(os.Stderr, "error reading input:", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)

		if len(parts) > 0 {
			if parts[0] == "pwd" {
				dir, err := os.Getwd()
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, "error getting current directory:", err)
					continue
				}

				fmt.Println(dir)
			} else if parts[0] == "cd" {
				args := parts[1:]
				err := shell.BuiltinCD(args)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					continue
				}
			}
		}
	}
}
