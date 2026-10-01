// Package main provides the mogo CLI, a toolkit of small utilities backed by
// the github.com/grokify/mogo library.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mogo",
	Short: "mogo utility CLI",
	Long: `mogo is a command-line toolkit of small utilities backed by the
github.com/grokify/mogo Go library.

Commands:
  ssh  - SSH public key utilities`,
	SilenceErrors: true, // main() prints the error itself
	SilenceUsage:  true, // don't dump usage for business-logic errors like "no match"
}

func init() {
	rootCmd.Version = versionString()
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(versionCmd)
}
