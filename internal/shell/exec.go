package shell

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Run executes a pipeline and handles logical operators (&& and ||).
// - If a pipeline fails and there is an OR branch (`||`), it continues there.
// - If a pipeline succeeds and there is an AND branch (`&&`), it continues there.
// Returns the last encountered error, or nil if all succeeded.
func (s *Shell) Run(p *Pipeline) error {
	cur := p
	var lastErr error

	for cur != nil {
		err := s.RunPipeline(cur)
		if err != nil {
			lastErr = err

			// If command failed but we have an "OR" pipeline, go there.
			if cur.OrNext != nil {
				cur = cur.OrNext
				continue
			}
			break
		} else {
			lastErr = nil

			// If command succeeded and we have an "AND" pipeline, go there.
			if cur.AndNext != nil {
				cur = cur.AndNext
				continue
			}
			break
		}
	}

	return lastErr
}

// RunPipeline executes a single pipeline (possibly multiple commands connected with pipes).
// It sets up pipes between processes, redirects input/output, and manages signal handling.
// If the command is a builtin and the pipeline has only one command, it is executed directly.
func (s *Shell) RunPipeline(p *Pipeline) error {
	var cmds []*exec.Cmd
	var prevStdout *os.File
	var files []*os.File // open files for redirection
	var pipes []*os.File // pipes between commands

	// Listen for SIGINT (Ctrl+C) so we can forward it to child processes.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	defer signal.Stop(sigCh)

	var currentCmds []*exec.Cmd

	for i, c := range p.Commands {
		// If it's a builtin and the only command in a pipeline: run directly
		if IsBuiltin(c.Name) && len(p.Commands) == 1 {
			return RunCommand(c)
		}

		// Create an external command process.
		cmd := exec.Command(c.Name, c.Args...)
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // put command in its own process group

		// --- Setup stdin ---
		if i == 0 {
			// The first command may have input redirection.
			if c.Input != "" {
				inFile, err := os.Open(c.Input)
				if err != nil {
					return err
				}
				cmd.Stdin = inFile
				files = append(files, inFile)
			}
		} else {
			// Subsequent command gets stdin from previous pipe.
			cmd.Stdin = prevStdout
		}

		// --- Setup stdout ---
		if i == len(p.Commands)-1 {
			// Last command may redirect output to file.
			if c.Output != "" {
				outFile, err := os.Create(c.Output)
				if err != nil {
					return err
				}
				cmd.Stdout = outFile
				files = append(files, outFile)
			} else {
				cmd.Stdout = os.Stdout
			}
		} else {
			// Create a pipe for communication with the next command.
			r, w, err := os.Pipe()
			if err != nil {
				return err
			}
			cmd.Stdout = w
			prevStdout = r
			pipes = append(pipes, w, r)
		}

		cmds = append(cmds, cmd)
		currentCmds = append(currentCmds, cmd)
	}

	// Start all commands.
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	// Close all pipe file descriptors in a parent process.
	for _, f := range pipes {
		_ = f.Close()
	}
	// Close redirection files.
	for _, f := range files {
		_ = f.Close()
	}
	// Ensure last read end of pipe is closed eventually.
	if prevStdout != nil {
		defer func() { _ = prevStdout.Close() }()
	}

	// Goroutine to handle SIGINT (Ctrl+C).
	done := make(chan struct{})
	go func() {
		for sig := range sigCh {
			if sig == syscall.SIGINT {
				if len(currentCmds) > 0 {
					for _, c := range currentCmds {
						if c.Process != nil {
							pgid, _ := syscall.Getpgid(c.Process.Pid)
							_ = syscall.Kill(-pgid, syscall.SIGINT)
						}
					}
				}
			}
		}
		close(done)
	}()

	// Wait for all commands in pipeline to finish.
	var err error
	for _, cmd := range cmds {
		if werr := cmd.Wait(); werr != nil {
			// Preserve the first error if multiple fail.
			if err == nil {
				err = werr
			}
		}
	}

	// Stop listening for signals and wait for goroutine to exit.
	close(sigCh)
	<-done

	return err
}
