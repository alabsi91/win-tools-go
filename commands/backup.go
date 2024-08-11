package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alabsi91/win-tools/commands/utils"
)

var Powershell = utils.Powershell
var Log = utils.Log
var AssetsPath = utils.AssetsPath
var Chocolatey = utils.Chocolatey

func BackupData(configFilePath *string) {

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

	// create the target path
	isTargetPathExists := utils.IsPathExists(yamlData.Backup.Target)
	if !isTargetPathExists {
		Log.Info("\nthe target path does not exist. It will be created\n")
		// try to create the target path
		err := os.MkdirAll(yamlData.Backup.Target, os.ModePerm)
		if err != nil {
			Log.Error("\nthe YAML file does not contain any backup paths\n")
			return
		}
	}

	Log.Warning("\nFiles and folders with the same name will be overwritten.\n")
	Log.Info(fmt.Sprintf(`The target path is: "%s"`, yamlData.Backup.Target), "\n")

	// loop over paths and copy the files and folders to the target path
	utils.PreparePathsString(yamlData.Backup.Paths)
	for _, path := range yamlData.Backup.Paths {

		Log.Info(fmt.Sprintf(`Copying "%s"`, path))

		baseFilename := filepath.Base(path)
		targetPath := filepath.Join(yamlData.Backup.Target, baseFilename)

		err := Powershell.RunPathThroughCmd(
			"Copy-Item",
			"-Path", fmt.Sprintf(`"%s"`, path),
			"-Destination", fmt.Sprintf(`"%s"`, targetPath),
			"-Recurse", "-Force",
		)

		if err != nil {
			Log.Fatal(err.Error())
			os.Exit(1)
		}
	}

	Log.Success("\nBackup completed\n")
}
