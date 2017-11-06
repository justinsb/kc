package main

import (
	"bytes"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"kope.io/kc/pkg/cmd"
)

type CompletionOptions struct {
	Shell string
}

func NewCmdCompletion(f cmd.Factory, out io.Writer) *cobra.Command {
	options := &CompletionOptions{}

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "completion",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCompletion(f, cmd, args, out, options)
			if err != nil {
				glog.Warningf("unexpected error: %v", err)
				exitWithError(err)
			}
		},
	}

	cmd.Flags().StringVar(&options.Shell, "shell", "", "target shell (bash).")

	return cmd
}

func RunCompletion(f cmd.Factory, c *cobra.Command, args []string, out io.Writer, options *CompletionOptions) error {
	//if len(args) != 0 {
	//	if c.Shell != "" {
	//		return fmt.Errorf("cannot specify shell both as a flag and a positional argument")
	//	}
	//	c.Shell = args[0]
	//}
	//
	//if c.Shell == "" {
	//	return fmt.Errorf("shell is required")
	//}
	//
	//switch c.Shell {
	//case "bash":
	//	return runCompletionBash(out, cmd.Parent())
	//
	//default:
	//	return fmt.Errorf("Unsupported shell type %q.", args[0])
	//}

	goflag.Set("logtostderr", "false")
	goflag.CommandLine.Parse([]string{})

	envCompCword := os.Getenv("COMP_CWORD")
	glog.Infof("COMP_CWORD=%q", envCompCword)

	if envCompCword == "" {
		return fmt.Errorf("COMP_CWORD not provided")
	}

	currentWordIndex, err := strconv.Atoi(envCompCword)
	if err != nil {
		return fmt.Errorf("cannot parse COMP_CWORD=%q", envCompCword)
	}

	var completions []string
	words := strings.Split(os.Getenv("COMP_WORDS"), " ")
	glog.Infof("COMP_WORDS=%v", words)

	if currentWordIndex == 1 {
		commands := []string{
			"namespace",
			"logs",
			"get",
			"ssh",
		}

		prefix := ""
		if len(words) >= 1 {
			prefix = words[1]
		}

		var matches []string
		for _, c := range commands {
			if strings.HasPrefix(c, prefix) {
				matches = append(matches, c)
			}
		}

		completions = matches
	} else if len(words) >= 2 {
		switch words[1] {
		case "namespace", "ns":
			prefix := ""
			if len(words) >= 3 {
				prefix = words[2]
			}
			completions, err = cmd.CompleteNamespaces(f, prefix)
		case "logs":
			prefix := ""
			if len(words) >= 3 {
				prefix = words[2]
			}
			completions, err = cmd.CompletePods(f, prefix)
		case "ssh":
			prefix := ""
			if len(words) >= 3 {
				prefix = words[2]
			}
			completions, err = cmd.CompleteNodes(f, prefix)
		}
	}

	if err != nil {
		return err
	}

	var b bytes.Buffer
	b.WriteString(strings.Join(completions, "\n"))
	//b.WriteString("aCOMP_WORDS=" +strings.Replace(os.Getenv("COMP_WORDS"), " ", "-", -1) + "\n")
	//b.WriteString("bCOMP_CWORD=" + strings.Replace(os.Getenv("COMP_CWORD"), " ", "-", -1) + "\n")
	//b.WriteString("cCOMP_LINE=" + strings.Replace(os.Getenv("COMP_LINE"), " ", "-", -1)  + "\n")

	_, err = b.WriteTo(out)
	if err != nil {
		return err
	}

	return nil
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

const bash = `
_complete_kc()
{
    local cur=${COMP_WORDS[COMP_CWORD]}
  	COMPREPLY=( $(COMP_CWORD="${COMP_CWORD}" COMP_WORDS="${COMP_WORDS[@]}" COMP_LINE="${COMP_LINE[@]}" kc completion do) )
}
complete -F _complete_kc kc
`
