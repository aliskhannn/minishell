package shell

import (
	"fmt"
	"strings"
	"unicode"
)

// Command represents a single shell command (e.g., "ls", "echo"),
// with its arguments and optional input/output redirection.
type Command struct {
	Name   string   // Command name, e.g., "ls"
	Args   []string // Command arguments
	Input  string   // Input redirection ("<")
	Output string   // Output redirection (">")
}

// Pipeline represents a sequence of commands connected via pipes (|),
// and conditional execution with AND (&&) or OR (||) operators.
type Pipeline struct {
	Commands []*Command // Commands in the current pipeline
	AndNext  *Pipeline  // Next pipeline to execute on success (&&)
	OrNext   *Pipeline  // Next pipeline to execute on failure (||)
}

// Parse takes a line of shell input and returns a Pipeline structure
// representing commands, pipes, and conditional execution.
func Parse(line string) (*Pipeline, error) {
	tokens := tokenize(line)

	return parseConditional(tokens)
}

// tokenize splits the input string into individual tokens,
// including commands, arguments, operators (|, &&, ||), and redirection symbols (<, >).
func tokenize(input string) []string {
	var tokens []string
	var b strings.Builder

	flush := func() {
		// Flush the current token buffer into the tokens slice.
		if b.Len() > 0 {
			tokens = append(tokens, b.String())
			b.Reset()
		}
	}

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsSpace(r) {
			// Space indicates token boundary.
			flush()
			continue
		}

		if r == '|' {
			flush()
			// Handle "||" as a single token
			if i+1 < len(runes) && runes[i+1] == '|' {
				tokens = append(tokens, "||")
				i++
			} else {
				tokens = append(tokens, "|")
			}
			continue
		}

		if r == '&' {
			flush()
			// Handle "&&" as a single token
			if i+1 < len(runes) && runes[i+1] == '&' {
				tokens = append(tokens, "&&")
				i++
			} else {
				tokens = append(tokens, "&") // standalone '&' not used in our shell
			}
			continue
		}

		if r == '>' || r == '<' {
			flush()
			tokens = append(tokens, string(r))
			continue
		}

		// Regular character: append to current token buffer.
		b.WriteRune(r)
	}

	flush() // flush last token

	return tokens
}

// indexOfCondOp returns the index of the first conditional operator (&& or ||)
// in the token slice, or -1 if none are found.
func indexOfCondOp(tokens []string) int {
	for i, t := range tokens {
		if t == "&&" || t == "||" {
			return i
		}
	}
	return -1
}

// parseConditional parses tokens into a Pipeline, handling conditional operators (&& and ||).
// It recursively splits tokens at the first conditional operator and builds linked Pipelines.
func parseConditional(tokens []string) (*Pipeline, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	idx := indexOfCondOp(tokens)
	if idx == -1 {
		// No conditional operator: just a single pipeline.
		return parsePipeline(tokens)
	}

	// Parse the segment before the first conditional operator.
	if idx == 0 {
		return nil, fmt.Errorf("syntax error near %q", tokens[idx])
	}
	left, err := parsePipeline(tokens[:idx])
	if err != nil {
		return nil, err
	}

	cur := left
	rest := tokens[idx:]

	// Iteratively process remaining segments: operator + next segment.
	for len(rest) > 0 {
		op := rest[0] // "&&" or "||"
		rest = rest[1:]

		if len(rest) == 0 {
			return nil, fmt.Errorf("syntax error: expected command after %q", op)
		}

		nextIdx := indexOfCondOp(rest)
		var seg []string
		if nextIdx == -1 {
			seg = rest
			rest = nil
		} else {
			seg = rest[:nextIdx]
			rest = rest[nextIdx:]
		}

		if len(seg) == 0 {
			return nil, fmt.Errorf("syntax error: empty command after %q", op)
		}

		right, err := parsePipeline(seg)
		if err != nil {
			return nil, err
		}

		// Link the current pipeline to the next based on the operator.
		if op == "&&" {
			cur.AndNext = right
		} else {
			cur.OrNext = right
		}
		cur = right
	}

	return left, nil
}

// parsePipeline converts tokens into a single Pipeline with multiple commands connected by pipes.
// For example, "echo hi | wc -w" will produce a Pipeline with two Commands.
func parsePipeline(tokens []string) (*Pipeline, error) {
	var cmds []*Command
	var current []string

	for _, token := range tokens {
		if token == "|" {
			// End of current command, create Command object.
			cmd := parseCommand(current)
			cmds = append(cmds, cmd)
			current = nil
		} else {
			current = append(current, token)
		}
	}

	if len(current) > 0 {
		// Last command in the pipeline.
		cmd := parseCommand(current)
		cmds = append(cmds, cmd)
	}

	return &Pipeline{Commands: cmds}, nil
}

// parseCommand converts a slice of tokens into a Command structure,
// extracting the command name, arguments, and input/output redirection.
func parseCommand(tokens []string) *Command {
	cmd := &Command{}
	i := 0

	for i < len(tokens) {
		token := tokens[i]

		switch token {
		case ">":
			// Output redirection.
			if i+1 < len(tokens) {
				cmd.Output = tokens[i+1]
				i += 2
			} else {
				i++
			}
		case "<":
			// Input redirection.
			if i+1 < len(tokens) {
				cmd.Input = tokens[i+1]
				i += 2
			} else {
				i++
			}
		default:
			//// Expand environment variables like $HOME.
			//if strings.HasPrefix(token, "$") {
			//	envName := token[1:]
			//	token = os.Getenv(envName)
			//}

			// First token is the command name, subsequent tokens are arguments.
			if cmd.Name == "" {
				cmd.Name = token
			} else {
				cmd.Args = append(cmd.Args, token)
			}
			i++
		}
	}
	return cmd
}
