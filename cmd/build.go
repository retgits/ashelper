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

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a container for Akka Serverless using Docker",
	Run:   runBuild,
}

var (
	configFileLocation string
	withPush           bool
)

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&configFileLocation, "config-file", "c", "", "The location of the .ashelper.yaml file.")
	buildCmd.Flags().BoolVar(&withPush, "with-push", false, "Determines whether the resulting built container should be pushed to the registry.")
}

// runInit is the actual execution of the command
func runBuild(cmd *cobra.Command, args []string) {
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
		doBuild(ashelper.Ashelper.Docker, svc, withPush)
	}
}

func doBuild(docker Docker, svc Service, withPush bool) {
	var result []byte
	var err error

	if len(svc.Folder) == 0 {
		result, err = exec.Run("docker", "build", ".", "-t", fmt.Sprintf("%s/%s/%s:%s", docker.Registry, docker.Username, svc.Name, svc.Version))
	} else {
		result, err = exec.RunInDir(svc.Folder, "docker", "build", ".", "-t", fmt.Sprintf("%s/%s/%s:%s", docker.Registry, docker.Username, svc.Name, svc.Version))
	}

	Info("Build result: %s", string(result))
	CheckIfError(err)

	if withPush {
		result, err = exec.Run("docker", "push", fmt.Sprintf("%s/%s/%s:%s", docker.Registry, docker.Username, svc.Name, svc.Version))
		Info("Push result: %s", string(result))
		CheckIfError(err)
	}
}
