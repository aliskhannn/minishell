package builtins

import (
	"fmt"
	"github.com/aliskhannn/minishell/internal/utils"
	"strings"
)

type echoFlags struct {
	Escape    bool
	NoNewLine bool
}

func builtinEcho(args []string) error {
	// If no arguments are provided, print a newline and return.
	if len(args) == 0 {
		fmt.Println()
		return nil
	}

	var flags echoFlags
	flags, args = parseFlags(args)

	// Determine if the arguments are quoted (single or double).
	singleQuoted := strings.HasPrefix(args[0], "'") && strings.HasSuffix(args[len(args)-1], "'")
	doubleQuoted := strings.HasPrefix(args[0], "\"") && strings.HasSuffix(args[len(args)-1], "\"")
	quoted := singleQuoted || doubleQuoted

	if flags.Escape && len(args) > 0 {
		argsStr := strings.Join(args, " ")     // join arguments into a single string
		unescaped := unescape(argsStr, quoted) // unescape the escape sequences
		args = strings.Split(unescaped, " ")
	}

	// If the arguments are enclosed in single quotes, join them and print as a single string
	// without interpreting any $ variables inside.
	if singleQuoted {
		joined := strings.Join(args, " ")
		trimmed := strings.Trim(joined, "'")
		fmt.Println(trimmed)

		return nil
	}

	// If the arguments are enclosed in double quotes, join them and then split by spaces
	// to handle environment variable expansion correctly.
	// For example, echo "Hello $USER" should expand $USER.
	if doubleQuoted {
		joined := strings.Join(args, " ")
		trimmed := strings.Trim(joined, "\"")
		args = strings.Split(trimmed, " ")
	}

	argsStr := strings.Join(args, " ")   // join arguments into a single string
	expanded := utils.ExpandEnv(argsStr) // expand environment variables

	if flags.NoNewLine {
		fmt.Print(expanded)
		return nil
	}

	fmt.Println(expanded) // print the final result

	return nil
}

func parseFlags(args []string) (echoFlags, []string) {
	var flags echoFlags
	var rest []string

	var allowedFlags = map[rune]bool{
		'e': true,
		'n': true,
	}

	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			rest = args[i:]
			break
		}

		valid := true
		for _, r := range arg[1:] {
			if _, ok := allowedFlags[r]; !ok {
				valid = false
				break
			}

			switch r {
			case 'e':
				flags.Escape = true
			case 'n':
				flags.NoNewLine = true
			}
		}

		if !valid {
			rest = args[i:]
			break
		}

		if i == len(args)-1 {
			rest = []string{}
		}
	}

	if len(rest) == 0 && len(args) > 0 {
		rest = args[len(args):]
	}

	return flags, rest
}

func unescape(s string, quoted bool) string {
	//If not quoted, just remove backslashes.
	if !quoted {
		var b strings.Builder
		i := 0

		for i < len(s) {
			// If a backslash is found,
			if s[i] == '\\' && i+1 < len(s) {
				b.WriteByte(s[i+1]) // skip it and add the next character
				i += 2
			} else {
				b.WriteByte(s[i]) // otherwise, just add the current character
				i++
			}
		}

		return b.String()
	}

	// If quoted, replace escape sequences with their actual characters.
	replacer := strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
		`\r`, "\r",
		`\a`, "\a",
		`\b`, "\b",
		`\033`, "\033", // ESC for color codes
		`\\`, "\\",
		`\"`, `"`,
		`\'`, "'",
	)
	return replacer.Replace(s)
}
