// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"github.com/retgits/ashelper/exec"
	"github.com/spf13/cobra"
)

// localStopCmd represents the build command
var localStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stops a running service locally",
	Run:   runstopLocal,
}

var (
	destroyNetwork bool
)

func init() {
	localCmd.AddCommand(localStopCmd)
	localStopCmd.Flags().BoolVar(&destroyNetwork, "destroy-network", false, "Removes the Docker bridged network.")
}

// runstopLocal is the actual execution of the command
func runstopLocal(cmd *cobra.Command, args []string) {
	result, err := exec.Run("docker", "stop", "userfunction")
	Info(string(result))
	CheckIfError(err)

	result, err = exec.Run("docker", "rm", "userfunction")
	Info(string(result))
	CheckIfError(err)

	result, err = exec.Run("docker", "stop", "proxy")
	Info(string(result))
	CheckIfError(err)

	result, err = exec.Run("docker", "rm", "proxy")
	Info(string(result))
	CheckIfError(err)

	if destroyNetwork {
		result, err = exec.Run("docker", "network", "rm", "akkasls")
		Info(string(result))
		CheckIfError(err)
	}
}
