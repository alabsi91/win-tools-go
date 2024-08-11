package commands

import (
	"fmt"

	"github.com/alabsi91/win-tools/commands/utils"
)

func SetEnvs(configFilePath *string) {
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

	// envs is empty, exit
	if len(yamlData.EnvironmentVariables) == 0 {
		Log.Error("\nthe YAML file does not contain any environment variables to set\n")
		return
	}

	// has admin privileges
	isAdmin := Powershell.IsAdmin()
	if !isAdmin {
		Log.Warning("\nEnvironment variables with \"Machine\" scope require admin privileges\n")
	}

	// loop through the envs
	for _, env := range yamlData.EnvironmentVariables {
		Log.Info(fmt.Sprintf(`Setting environment variable: %s="%s"`, env.Key, env.Value))
		Powershell.SetEnvVariable(env.Key, env.Value, env.Scope)
	}

	Log.Success("\nEnvironment variables set successfully\n")
}
