package main

import (
	"fmt"
	"os"

	"github.com/kubeopencode/kubeopencode/cmd"
)

// main is the entry point for the kubeopencode CLI tool.
// It delegates execution to the root cobra command defined in cmd/root.go.
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
