package ui

import (
	"fmt"
	"os/exec"
)

func isValidServerName(name string) bool {
	// Check for empty string
	if len(name) == 0 {
		return false
	}

	// Check if name starts with a number
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}

	// Check if name contains only alphanumeric characters and underscore
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}
	return true
}

func runCommand(name string, args []string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)

	// Set up pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	// Handle stdout in a goroutine
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if n > 0 {
				fmt.Print(string(buffer[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// Handle stderr in a goroutine
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stderr.Read(buffer)
			if n > 0 {
				fmt.Print(string(buffer[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// To kill the process later:
	// cmd.Process.Kill()

	// To wait for completion:
	// err = cmd.Wait()

	return cmd, nil
}
