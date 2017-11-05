package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/kc/pkg/cmd"
)

var (
	get_long = longDescription(`
	get
	`)

	get_short = shortDescription(`get`)
)

func NewCmdGet(f cmd.Factory, out io.Writer, stderr io.Writer) *cobra.Command {
	options := &cmd.ShellKubectlOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: get_short,
		Long:  get_long,
		Run: func(c *cobra.Command, args []string) {
			options.Args = []string{"get"}
			options.Args = append(options.Args, args...)

			err := cmd.RunShellKubectl(f, out, stderr, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}
