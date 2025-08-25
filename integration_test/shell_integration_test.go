package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// runShell executes the minishell binary with given input lines
// and returns combined stdout/stderr output.
func runShell(t *testing.T, input string) string {
	t.Helper()

	cmd := exec.Command("../bin/minishell")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Send commands line by line.
	for _, line := range strings.Split(input, "\n") {
		_, _ = stdin.Write([]byte(line + "\n"))
	}
	_ = stdin.Close()

	if err = cmd.Wait(); err != nil {
		// Ignore exit errors caused by Ctrl+D EOF simulation.
		_ = "" // for linter
	}

	return out.String()
}

func TestEchoBuiltin(t *testing.T) {
	output := runShell(t, "echo hello\necho -n world\necho 'My $HOME'\n")
	if !strings.Contains(output, "hello") {
		t.Error("echo hello failed")
	}
	if !strings.Contains(output, "world") {
		t.Error("echo -n world failed")
	}
	if !strings.Contains(output, "My $HOME") {
		t.Error("echo with single quotes failed")
	}
}

func TestPwdAndCd(t *testing.T) {
	home, _ := os.UserHomeDir()
	output := runShell(t, "cd ~\npwd\n")
	if !strings.Contains(output, home) {
		t.Errorf("expected %q in output, got %q", home, output)
	}
}

func TestPipeline(t *testing.T) {
	output := runShell(t, "echo hello world | wc -w\n")
	if !strings.Contains(output, "2") {
		t.Errorf("expected '2' in output, got %q", output)
	}
}

func TestConditionalOperators(t *testing.T) {
	output := runShell(t, "echo a && echo b\nfalse || echo ok\n")
	if !strings.Contains(output, "b") || !strings.Contains(output, "ok") {
		t.Errorf("conditional operators failed, got %q", output)
	}
}

func TestRedirection(t *testing.T) {
	file := "testfile.txt"
	defer func() {
		_ = os.Remove(file)
	}()

	runShell(t, "echo hello > "+file+"\n")
	data, _ := os.ReadFile(file)
	if strings.TrimSpace(string(data)) != "hello" {
		t.Errorf("file redirection failed, got %q", string(data))
	}
}

func TestPsAndKill(t *testing.T) {
	// Start a separate process to kill
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid

	// Use shell to kill that PID
	runShell(t, fmt.Sprintf("kill %d\n", pid))

	// Wait for process to exit
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(2 * time.Second):
		t.Errorf("process %d did not exit after kill", pid)
	case err := <-done:
		if err != nil {
			t.Logf("process %d exited with error: %v (expected for killed process)", pid, err)
		}
	}
}
