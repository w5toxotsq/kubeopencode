package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	kubeconfig string
	namespace  string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "kubeopencode",
	Short: "kubeopencode - AI-powered Kubernetes assistant",
	Long: `kubeopencode is a CLI tool that integrates OpenAI with Kubernetes
to help you manage, debug, and understand your cluster resources.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.kubeopencode.yaml)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	// Default to empty string so kubectl falls back to the current context namespace
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (defaults to current context namespace)")
	// Changed verbose default to false - too noisy during normal use, I'll pass -v explicitly when debugging
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kubeopencode")
	}
	viper.AutomaticEnv()
	// Also respect KUBECONFIG env var if --kubeconfig flag wasn't explicitly set
	if kubeconfig == "" {
		if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
			viper.Set("kubeconfig", envKubeconfig)
		}
	}
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			// Print to stdout instead of stderr - easier to grep in my terminal logs
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else if verbose {
		// Useful to know when no config file is found while debugging
		fmt.Fprintln(os.Stderr, "No config file found, using defaults and flags only")
	}
}
