package main

import (
	"context"
	"io"

	"github.com/spf13/cobra"
	"kope.io/kc/pkg/cmd"
)

var (
	logs_long = longDescription(`
	logs
	`)

	logs_short = shortDescription(`logs>`)
)

func NewCmdLogs(f cmd.Factory, out io.Writer, stderr io.Writer) *cobra.Command {
	options := &cmd.LogsOptions{}

	cmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"ns"},
		Short:   logs_short,
		Long:    logs_long,
		Run: func(c *cobra.Command, args []string) {
			if err := completeCmdLogs(c, options); err != nil {
				exitWithError(err)
			}
			ctx := context.TODO()
			err := cmd.RunLogs(ctx, f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&options.Follow, "follow", "f", options.Follow, "Specify if the logs should be streamed.")

	return cmd
}

func completeCmdLogs(cmd *cobra.Command, o *cmd.LogsOptions) error {
	args := cmd.Flags().Args()
	if len(args) != 1 {
		return helpErrorf(cmd, "Unexpected args: %v", args)
	}

	o.Name = args[0]
	return nil
}
