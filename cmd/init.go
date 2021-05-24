// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Akka Serverless service with required parameters passed as parameters",
	Run:   runInit,
}

const (
	gitURL = "https://github.com/retgits/akkasls-templates"
)

var (
	language  string
	name      string
	outputDir string
	template  string

	// Supported languages by this tool
	languages = []string{"typescript", "nodejs"}
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&language, "language", "l", "", "The language you want to use for your service.")
	initCmd.Flags().StringVarP(&name, "name", "n", "", "The name you want to use for your service (will also be the name of the root directory in --output-dir).")
	initCmd.Flags().StringVar(&outputDir, "output-dir", "", "The location where the initialized service is output.")
	initCmd.Flags().StringVarP(&template, "template", "t", "", "The template you want to use for your service.")
	initCmd.MarkFlagRequired("language")
	initCmd.MarkFlagRequired("template")
	initCmd.MarkFlagRequired("name")
}

// runInit is the actual execution of the command
func runInit(cmd *cobra.Command, args []string) {
	if strings.EqualFold("javascript", language) {
		language = "nodejs"
	}

	// Check if the language is supported or not
	CheckIfError(languageExists(language))

	// Clone the templates into a temp directory
	tempDir, err := ioutil.TempDir("", "ashelper")
	CheckIfError(err)
	defer cleanup(tempDir)

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:   gitURL,
		Depth: 1,
	})
	CheckIfError(err)

	// Check if the template exists
	CheckIfError(templateExists(filepath.Join(tempDir, language), template))

	// Make sure there is an output dir we can use
	if len(outputDir) == 0 {
		currDir, _ := os.Getwd()
		outputDir = filepath.Join(currDir, name)
		Warning("outputDir was not provided, using %s directory as output", outputDir)
	}

	_, err = os.Stat(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(outputDir, 0755)
		} else {
			CheckIfError(err)
		}
	}

	// Copy the contents of the template to the output dir
	err = copy.Copy(filepath.Join(tempDir, language, template), outputDir)
	CheckIfError(err)
}

// templateExists checks if there is an existing template
func templateExists(dir string, template string) error {
	files, err := ioutil.ReadDir(dir)
	CheckIfError(err)

	tmpls := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			if strings.EqualFold(file.Name(), template) {
				return nil
			}
			tmpls = append(tmpls, file.Name())
		}
	}

	return fmt.Errorf("template [%s] is not a valid template. Available templates are [%s]", template, strings.Join(tmpls, ","))
}

// languageExists validates whether a language runtime is supported by this tool
func languageExists(l string) error {
	for _, language := range languages {
		if strings.EqualFold(language, l) {
			return nil
		}
	}
	return fmt.Errorf("language %s is not yet supported by ashelper", l)
}

func cleanup(path string) {
	Info("cleaning up %s", path)
	os.RemoveAll(path)
}
