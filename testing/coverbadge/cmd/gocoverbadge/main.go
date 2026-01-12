// Command gocoverbadge generates a Shields.io coverage badge for Go projects.
//
// It runs go test with coverage, excludes specified directories (like cmd/),
// and outputs a Markdown badge snippet.
//
// Usage:
//
//	gocoverbadge -dir ./ -out coverage.md
//	gocoverbadge -dir ./ -exclude cmd,examples
//	gocoverbadge -dir ./ -badge-only
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	dir := flag.String("dir", ".", "directory to run coverage on")
	out := flag.String("out", "", "output file for markdown snippet (stdout if empty)")
	exclude := flag.String("exclude", "cmd", "comma-separated directory prefixes to exclude (e.g., cmd,examples)")
	badgeOnly := flag.Bool("badge-only", false, "only output the badge markdown, no progress messages")
	label := flag.String("label", "coverage", "badge label text")
	flag.Parse()

	if !*badgeOnly {
		fmt.Println("Calculating coverage...")
	}

	// Step 1: Get list of packages, excluding specified directories
	packages, err := getPackages(*dir, *exclude)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing packages: %v\n", err)
		os.Exit(1)
	}

	if len(packages) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no packages found after exclusions\n")
		os.Exit(1)
	}

	if !*badgeOnly {
		fmt.Printf("Testing %d packages (excluding: %s)\n", len(packages), *exclude)
	}

	// Step 2: Run tests with coverage on selected packages
	coverFile := "coverage.out"
	args := append([]string{"test", "-coverprofile=" + coverFile}, packages...)
	testCmd := exec.Command("go", args...)
	testCmd.Dir = *dir
	if !*badgeOnly {
		testCmd.Stdout = os.Stdout
		testCmd.Stderr = os.Stderr
	}
	if err := testCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: tests failed: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Parse coverage
	coverCmd := exec.Command("go", "tool", "cover", "-func="+coverFile)
	coverCmd.Dir = *dir
	outBuf := &bytes.Buffer{}
	coverCmd.Stdout = outBuf
	coverCmd.Stderr = os.Stderr
	if err := coverCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: coverage analysis failed: %v\n", err)
		os.Exit(1)
	}

	// Extract total coverage
	re := regexp.MustCompile(`total:\s+\(statements\)\s+([\d.]+)%`)
	match := re.FindStringSubmatch(outBuf.String())
	if len(match) < 2 {
		fmt.Fprintf(os.Stderr, "Error: could not parse coverage output\n")
		os.Exit(1)
	}
	percentStr := match[1]
	percentVal, _ := strconv.ParseFloat(percentStr, 64)

	// Step 4: Determine badge color
	color := badgeColor(percentVal)

	// Step 5: Generate Markdown badge
	// URL-encode special characters for Shields.io
	badgePercent := strings.ReplaceAll(percentStr+"%", "%", "%25")
	badgeLabel := strings.ReplaceAll(*label, " ", "%20")
	badgeLabel = strings.ReplaceAll(badgeLabel, "-", "--")
	badgeLabel = strings.ReplaceAll(badgeLabel, "_", "__")
	markdown := fmt.Sprintf("![%s](https://img.shields.io/badge/%s-%s-%s)",
		*label, badgeLabel, badgePercent, color)

	// Step 6: Output
	if *out != "" {
		if err := os.WriteFile(*out, []byte(markdown+"\n"), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write file: %v\n", err)
			os.Exit(1)
		}
		if !*badgeOnly {
			fmt.Printf("Coverage: %s%% (%s)\n", percentStr, color)
			fmt.Printf("Badge written to: %s\n", *out)
		}
	} else {
		if !*badgeOnly {
			fmt.Printf("Coverage: %s%% (%s)\n", percentStr, color)
			fmt.Println()
		}
		fmt.Println(markdown)
	}
}

// getPackages returns a list of packages, excluding those matching the given prefixes.
func getPackages(dir, excludeStr string) ([]string, error) {
	// Get all packages
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse exclude patterns
	var excludes []string
	for _, e := range strings.Split(excludeStr, ",") {
		e = strings.TrimSpace(e)
		if e != "" {
			excludes = append(excludes, "/"+e)
		}
	}

	// Filter packages
	var packages []string
	for _, pkg := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if pkg == "" {
			continue
		}
		excluded := false
		for _, ex := range excludes {
			if strings.Contains(pkg, ex+"/") || strings.HasSuffix(pkg, ex) {
				excluded = true
				break
			}
		}
		if !excluded {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// badgeColor returns the Shields.io color based on coverage percentage.
func badgeColor(percent float64) string {
	switch {
	case percent >= 90:
		return "brightgreen"
	case percent >= 80:
		return "green"
	case percent >= 70:
		return "yellowgreen"
	case percent >= 60:
		return "yellow"
	case percent >= 50:
		return "orange"
	default:
		return "red"
	}
}
