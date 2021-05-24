// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Initializes the YAML template for ashelper",
	Run:   runGenerate,
}

var (
	project  string
	username string
	registry string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&project, "as-project", "", "The ID of the Akka Serverless project to deploy to.")
	generateCmd.Flags().StringVar(&username, "docker-user", "", "The username for your container registry.")
	generateCmd.Flags().StringVar(&registry, "docker-registry", "", "The registr for your containers.")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "The location where the YAML template is output.")
	generateCmd.MarkFlagRequired("as-project")
	generateCmd.MarkFlagRequired("docker-user")
}

// runGenerate is the actual execution of the command
func runGenerate(cmd *cobra.Command, args []string) {
	if len(registry) == 0 {
		registry = "docker.io"
	}
	// Generate the YAML specification
	spec := ASHelper{
		Ashelper: Ashelper{
			Akkaserverless: Akkaserverless{
				Project: project,
			},
			Docker: Docker{
				Username: username,
				Registry: registry,
			},
		},
	}

	d, err := yaml.Marshal(&spec)
	CheckIfError(err)

	// Make sure there is an output dir we can use
	if len(outputDir) == 0 {
		outputDir, _ = os.Getwd()
		Warning("outputDir was not provided, using %s directory as output", outputDir)
	}

	// Write to file
	file, err := os.Create(filepath.Join(outputDir, ".ashelper.yaml"))
	CheckIfError(err)
	_, err = file.Write(d)
	CheckIfError(err)
}
