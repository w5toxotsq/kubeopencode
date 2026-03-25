// Copyright Contributors to the KubeOpenCode project

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeopenv1alpha1 "github.com/kubeopencode/kubeopencode/api/v1alpha1"
	"github.com/kubeopencode/kubeopencode/internal/controller"
)

func init() {
	rootCmd.AddCommand(newAgentCmd())
}

func newAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Interact with KubeOpenCode agents",
	}
	cmd.AddCommand(newAgentAttachCmd())
	return cmd
}

func newAgentAttachCmd() *cobra.Command {
	var namespace string
	var localPort int

	cmd := &cobra.Command{
		Use:   "attach <agent-name>",
		Short: "Attach to a server-mode agent via OpenCode TUI",
		Long: `Attach to a server-mode agent with a single command.

This command:
  1. Verifies kubectl and opencode are installed
  2. Looks up the Agent CR to find the server deployment and port
  3. Starts kubectl port-forward in the background
  4. Launches opencode attach to connect to the agent
  5. Cleans up port-forward on exit

Examples:
  koc agent attach server-agent -n test
  koc agent attach my-agent -n production --local-port 5000`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgentAttach(cmd.Context(), namespace, args[0], localPort)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Agent namespace")
	cmd.Flags().IntVar(&localPort, "local-port", 0, "Local port for port-forward (default: same as server port)")

	return cmd
}

func runAgentAttach(ctx context.Context, namespace, agentName string, localPort int) error {
	// Step 1: Check prerequisites
	fmt.Println("🔍 Checking prerequisites...")

	if err := checkBinary("kubectl"); err != nil {
		return fmt.Errorf("kubectl not found: %w\n  Install: https://kubernetes.io/docs/tasks/tools/", err)
	}

	if err := checkBinary("opencode"); err != nil {
		return fmt.Errorf("opencode not found: %w\n  Install: https://opencode.ai", err)
	}

	// Step 2: Check kubeconfig / cluster connectivity
	fmt.Println("🔗 Connecting to cluster...")

	cfg, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("cannot connect to cluster: %w\n  Make sure KUBECONFIG is set or ~/.kube/config exists", err)
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Step 3: Look up Agent CR
	fmt.Printf("📋 Looking up agent %s/%s...\n", namespace, agentName)

	var agent kubeopenv1alpha1.Agent
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: agentName, Namespace: namespace}, &agent); err != nil {
		return fmt.Errorf("agent %q not found in namespace %q: %w", agentName, namespace, err)
	}

	if !controller.IsServerMode(&agent) {
		return fmt.Errorf("agent %q is not in Server mode (no serverConfig)\n  Only server-mode agents support interactive attach", agentName)
	}

	serverPort := controller.GetServerPort(&agent)
	deploymentName := controller.ServerDeploymentName(agentName)

	if localPort == 0 {
		localPort = int(serverPort)
	}

	// Check if the server is ready
	if agent.Status.ServerStatus == nil || agent.Status.ServerStatus.ReadyReplicas == 0 {
		return fmt.Errorf("agent %q server is not ready (0 ready replicas)\n  Check: kubectl get deployment %s -n %s", agentName, deploymentName, namespace)
	}

	fmt.Printf("✅ Agent found: %s (port %d, %d ready replicas)\n", deploymentName, serverPort, agent.Status.ServerStatus.ReadyReplicas)

	// Step 4: Check if local port is available
	if !isPortAvailable(localPort) {
		return fmt.Errorf("local port %d is already in use\n  Use --local-port to specify a different port", localPort)
	}

	// Step 5: Start port-forward
	fmt.Printf("🔀 Starting port-forward (localhost:%d → %s:%d)...\n", localPort, deploymentName, serverPort)

	pfCtx, pfCancel := context.WithCancel(ctx)
	defer pfCancel()

	pfCmd := exec.CommandContext(pfCtx, "kubectl", "port-forward",
		"-n", namespace,
		fmt.Sprintf("deployment/%s", deploymentName),
		fmt.Sprintf("%d:%d", localPort, serverPort),
	)
	pfCmd.Stderr = os.Stderr

	if err := pfCmd.Start(); err != nil {
		return fmt.Errorf("failed to start port-forward: %w", err)
	}

	// Wait for port-forward to be ready
	serverURL := fmt.Sprintf("http://localhost:%d", localPort)
	if err := waitForPort(localPort, 15*time.Second); err != nil {
		pfCancel()
		_ = pfCmd.Wait()
		return fmt.Errorf("port-forward failed to start: %w", err)
	}

	fmt.Printf("✅ Port-forward ready: %s\n", serverURL)

	// Step 6: Launch opencode attach
	fmt.Printf("🚀 Launching opencode attach...\n\n")

	attachCmd := exec.CommandContext(ctx, "opencode", "attach", serverURL)
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	attachCmd.Stderr = os.Stderr

	// Handle signals for cleanup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		pfCancel()
	}()

	err = attachCmd.Run()

	// Cleanup
	pfCancel()
	_ = pfCmd.Wait()

	if err != nil {
		// Don't show error for normal exit (user pressed Ctrl+C)
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("opencode attach exited with error: %w", err)
	}

	fmt.Println("\n👋 Session ended. Port-forward cleaned up.")
	return nil
}

// checkBinary verifies a binary is available in PATH.
func checkBinary(name string) error {
	_, err := exec.LookPath(name)
	return err
}

// isPortAvailable checks if a TCP port is available locally.
func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// waitForPort waits for a TCP port to become available.
func waitForPort(port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("port %d did not become available within %s", port, timeout)
}
