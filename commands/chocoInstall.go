package commands

import (
	"fmt"
	"strings"

	"github.com/alabsi91/win-tools/commands/utils"
)

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
		answer, err := utils.AskForConfigFilePath()
		if err != nil {
			Log.Error("failed to get user input\n")
			return
		}

		configFilePath = &answer
	}

	// config file path provided does not exist, ask for a new one
	if !utils.IsPathExists(*configFilePath) {
		Log.Error("\nfile not found. Please enter a valid path\n")

		answer, err := utils.AskForConfigFilePath()
		if err != nil {
			Log.Error("failed to get user input\n")
			return
		}

		configFilePath = &answer
	}

	yamlData := utils.ReadConfigFile(*configFilePath)

	// packages is empty, exit
	if len(yamlData.Packages) == 0 {
		Log.Error("\nthe YAML file does not contain any packages\n")
		return
	}

	// check if chocolatey is installed
	isChocolateyInstalled := Chocolatey.IsInstalled()
	if !isChocolateyInstalled {
		answer, err := Chocolatey.AskForInstallConfirmation()
		if !answer || err != nil {
			return
		}

		// install chocolatey
		Chocolatey.InstallSelf()
	}

	Log.Info("\n" + fmt.Sprintf(`Found "%d" packages`, len(yamlData.Packages)))

	// loop through packages and install them
	for _, packageName := range yamlData.Packages {

		openInNewWindow := false
		if strings.Contains(packageName, "--new-window") {
			openInNewWindow = true
			packageName = strings.ReplaceAll(packageName, "--new-window", "")
		}

		Log.Info("\n"+fmt.Sprintf(`Installing package: "%s"`, packageName), "\n")

		Chocolatey.InstallPackage(packageName, openInNewWindow)
	}

	Log.Success("\nDone\n")
}
