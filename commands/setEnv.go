package commands

import "fmt"

func SetEnvs(configFilePath *string) {
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
