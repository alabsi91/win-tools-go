package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

func askForSavePath() (string, error) {
	var results string

	validate := func(str string) error {
		if !strings.HasSuffix(str, ".yaml") {
			return errors.New("file extension must be .yaml")
		}
		return nil
	}

	err := huh.NewInput().
		Title("Please enter a path to save the config file template").
		Placeholder("Example: F:\\config.yaml").
		Validate(validate).
		Value(&results).Run()

	return results, err
}

func CreateConfigTemplate(savePath *string) {

	// check if the savePath is provided, if not, ask for it
	if savePath == nil {

		answer, err := askForSavePath()
		if err != nil {
			Log.Error("\nFailed to get user input\n")
			return
		}

		savePath = &answer
	}

	// check if the savePath lead to a .yaml file, if not, ask for a new one
	if !strings.HasSuffix(*savePath, ".yaml") {

		Log.Error("\nfile extension must be .yaml\n")
		answer, err := askForSavePath()
		if err != nil {
			Log.Error("\nFailed to get user input\n")
			return
		}

		savePath = &answer
	}

	configTemplate := `backup:
  # Files and folders paths to backup
  paths:
    - D:\data # Example: a folder path
    - F:\importantText.txt # Example: a file path
    - "%localappdata%\\app" # Example: a path with environment variable
    - C:\Users\%USERNAME%\Saved Games # Example: a path with environment variable

  # backup/restore paths to/from this path
  target: F:\backup # Example: a folder path

# A list of environment variables to be set
environmentVariables:
  - key: ANDROID_HOME
    value: F:\Android\Sdk
    scope: User # or Machine (Needs admin privileges)

  # Example: add a new entry to the "PATH" environment variable
  - key: PATH
    value: F:\Android\Sdk\platform-tools
    scope: User

# A list of packages to be installed using Chocolatey
packages:
  # --- BROWSERS ---
  - googlechrome # Example: install Google Chrome

  # --- PROGRAMS ---

  # --- GAMING ---

  # --- MEDIA ---

  # --- DRIVERS ---

  # --- TOOLS ---

  # --- DEV ---

# A list of scripts to be executed
scripts:
  # Example: single line with cmd shell
  - echo "Hello World!"

  # Example: multiline script with cmd shell
  - >
    echo "Hello World!"
    && echo "What's up?"

  # Example: single line with powershell
  - powershell echo "Hello World!"

  # Example: multiline script with powershell
  - >
    powershell $name = "David";
    echo "Hello $name!";
`

	// Create the file
	file, err := os.Create(*savePath)
	if err != nil {
		Log.Error("\n"+fmt.Sprintf(`error while creating the config file template at: "%s"`, *savePath), "\n")
		return
	}
	// Ensure the file is closed properly after writing
	defer file.Close()

	// Write some text to the file
	_, err = file.WriteString(configTemplate)
	if err != nil {
		Log.Error("\n"+fmt.Sprintf(`error while creating the config file template at: "%s"`, *savePath), "\n")
		return
	}

	Log.Success("\n" + fmt.Sprintf(`config file template created at: "%s"`, *savePath) + "\n")
}
