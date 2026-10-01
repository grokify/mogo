package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

// Version information. These are stamped by goreleaser via ldflags for
// released binaries; for local `go build`/`go install` builds they keep their
// sentinel defaults and are filled from the Go build info embedded by the
// toolchain (see resolveVersion).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), versionString())
	},
}

// resolveVersion returns version details, preferring the ldflags values set by
// goreleaser and falling back to the build info the Go toolchain embeds
// automatically. For a plain `go install`/`go build` from a working tree this
// yields the module version (or "(devel)"), the commit, the build time, and
// whether the tree had uncommitted changes.
func resolveVersion() (ver, rev, when string, modified, devel bool) {
	ver, rev, when = version, commit, date

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ver, rev, when, false, ver == "dev"
	}

	if ver == "dev" { // not stamped by ldflags
		if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			ver = bi.Main.Version // e.g. installed via `@v0.5.0`
		} else {
			ver, devel = "(devel)", true
		}
	}

	var biRev, biTime string
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			biRev = s.Value
		case "vcs.time":
			biTime = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if rev == "none" && biRev != "" {
		rev = biRev
	}
	if when == "unknown" && biTime != "" {
		when = biTime
	}
	return ver, rev, when, modified, devel
}

// versionString renders the multi-line version report.
func versionString() string {
	ver, rev, when, modified, devel := resolveVersion()

	short := rev
	if len(short) > 12 {
		short = short[:12]
	}

	var b strings.Builder
	fmt.Fprintf(&b, "mogo %s", ver)
	if devel || modified {
		b.WriteString("  [development build]")
	}
	b.WriteByte('\n')
	if short != "" && short != "none" {
		fmt.Fprintf(&b, "  commit: %s", short)
		if modified {
			b.WriteString(" (modified)")
		}
		b.WriteByte('\n')
	}
	if when != "" && when != "unknown" {
		fmt.Fprintf(&b, "  built:  %s\n", when)
	}
	fmt.Fprintf(&b, "  go:     %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	return b.String()
}
