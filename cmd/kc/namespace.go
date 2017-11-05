package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/kc/pkg/cmd"
)

var (
	namespace_long = longDescription(`
	namespace <namespace>
	`)

	namespace_short = shortDescription(`namespace <namespace>`)
)

func NewCmdNamespace(f cmd.Factory, out io.Writer) *cobra.Command {
	options := &cmd.NamespaceOptions{}

	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"ns"},
		Short:   namespace_short,
		Long:    namespace_long,
		Run: func(c *cobra.Command, args []string) {
			if err := complete(c, options); err != nil {
				exitWithError(err)
			}
			err := cmd.RunNamespace(f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}

func complete(cmd *cobra.Command, o *cmd.NamespaceOptions) error {
	args := cmd.Flags().Args()
	if len(args) != 1 {
		return helpErrorf(cmd, "Unexpected args: %v", args)
	}

	o.Namespace = args[0]
	return nil
}
