package shell

//func parseCommandLine(args []string) ([]string, error) {
//	var parsedLine []string
//	var current strings.Builder
//	escaped := false
//	enableEscapes := false
//
//	if len(args) > 1 && args[1] == "-e" {
//		enableEscapes = true
//		args = args[2:]
//	}
//
//	var singleQuoted bool
//	var doubleQuoted bool
//
//	if enableEscapes {
//		singleQuoted = strings.HasPrefix(args[0], "'") && strings.HasSuffix(args[len(args)-1], "'")
//		doubleQuoted = strings.HasPrefix(args[0], "\"") && strings.HasSuffix(args[len(args)-1], "\"")
//	} else {
//		singleQuoted = strings.HasPrefix(args[1], "'") && strings.HasSuffix(args[len(args)-1], "'")
//		doubleQuoted = strings.HasPrefix(args[1], "\"") && strings.HasSuffix(args[len(args)-1], "\"")
//	}
//
//	line := strings.Join(args, " ")
//
//	for _, r := range line {
//		c := string(r)
//
//		if enableEscapes && escaped {
//			// If the character is escaped and we're not inside quotes,
//			if !singleQuoted && !doubleQuoted {
//				current.WriteRune(r) // add the escaped character as is
//			} else {
//				current.WriteString("\\" + c) // else keep the backslash for quoted strings
//			}
//
//			escaped = false
//			continue
//		}
//
//		switch c {
//		case "\\":
//			escaped = true
//
//		//case "'":
//		//	if !doubleQuoted {
//		//		singleQuoted = !singleQuoted
//		//		continue
//		//	}
//		//	current.WriteRune(r)
//		//
//		//case `"`:
//		//	if !singleQuoted {
//		//		doubleQuoted = !doubleQuoted
//		//		continue
//		//	}
//		//	current.WriteRune(r)
//
//		case " ":
//			if singleQuoted || doubleQuoted {
//				current.WriteRune(r)
//			} else if current.Len() > 0 {
//				parsedLine = append(parsedLine, current.String())
//				current.Reset()
//			}
//
//		default:
//			current.WriteRune(r)
//		}
//	}
//
//	if current.Len() > 0 {
//		parsedLine = append(parsedLine, current.String())
//	}
//
//	//if singleQuoted || doubleQuoted {
//	//	return nil, fmt.Errorf("unclosed quotes in command line")
//	//}
//
//	if singleQuoted {
//		joined := strings.Join(parsedLine, " ")
//		trimmed := strings.Trim(joined, "'")
//		parsedLine = strings.Split(trimmed, " ")
//		return parsedLine, nil
//	}
//
//	parsedLineStr := strings.Join(parsedLine, " ")
//	expanded := expandEnv(parsedLineStr)
//	parsedLine = strings.Fields(expanded)
//
//	return parsedLine, nil
//}
//
//func expandEnv(argsStr string) string {
//	expander := func(key string) string {
//		if val, ok := os.LookupEnv(key); ok {
//			return val
//		}
//
//		return ""
//	}
//
//	return os.Expand(argsStr, expander)
//}

//func builtinEcho(args []string) error {
//	// If no arguments are provided, print a newline and return.
//	if len(args) == 0 {
//		fmt.Println()
//		return nil
//	}
//
//	var flags echoFlags
//	flags, args = parseFlags(args)
//
//	// Determine if the arguments are quoted (single or double).
//	singleQuoted := strings.HasPrefix(args[0], "'") && strings.HasSuffix(args[len(args)-1], "'")
//	doubleQuoted := strings.HasPrefix(args[0], "\"") && strings.HasSuffix(args[len(args)-1], "\"")
//	quoted := singleQuoted || doubleQuoted
//
//	if flags.Escape && len(args) > 0 {
//		argsStr := strings.Join(args, " ")     // join arguments into a single string
//		unescaped := unescape(argsStr, quoted) // unescape the escape sequences
//		args = strings.Split(unescaped, " ")
//	}
//
//	// If the arguments are enclosed in single quotes, join them and print as a single string
//	// without interpreting any $ variables inside.
//	if singleQuoted {
//		joined := strings.Join(args, " ")
//		trimmed := strings.Trim(joined, "'")
//		fmt.Println(trimmed)
//
//		return nil
//	}
//
//	// If the arguments are enclosed in double quotes, join them and then split by spaces
//	// to handle environment variable expansion correctly.
//	// For example, echo "Hello $USER" should expand $USER.
//	if doubleQuoted {
//		joined := strings.Join(args, " ")
//		trimmed := strings.Trim(joined, "\"")
//		args = strings.Split(trimmed, " ")
//	}
//
//	argsStr := strings.Join(args, " ") // join arguments into a single string
//	expanded := expandEnv(argsStr)     // expand environment variables
//
//	if flags.NoNewLine {
//		fmt.Print(expanded)
//		return nil
//	}
//
//	fmt.Println(expanded) // print the final result
//
//	return nil
//}
