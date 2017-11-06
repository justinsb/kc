package main

import (
	"context"
	"io"

	"fmt"

	"github.com/spf13/cobra"
	"kope.io/kc/pkg/cmd"
)

var (
	ssh_long = longDescription(`
	ssh
	`)

	ssh_short = shortDescription(`ssh>`)
)

func NewCmdSsh(f cmd.Factory, out io.Writer, stderr io.Writer) *cobra.Command {
	options := &cmd.SshOptions{}

	cmd := &cobra.Command{
		Use:   "ssh",
		Short: ssh_short,
		Long:  ssh_long,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("must specify name of node to SSH to")
			}
			options.NodeName = args[0]
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			ctx := context.TODO()

			// TODO: Unclear how to source the username.  One option is to look at the reported OS image,
			// (and maybe the cloud) but that might be too automagical
			options.Username = "admin"
			err := cmd.RunSsh(ctx, f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}
