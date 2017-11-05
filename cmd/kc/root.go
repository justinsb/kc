package main

import (
	goflag "flag"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kope.io/kc/pkg/cmd"
)

var (
	root_long = longDescription(`
	kc
	`)

	root_short = shortDescription(`kc`)
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func NewCmdRoot(out io.Writer, stderr io.Writer) *cobra.Command {
	options := &cmd.FactoryOptions{}

	home := homeDir()
	if home != "" {
		options.Kubeconfig = filepath.Join(home, ".kube", "config")
	}

	if os.Getenv("KUBECONFIG") != "" {
		options.Kubeconfig = os.Getenv("KUBECONFIG")
	}

	f := cmd.NewDefaultFactory(options)

	cmd := &cobra.Command{
		Use:   "kc",
		Short: root_short,
		Long:  root_long,
	}

	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	cmd.PersistentFlags().StringVar(&options.Kubeconfig, "kubeconfig", options.Kubeconfig, "Path to the kubeconfig file to use for CLI requests.")

	// create subcommands
	//cmd.AddCommand(NewCmdCompletion(f, out))
	cmd.AddCommand(NewCmdNamespace(f, out, stderr))
	cmd.AddCommand(NewCmdGet(f, out, stderr))
	cmd.AddCommand(NewCmdLogs(f, out, stderr))

	return cmd
}
