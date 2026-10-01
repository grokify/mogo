package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/grokify/mogo/crypto/sshutil"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH public key utilities",
}

var sshListCmd = &cobra.Command{
	Use:   "list [dir]",
	Short: "List .pub files in a directory with their SHA256 fingerprints",
	Long: `List every "*.pub" file directly inside dir (default: ~/.ssh) along
with its key type, SHA256 fingerprint, and comment.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSSHList,
}

var sshFindCmd = &cobra.Command{
	Use:   "find <sha256-fingerprint> [dir]",
	Short: "Find the .pub file matching a SHA256 fingerprint",
	Long: `Scan "*.pub" files directly inside dir (default: ~/.ssh) and print
the ones whose SHA256 fingerprint matches <sha256-fingerprint>. The
"SHA256:" prefix is optional.

Example:
  mogo ssh find SHA256:AbCdEfGhIjKlMnOpQrStUvWxYz0123456789ABCDEF`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runSSHFind,
}

func init() {
	sshCmd.AddCommand(sshListCmd)
	sshCmd.AddCommand(sshFindCmd)
}

func sshResolveDir(args []string, idx int) (string, error) {
	if len(args) > idx {
		return args[idx], nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".ssh"), nil
}

func printSSHScanErrs(cmd *cobra.Command, errs []error) {
	for _, e := range errs {
		fmt.Fprintln(cmd.ErrOrStderr(), "warning:", e)
	}
}

func runSSHList(cmd *cobra.Command, args []string) error {
	dir, err := sshResolveDir(args, 0)
	if err != nil {
		return err
	}

	keys, errs, err := sshutil.ScanDir(dir)
	if err != nil {
		return err
	}
	printSSHScanErrs(cmd, errs)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s  %-20s %-40s %s\n", k.Fingerprint, k.Type, k.Path, k.Comment)
	}
	return nil
}

func runSSHFind(cmd *cobra.Command, args []string) error {
	target := args[0]
	dir, err := sshResolveDir(args, 1)
	if err != nil {
		return err
	}

	matches, errs, err := sshutil.FindByFingerprint(dir, target)
	if err != nil {
		return err
	}
	printSSHScanErrs(cmd, errs)

	if len(matches) == 0 {
		return fmt.Errorf("no .pub file in %s matches fingerprint %s", dir, sshutil.NormalizeFingerprint(target))
	}
	for _, k := range matches {
		fmt.Fprintf(cmd.OutOrStdout(), "%-40s %-20s %s\n", k.Path, k.Type, k.Comment)
	}
	return nil
}
