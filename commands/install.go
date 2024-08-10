package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

func askToInstallChocolatey() bool {
	var answer bool = false

	huh.NewConfirm().
		Title("Chocolatey is not installed. Do you want to install it?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&answer).Run()

	return answer
}

func InstallPackages(configFilePath *string) {
	// has admin privileges
	isAdmin := Powershell.IsAdmin()

	if !isAdmin {
		Log.Error("\nyou need admin privileges to install packages\n")
		Log.Info("Please run this command from an elevated powershell session\n")
		return
	}

	// no config file path provided, ask for it
	if configFilePath == nil {
		answer := AskForConfigFilePath()
		configFilePath = &answer

		// when the user exit the prompt using CTRL + C
		if !Utils.IsPathExists(answer) {
			Log.Error("\nfile not found. Please enter a valid path\n")
			return
		}
	}

	// config file path provided does not exist, ask for a new one
	if !Utils.IsPathExists(*configFilePath) {
		Log.Error("\nfile not found. Please enter a valid path\n")
		answer := AskForConfigFilePath()

		// when the user exit the prompt using CTRL + C
		if !Utils.IsPathExists(answer) {
			Log.Error("\nfile not found. Please enter a valid path\n")
			return
		}

		configFilePath = &answer
	}

	yamlData := ReadConfigFile(*configFilePath)

	// packages is empty, exit
	if len(yamlData.Packages) == 0 {
		Log.Error("\nthe YAML file does not contain any packages\n")
		return
	}

	// check if chocolatey is installed
	isChocolateyInstalled := Chocolatey.IsChocolateyInstalled()
	if !isChocolateyInstalled {
		answer := askToInstallChocolatey()
		if !answer {
			return
		}

		// install chocolatey
		Chocolatey.InstallChocolatey()
	}

	Log.Info("\n" + fmt.Sprintf(`Found "%d" packages`, len(yamlData.Packages)))

	Chocolatey.InstallChocolateyPackage(strings.Join(yamlData.Packages, " "))

	Log.Success("\nDone\n")
}
