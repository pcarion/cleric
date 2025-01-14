package ui

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/pcarion/cleric/internal/configuration"
	"github.com/skratchdot/open-golang/open"
)

func ShowInspectorDialog(window fyne.Window, mcpServer *configuration.McpServerDescription) {
	inspectorArgs := []string{}
	inspectorArgs = append(inspectorArgs, "@modelcontextprotocol/inspector")
	inspectorArgs = mcpServer.Configuration.GetMcpInspectorArgs(inspectorArgs)
	fmt.Println(strings.Join(inspectorArgs, " "))

	// Create output text widget and dialog
	outputText := widget.NewTextGrid()
	outputText.SetText(">> MCP inspector for " + mcpServer.Name + ":\n")

	// Create dialog with output and kill button
	content := container.NewBorder(nil, nil, nil, nil,
		container.NewVScroll(outputText))
	d := dialog.NewCustom("MCP Inspector Output", "Close", content, window)
	d.Resize(fyne.NewSize(600, 400))
	d.Show()

	// Run command in goroutine
	go func() {
		cmdRunner, err := runCommand("npx", inspectorArgs, func(line string) {
			outputText.SetText(outputText.Text() + line + "\n")
		}, func(url string) {
			open.Run(url)
		})
		if err != nil {
			outputText.SetText(outputText.Text() + fmt.Sprintf("Error running inspector: %v\n", err))
			return
		}
		d.SetOnClosed(func() {
			cmdRunner.Kill()
		})

		cmdRunner.cmd.Wait()
		outputText.SetText(outputText.Text() + "Process completed.\n")
	}()
}

type CommandRunner struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func runCommand(name string, args []string, onOutput func(line string), onUrl func(url string)) (*CommandRunner, error) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
				line := string(buffer[:n])
				// check if line contains "MCP Inspector is up and running"
				if strings.Contains(line, "MCP Inspector is up and running") {
					// Find and extract URL from line
					if idx := strings.Index(line, "http://"); idx != -1 {
						// Extract from http:// until the next space or end of line
						urlEnd := strings.Index(line[idx:], " ")
						if urlEnd == -1 {
							url := strings.TrimSpace(line[idx:])
							if url != "" {
								onUrl(url)
							}
						} else {
							url := strings.TrimSpace(line[idx : idx+urlEnd])
							if url != "" {
								onUrl(url)
							}
						}
					}
				}
				onOutput(line)
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
				onOutput(string(buffer[:n]))
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

	return &CommandRunner{cmd: cmd, stdout: stdout, stderr: stderr}, nil
}

func (cr *CommandRunner) Kill() {
	cr.stdout.Close()
	cr.stderr.Close()
	if cr.cmd != nil && cr.cmd.Process != nil {
		// Kill the entire process group
		syscall.Kill(-cr.cmd.Process.Pid, syscall.SIGKILL)
		cr.cmd.Process.Kill()
	}
}
