package executil

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// ExecSimple provides a simple interface to execute a system command.
func ExecSimple(command string) (bytes.Buffer, bytes.Buffer, error) {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...) // #nosec G204
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return stdout, stderr, cmd.Run()
}

// ExecToFiles provides a simple interface to execute a system command.
// Redirects for STDOUT and STDERR must be passed in as file names,
// not as `>` and `2>` UNIX file descriptors.
func ExecToFiles(command, stdoutFile, stderrFile string, perm os.FileMode) (bytes.Buffer, bytes.Buffer, error) {
	stdout, stderr, err := ExecSimple(command)
	if err != nil {
		return stdout, stderr, err
	}
	stdoutFile = strings.TrimSpace(stdoutFile)
	stderrFile = strings.TrimSpace(stderrFile)
	if len(stdoutFile) > 0 {
		err := os.WriteFile(stdoutFile, stdout.Bytes(), perm)
		if err != nil {
			return stdout, stderr, err
		}
	}
	if len(stderrFile) > 0 {
		err := os.WriteFile(stderrFile, stdout.Bytes(), perm)
		if err != nil {
			return stdout, stderr, err
		}
	}
	return stdout, stderr, nil
}
