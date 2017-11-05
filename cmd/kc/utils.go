package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// longDescription is used to format long-description help text.  kubectl does i18n here.
func longDescription(s string) string {
	return s
}

// shortDescription is used to format short-description help text.  kubectl does i18n here.
func shortDescription(s string) string {
	return s
}

// helpErrorf prints usage help when returning an error from CLI parsing.  Based on k8s code.
func helpErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	cmd.Help()
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\n", msg)
}
