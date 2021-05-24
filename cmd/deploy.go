// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/retgits/ashelper/exec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// deployCmd represents the build command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a service to Akka Serverless",
	Run:   runDeploy,
}

var (
	withBuild bool
)

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&configFileLocation, "config-file", "c", "", "The location of the .ashelper.yaml file.")
	deployCmd.Flags().BoolVar(&withBuild, "with-build", false, "Determines whether to build a container first.")
	deployCmd.Flags().BoolVar(&withPush, "with-push", false, "Determines whether the resulting built container should be pushed to the registry (only used with --with-build).")
}

// runInit is the actual execution of the command
func runDeploy(cmd *cobra.Command, args []string) {
	if len(configFileLocation) == 0 {
		configFileLocation, _ = os.Getwd()
	}

	data, err := ioutil.ReadFile(filepath.Join(configFileLocation, ".ashelper.yaml"))
	CheckIfError(err)

	ashelper := ASHelper{}
	yaml.Unmarshal(data, &ashelper)

	if len(ashelper.Ashelper.Akkaserverless.Service) == 0 {
		CheckIfError(fmt.Errorf("no services found in .ashelper.yaml"))
	}

	for _, svc := range ashelper.Ashelper.Akkaserverless.Service {
		var result []byte

		if withBuild {
			doBuild(ashelper.Ashelper.Docker, svc, withPush)
		}

		result, err = exec.Run("akkasls", "services", "deploy", svc.Name, fmt.Sprintf("%s/%s/%s:%s", ashelper.Ashelper.Docker.Registry, ashelper.Ashelper.Docker.Username, svc.Name, svc.Version), "--project", ashelper.Ashelper.Akkaserverless.Project)
		Info("Deploy result: %s", string(result))
		CheckIfError(err)
	}
}
