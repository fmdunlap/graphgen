package cmd

import (
	"github.com/spf13/cobra"
	"graphgen/internal/config"
	"os"
)

var EnvConfig *config.Config

var rootCmd = &cobra.Command{
	Use:   "graphgen",
	Short: "Graphgen CLI interface",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	EnvConfig = config.ParseConfig()
}
