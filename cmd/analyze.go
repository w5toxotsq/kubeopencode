package cmd

import (
	"context"
	"fmt"

	"github.com/kubeopencode/kubeopencode/internal/analyzer"
	"github.com/kubeopencode/kubeopencode/internal/k8s"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resourceType string
	resourceName string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze Kubernetes resources using AI",
	Long:  `Analyze Kubernetes resources and get AI-powered insights and recommendations.`,
	Example: `  kubeopencode analyze --type pod --name my-pod
  kubeopencode analyze --type deployment --name my-deployment -n production`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		client, err := k8s.NewClient(viper.GetString("kubeconfig"))
		if err != nil {
			return fmt.Errorf("failed to create kubernetes client: %w", err)
		}

		resource, err := client.GetResource(ctx, resourceType, resourceName, namespace)
		if err != nil {
			return fmt.Errorf("failed to get resource: %w", err)
		}

		apiKey := viper.GetString("openai_api_key")
		if apiKey == "" {
			return fmt.Errorf("OPENAI_API_KEY is not set")
		}

		a := analyzer.New(apiKey)
		result, err := a.Analyze(ctx, resource)
		if err != nil {
			return fmt.Errorf("analysis failed: %w", err)
		}

		fmt.Println(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	// Default to "deployment" since that's what I use most often in my clusters
	analyzeCmd.Flags().StringVarP(&resourceType, "type", "t", "deployment", "resource type (pod, deployment, service, etc.)")
	analyzeCmd.Flags().StringVar(&resourceName, "name", "", "resource name")
	analyzeCmd.MarkFlagRequired("name")
	viper.BindEnv("openai_api_key", "OPENAI_API_KEY")
}
