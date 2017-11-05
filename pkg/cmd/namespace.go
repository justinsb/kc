package cmd

import (
	"fmt"
	"io"

	"k8s.io/client-go/tools/clientcmd"
)

type NamespaceOptions struct {
	Namespace string
}

func RunNamespace(f Factory, out io.Writer, o *NamespaceOptions) error {
	if len(o.Namespace) == 0 {
		return fmt.Errorf("you must specify a non-empty namespace name")
	}

	configAccess, err := f.ConfigAccess()
	if err != nil {
		return err
	}

	config, err := configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	contextName := config.CurrentContext
	if contextName == "" {
		return fmt.Errorf("current-context is not set\n")
	}

	startingStanza, exists := config.Contexts[contextName]
	if !exists {
		return fmt.Errorf("no configuration is currently active in your configuration")
	}
	context := *startingStanza
	context.Namespace = o.Namespace
	config.Contexts[contextName] = &context

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		return fmt.Errorf("error modifying configuration: %v", err)
	}

	return nil
}
