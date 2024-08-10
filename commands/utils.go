package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alabsi91/win-tools/app"
	"github.com/charmbracelet/huh"
	"github.com/goccy/go-yaml"
)

var Log = app.Log
var Utils = app.Utils
var Chocolatey = app.Chocolatey
var Powershell = app.Powershell

var AssetsPath = func() string {
	execPath, err := os.Executable()
	if err != nil {
		Log.Error(fmt.Sprintf("Error getting executable path: %v\n", err))
		return ""
	}

	execDir := filepath.Dir(execPath)

	return filepath.Join(execDir, "assets")
}()

type ConfigYamlType struct {
	Backup struct {
		Paths  []string
		Target string
	}
	EnvironmentVariables []struct {
		Key   string
		Value string
		Scope string
	} `yaml:"environmentVariables"`
	Packages []string
	Scripts  []string
}

func ReadConfigFile(path string) ConfigYamlType {
	dat, err := os.ReadFile(path)
	if err != nil {
		Log.Fatal(fmt.Sprintf(`Failed to read config file: "%s"`, path))
		os.Exit(1)
	}

	var config ConfigYamlType

	if err := yaml.Unmarshal(dat, &config); err != nil {
		Log.Fatal(fmt.Sprintf(`Failed to unmarshal config file: "%s"`, path))
		os.Exit(1)
	}

	return config
}

func AskForConfigFilePath() string {
	var results string

	validate := func(str string) error {
		if !Utils.IsPathExists(str) {
			return errors.New("file not found. Please enter a valid path")
		}
		return nil
	}

	huh.NewInput().
		Title("Please enter your config file path").
		Placeholder("Example: \"F:\\config.yaml\"").
		Validate(validate).
		Value(&results).Run()

	return results
}

// Replace all environment variables with their values
func PreparePathsString(paths []string) {
	r := regexp.MustCompile("%.+%")

	for i, path := range paths {
		paths[i] = r.ReplaceAllStringFunc(path, func(m string) string {
			return os.Getenv(m[1 : len(m)-1])
		})
	}
}
