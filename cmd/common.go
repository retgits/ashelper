// Package cmd defines and implements command-line commands and flags
// used by ashelper. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"os"
)

type ASHelper struct {
	Ashelper Ashelper `json:"ashelper"`
}

type Ashelper struct {
	Akkaserverless Akkaserverless `json:"akkaserverless"`
	Docker         Docker         `json:"docker"`
}

type Akkaserverless struct {
	Project string    `json:"project"`
	Service []Service `json:"service"`
}

type Service struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Folder  string `json:"folder"`
}

type Docker struct {
	Registry string `json:"registry"`
	Username string `json:"username"`
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
