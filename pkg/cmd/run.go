package cmd

import (
	"io"
	"os"
	"os/exec"
	"syscall"
)

type ShellKubectlOptions struct {
	Args []string
}

func RunShellKubectl(f Factory, out io.Writer, stderr io.Writer, o *ShellKubectlOptions) error {
	c := exec.Command("kubectl", o.Args...)

	c.Stdout = out
	c.Stderr = stderr

	err := c.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.Exited() {
				sys := exitErr.ProcessState.Sys()
				if waitStatus, ok := sys.(syscall.WaitStatus); ok {
					exitCode := waitStatus.ExitStatus()
					os.Exit(exitCode)
				}
			}
		}
		return err
	}
	return nil
}
