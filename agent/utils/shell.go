package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %s %s\nError: %s", name, strings.Join(args, " "), stderr.String())
	}
	return nil
}

func RunDockerExec(containerName string, command ...string) error {
	args := append([]string{"exec", "-i", containerName}, command...)
	return RunCommand("docker", args...)
}
