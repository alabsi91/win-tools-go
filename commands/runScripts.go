package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alabsi91/win-tools/commands/utils"
)

func RunScripts(configFilePath *string) {

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

	// scripts is empty, exit
	if len(yamlData.Scripts) == 0 {
		Log.Error("\nthe YAML file does not contain any scripts\n")
		return
	}

	// has admin privileges
	isAdmin := Powershell.IsAdmin()
	if !isAdmin {
		Log.Warning("\nYou may need admin privileges to run some scripts")
	}

	// loop through the scripts
	for i, script := range yamlData.Scripts {
		Log.Info("\n"+fmt.Sprintf(`Running the script with the index "%d"`, i), "\n")

		shell := "cmd"
		command := "/C"

		script, isPowershell := strings.CutPrefix(script, "powershell")
		if isPowershell {
			shell = Powershell.GetShellName()
			command = "-Command"
		}

		cmd := exec.Command(shell, command, script)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			Log.Error("\n" + fmt.Sprintf(`Failed to run the script with the index "%d"`, i))
			return
		}

		if err := cmd.Wait(); err != nil {
			Log.Error("\n" + fmt.Sprintf(`Failed to run the script with the index "%d"`, i))
			return
		}
	}

	Log.Success("\nAll scripts have been run successfully\n")
}
