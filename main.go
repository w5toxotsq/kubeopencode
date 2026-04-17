package main

import (
	"fmt"
	"os"

	"github.com/kubeopencode/kubeopencode/cmd"
)

// main is the entry point for the kubeopencode CLI tool.
// It delegates execution to the root cobra command defined in cmd/root.go.
//
// Personal fork: using exit code 2 for usage/argument errors to better
// distinguish between general errors (1) and misuse of the CLI (2),
// following the convention used by many Unix tools (e.g. grep, curl).
//
// Note: I also considered exit code 64 (EX_USAGE from sysexits.h) but
// kept 2 for broader compatibility with shell scripting conventions.
//
// TODO: Look into wrapping errors with more context before printing,
// e.g. including the subcommand name so failures are easier to trace
// in scripts that chain multiple kubeopencode calls.
//
// TODO: Consider adding a --version flag that prints build info (git sha +
// build date) so I can quickly tell which version is running in my env.
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}
