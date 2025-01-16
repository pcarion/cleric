//go:build windows

package ui

import (
	"os/exec"
	"syscall"
)

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func killProcess(cmd *exec.Cmd) {
	if cmd && cmd.Process != nil {
	cmd.Process.Kill()
}
