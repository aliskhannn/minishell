# Minishell

**Minishell** is a simple, lightweight Unix shell written in Go, designed to emulate the behavior of a real shell. It supports built-in commands, external command execution, pipelines, conditional execution, environment variables, and input/output redirection.

This shell is implemented as a Go program and can be run interactively or via scripts.

---

## Features

### Built-in Commands

The shell provides several essential built-in commands:

* `cd <path>` – Change the current working directory. Supports `~` and `-` for home and previous directories.
* `pwd` – Print the current working directory.
* `echo <args>` – Print arguments to stdout. Supports:

    * `-n` flag to suppress the newline.
    * `-e` flag to interpret escape sequences such as `\n`, `\t`.
    * Single `'` and double `"` quotes.
    * Environment variable expansion for `$VAR`.
* `kill <pid>` – Send the `SIGTERM` signal to a process by its PID.
* `ps` – Display currently running processes with PID and command name.

### External Commands

Any command not recognized as a built-in is executed as an external process using `os/exec`, similar to a normal shell. For example: `ls`, `grep`, `sleep`, etc.

### Pipelines

Commands can be chained together using the pipe (`|`) operator, passing the output of one command to the input of the next.
Example:

```bash
ps | grep go | wc -l
```

### Conditional Execution

Supports the logical operators:

* `&&` – Execute the next command **only if the previous succeeds**.
* `||` – Execute the next command **only if the previous fails**.

Example:

```bash
echo hello && echo success
false || echo ok
```

### Input/Output Redirection

* `>` – Redirect stdout to a file (overwrite).
* `<` – Redirect stdin from a file.

Example:

```bash
echo hello > output.txt
cat < output.txt
```

### Environment Variables

Variables of the form `$VAR` are expanded automatically:

```bash
echo My home is $HOME
```

### Signal Handling

* **Ctrl+D (EOF)** – Exit the shell gracefully.
* **Ctrl+C (SIGINT)** – Interrupt the currently running command without closing the shell.

---

## Installation

Clone the repository:

```bash
git clone https://github.com/aliskhannn/minishell.git
cd minishell
```

Run the shell:

```bash
go run cmd/minishell/main.go
```

---

## Usage Examples

### Directory Navigation

```bash
cd example/path
cd -
cd ~
```

### Echo

```bash
echo hello                # hello
echo hello world          # hello world
echo -n hello             # hello (without \n)
echo line1\nline2         # prints literally
echo -e line1\nline2      # line1nline2
echo -e 'line1\nline2'    # interprets escapes in quotes
echo "Home is $HOME"      # Home is home/example
echo My home is $HOME     # Home is home/example
echo 'My home is $HOME'   # no expansion in single quotes
```

### Process Management

```bash
ps
kill 12345   # terminate process with PID 12345
```

### Pipelines

```bash
echo hello world | wc -w      # outputs: 2
ps | grep go | wc -l
```

### Redirection

```bash
echo hello > sample.txt
cat < sample.txt
```

### Conditional Execution

```bash
echo hello && echo success    # prints both lines
false || echo ok              # prints ok
```

### Combining Features

```bash
echo hello | wc -w > count.txt
cat < count.txt
echo "My home is $HOME" > home.txt
cat home.txt
```

### Signal Handling

```bash
sleep 10      # press Ctrl+C to interrupt without closing shell
# prints ^C
```

---

## Development Notes

* Written entirely in **Go** using `os/exec`, `syscall`, `bufio`, and `strings`.
* Commands are parsed into a `Pipeline` structure to support conditional execution and piping.
* Each external command runs in its own process group to allow proper signal forwarding.
* Built-ins are executed directly in Go, enabling features like `cd` and `echo` to affect the shell environment.

---

## Testing

Integration tests cover:

* Built-in commands (`cd`, `pwd`, `echo`, `kill`, `ps`)
* Pipelines and redirection
* Conditional operators (`&&`, `||`)
* Environment variable expansion
* Signal handling (Ctrl+C, Ctrl+D)

Run tests:

```bash
make test
```