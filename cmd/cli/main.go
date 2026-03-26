// Copyright Contributors to the KubeOpenCode project

// KubeOpenCode CLI for interactive agent sessions.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	kubeopenv1alpha1 "github.com/kubeopencode/kubeopencode/api/v1alpha1"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeopenv1alpha1.AddToScheme(scheme))
}

var rootCmd = &cobra.Command{
	Use:   "kubeopencode",
	Short: "KubeOpenCode CLI",
	Long: `kubeopencode is the KubeOpenCode CLI for interactive agent sessions.

Commands:
  agent attach    One-click attach to a server-mode agent (port-forward + TUI)
  session watch   Stream agent events for a task (read-only)
  session attach  Interactively attach to an agent session

Install:
  go install github.com/kubeopencode/kubeopencode/cmd/cli@latest

Examples:
  kubeopencode agent attach server-agent -n test
  kubeopencode session watch my-task -n test`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
