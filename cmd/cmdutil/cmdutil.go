package cmdutil

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
)

// ExecSimple provides a simple interface to execute a system command.
func ExecSimple(command, stdoutFile, stderrFile string) (bytes.Buffer, bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout, stderr, err
	}
	stdoutFile = strings.TrimSpace(stdoutFile)
	stderrFile = strings.TrimSpace(stderrFile)
	if len(stdoutFile) > 0 {
		err := ioutil.WriteFile(stdoutFile, stdout.Bytes(), 0644)
		if err != nil {
			return stdout, stderr, err
		}
	}
	if len(stderrFile) > 0 {
		err := ioutil.WriteFile(stderrFile, stdout.Bytes(), 0644)
		if err != nil {
			return stdout, stderr, err
		}
	}
	return stdout, stderr, nil
}
