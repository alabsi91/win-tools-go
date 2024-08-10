package commands

import (
	"fmt"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

func RestoreData(configFilePath *string) {
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

	// check the target path
	isTargetPathExists := Utils.IsPathExists(yamlData.Backup.Target)
	if !isTargetPathExists {
		Log.Error(fmt.Sprintf(`the target path does not exist: "%s"`, yamlData.Backup.Target), "\n")
		return
	}

	Log.Warning("\nFiles and folders with the same name will be overwritten.\n")
	Log.Info(fmt.Sprintf(`Restoring data from: "%s"`, yamlData.Backup.Target), "\n")

	// loop over paths and copy the files and folders to the target path
	PreparePathsString(yamlData.Backup.Paths)
	for _, path := range yamlData.Backup.Paths {
		fromPath := filepath.Join(yamlData.Backup.Target, filepath.Base(path))

		Log.Info(fmt.Sprintf(`Copying "%s" to "%s"`, fromPath, path))

		info, err := os.Stat(fromPath)
		if err != nil {
			Log.Error("\nerror getting path info:", fromPath, "\n")
			return
		}

		// copy folder recursively
		if info.IsDir() {
			err = cp.Copy(fromPath, path)
			if err != nil {
				Log.Error("\nerror copying folder", path, "\n")
				return
			}

			continue
		}

		// copy file
		err = Utils.CopyFile(fromPath, path)
		if err != nil {
			Log.Error("\nerror copying file", path, "\n")
			return
		}
	}

	Log.Success("\nRestore completed\n")
}
