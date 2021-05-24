// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/retgits/ashelper/exec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// localStartCmd represents the build command
var localStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a service locally",
	Run:   runStartLocal,
}

var (
	createNetwork bool
)

func init() {
	localCmd.AddCommand(localStartCmd)
	localStartCmd.Flags().BoolVar(&withBuild, "with-build", false, "Determines whether to build a container first.")
	localStartCmd.Flags().BoolVar(&withPush, "with-push", false, "Determines whether the resulting built container should be pushed to the registry (only used with --with-build).")
	localStartCmd.Flags().BoolVar(&createNetwork, "create-network", false, "Create a Docker bridged network.")
	localStartCmd.Flags().StringVar(&name, "service", "", "The service to start.")
}

// runStartLocal is the actual execution of the command
func runStartLocal(cmd *cobra.Command, args []string) {
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

	svc, err := findService(name, ashelper.Ashelper.Akkaserverless.Service)
	CheckIfError(err)

	if withBuild {
		doBuild(ashelper.Ashelper.Docker, svc, withPush)
	}

	if createNetwork {
		result, err := exec.Run("docker", "network", "create", "-d", "bridged", "akkasls")
		Info(string(result))
		CheckIfError(err)
	}

	result, err := exec.Run("docker", "run", "-d", "--name", "userfunction", "--hostname", "userfunction", "--network", "akkasls", fmt.Sprintf("%s/%s/%s:%s", ashelper.Ashelper.Docker.Registry, ashelper.Ashelper.Docker.Username, svc.Name, svc.Version))
	Info(string(result))
	CheckIfError(err)

	result, err = exec.Run("docker", "run", "-d", "--name", "proxy", "--network", "akkasls", "-p", "9000:9000", "--env", "USER_FUNCTION_HOST=userfunction", "gcr.io/akkaserverless-public/akkaserverless-proxy:0.7.0-beta.8", "-Dconfig.resource=dev-mode.conf", "-Dcloudstate.proxy.protocol-compatibility-check=false")
	Info(string(result))
	CheckIfError(err)
}

func findService(s string, svcs []Service) (Service, error) {
	for _, svc := range svcs {
		if strings.EqualFold(svc.Name, s) {
			return svc, nil
		}
	}
	return Service{}, fmt.Errorf("service %s is not declared in .ashelper.yaml", s)
}
