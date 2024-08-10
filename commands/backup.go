package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func BackupData(configFilePath *string) {

	// no config file path provided, ask for it
	if configFilePath == nil {
		answer := AskForConfigFilePath()

		// when the user exit the prompt using CTRL + C
		if !Utils.IsPathExists(answer) {
			Log.Error("\nfile not found. Please enter a valid path\n")
			return
		}

		configFilePath = &answer
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

	// paths is empty, exit
	if len(yamlData.Backup.Paths) == 0 {
		Log.Error("\nthe YAML file does not contain any backup paths\n")
		return
	}

	// create the target path
	isTargetPathExists := Utils.IsPathExists(yamlData.Backup.Target)
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
	PreparePathsString(yamlData.Backup.Paths)
	for _, path := range yamlData.Backup.Paths {

		Log.Info(fmt.Sprintf(`Copying "%s"`, path))

		baseFilename := filepath.Base(path)
		targetPath := filepath.Join(yamlData.Backup.Target, baseFilename)

		powershell := Powershell.GetShellPath()

		cmd := exec.Command(
			powershell,
			"-Command",
			"Copy-Item",
			"-Path", fmt.Sprintf(`"%s"`, path),
			"-Destination", fmt.Sprintf(`"%s"`, targetPath),
			"-Recurse", "-Force",
		)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			Log.Fatal(err.Error())
			os.Exit(1)
		}

		if err := cmd.Wait(); err != nil {
			Log.Fatal(err.Error())
			os.Exit(1)
		}

	}

	Log.Success("\nBackup completed\n")
}
