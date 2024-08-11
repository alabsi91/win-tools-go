package commands

import (
	"fmt"
	"path/filepath"

	"github.com/alabsi91/win-tools/commands/utils"
)

func RestoreData(configFilePath *string) {
	// no config file path provided, ask for it
	if configFilePath == nil {
		answer, err := utils.AskForConfigFilePath()
		if err != nil {
			Log.Error("\nFailed to get user input\n")
			return
		}

		configFilePath = &answer
	}

	// config file path provided does not exist, ask for a new one
	if !utils.IsPathExists(*configFilePath) {
		Log.Error("\nfile not found. Please enter a valid path\n")

		answer, err := utils.AskForConfigFilePath()
		if err != nil {
			Log.Error("\nFailed to get user input\n")
			return
		}

		configFilePath = &answer
	}

	yamlData := utils.ReadConfigFile(*configFilePath)

	// paths is empty, exit
	if len(yamlData.Backup.Paths) == 0 {
		Log.Error("\nthe YAML file does not contain any backup paths\n")
		return
	}

	// check the target path
	isTargetPathExists := utils.IsPathExists(yamlData.Backup.Target)
	if !isTargetPathExists {
		Log.Error(fmt.Sprintf(`the target path does not exist: "%s"`, yamlData.Backup.Target), "\n")
		return
	}

	Log.Warning("\nFiles and folders with the same name will be overwritten.\n")
	Log.Info(fmt.Sprintf(`Restoring data from: "%s"`, yamlData.Backup.Target), "\n")

	// loop over paths and copy the files and folders to the target path
	utils.PreparePathsString(yamlData.Backup.Paths)
	for _, toPath := range yamlData.Backup.Paths {
		fromPath := filepath.Join(yamlData.Backup.Target, filepath.Base(toPath))

		if utils.IsPathExists(toPath) {
			toPath = filepath.Dir(toPath)
		}

		Log.Info(fmt.Sprintf(`Copying "%s" to "%s"`, fromPath, toPath))

		err := Powershell.RunPathThroughCmd(
			"Copy-Item",
			"-Path", fmt.Sprintf(`"%s"`, fromPath),
			"-Destination", fmt.Sprintf(`"%s"`, toPath),
			"-Recurse", "-Force",
		)

		if err != nil {
			Log.Error("\nFailed to copy the path.\n", err.Error(), "\n")
		}
	}

	Log.Success("\nRestore completed\n")
}
