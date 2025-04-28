package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sunbird-cli",
	Short: "Sunbird CLI for data products",
	Long: `A CLI tool for Sunbird data products that provides type-safe
operations for running jobs, replaying data, and managing configurations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Read in config file and ENV variables if set
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
