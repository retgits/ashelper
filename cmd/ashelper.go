// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ashelper",
	Short: "Akka Serverless Helper",
}

var (
	// Version of ashelper
	Version = "dev"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Set the version and template function to render the version text
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("\nashelper {{ .Version }}\n\n")
}
