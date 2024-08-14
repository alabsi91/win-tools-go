package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/huh"
	"github.com/goccy/go-yaml"
)

// AssetsPath is the path to the assets directory which lives alongside the executable
var AssetsPath = func() string {
	execPath, err := os.Executable()
	if err != nil {
		Log.Error(fmt.Sprintf("error getting executable path: %v\n", err))
		return ""
	}

	execDir := filepath.Dir(execPath)

	return filepath.Join(execDir, "assets")
}()

// ConfigYamlType defines the structure of the config YAML file
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

// ReadConfigFile reads and unmarshal the config file of type YAML
//   - On error, the program will exit with a fatal error message
//
// Returns: ConfigYamlType
func ReadConfigFile(path string) ConfigYamlType {
	dat, err := os.ReadFile(path)
	if err != nil {
		Log.Fatal(fmt.Sprintf(`Failed to read config file: "%s"`, path))
	}

	var config ConfigYamlType

	if err := yaml.Unmarshal(dat, &config); err != nil {
		Log.Fatal(fmt.Sprintf(`Failed to unmarshal config file: "%s"`, path))
	}

	return config
}

// AskForConfigFilePath prompts the user to enter a config file path
//   - Accepts relative and absolute paths
//   - Validates if the path exists
//   - Returns an error if the user cancels the prompt
func AskForConfigFilePath() (string, error) {
	var results string

	println("")

	validate := func(str string) error {
		if !IsPathExists(str) {
			return errors.New("file not found. Please enter a valid path")
		}
		return nil
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Please enter your config file path").
				Placeholder("Example: F:\\config.yaml").
				Validate(validate).
				Value(&results),
		),
	).Run()

	return results, err
}

// PreparePathsString prepares paths with environment variables by replacing them with their values
func PreparePathsString(paths []string) {
	r := regexp.MustCompile("%.+%")

	for i, path := range paths {
		paths[i] = r.ReplaceAllStringFunc(path, func(m string) string {
			return os.Getenv(m[1 : len(m)-1])
		})
	}
}

// IsPathExists checks if a path exists in the system
func IsPathExists(str string) bool {
	path := filepath.Clean(str)
	_, err := os.Stat(path)
	return err == nil
}
