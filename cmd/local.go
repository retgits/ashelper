// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"github.com/spf13/cobra"
)

// localCmd represents the build command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Starts or stops a service running locally",
}

func init() {
	rootCmd.AddCommand(localCmd)
}
