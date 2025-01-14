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

type SetCmdRunner func(cmdRunner *CommandRunner, err error)
type SetLaunchUrl func(url string)
type AddOutputLine func(line string)

func ShowInspectorDialog(window fyne.Window, mcpServer *configuration.McpServerDescription) {

	var cmdRunner *CommandRunner
	var launchUrl string = ""
	var startButton, launchButton *widget.Button

	inspectorArgs := []string{}
	inspectorArgs = append(inspectorArgs, "@modelcontextprotocol/inspector")
	inspectorArgs = mcpServer.Configuration.GetMcpInspectorArgs(inspectorArgs)
	fmt.Println(strings.Join(inspectorArgs, " "))

	// create a label with copy to hold the inspector args
	inspectorArgsLabel := NewTextWithCopy(fmt.Sprintf("npx %s", strings.Join(inspectorArgs, " ")), window)

	// Create output text widget and dialog
	outputText := widget.NewTextGrid()
	outputText.SetText("")

	// add button to launch the inspector URL
	launchButton = widget.NewButton("Launch Browser for MCP Inspector", func() {
		if launchUrl != "" {
			open.Run(launchUrl)
			if launchButton != nil {
				launchButton.Hide()
			}
		}
	})

	// make the launch button visible only if the URL is not empty
	launchButton.Hide()

	startButton = widget.NewButton("Start MCP Inspector", func() {
		doRunCommand("npx", inspectorArgs, func(theCmdRunner *CommandRunner, err error) {
			cmdRunner = theCmdRunner
			if err != nil {
				// show Dialog with error
				dialog.ShowError(err, window)
				return
			}
			if startButton != nil {
				startButton.Hide()
			}
		}, func(line string) {
			outputText.SetText(outputText.Text() + line + "\n")
		}, func(url string) {
			launchUrl = url
			// change the button text
			launchButton.SetText("Launch Browser: " + url)
			launchButton.Show()
		})
	})

	// create a vertical box with the start button and the launch button
	// only one button is visible at a time
	buttons := container.NewVBox(startButton, launchButton)

	// Create dialog with output and kill button
	content := container.NewBorder(inspectorArgsLabel, buttons, nil, nil,
		container.NewVScroll(outputText))
	d := dialog.NewCustom("MCP Inspector", "Stop & Close the MCP Inspector", content, window)
	d.Resize(fyne.NewSize(600, 400))
	d.Show()

	d.SetOnClosed(func() {
		cmdRunner.Kill()

		// display dialog box to tell the userr to close the browser
		dialog.ShowInformation("MCP Inspector", "Please close the browser to fullystop the MCP Inspector", window)
	})
}

func doRunCommand(name string, args []string, onSetCmdRunner SetCmdRunner, onOutput AddOutputLine, onUrl SetLaunchUrl) {
	// Run command in goroutine
	go func() {
		cmdRunner, err := runCommand(name, args, onOutput, onUrl)
		if err != nil {
			onSetCmdRunner(nil, err)
			return
		}
		onSetCmdRunner(cmdRunner, nil)
		cmdRunner.cmd.Wait()
	}()
}

type CommandRunner struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func runCommand(name string, args []string, onOutput func(line string), onUrl func(url string)) (*CommandRunner, error) {
	fmt.Println("@@ runCommand", name, args)
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
