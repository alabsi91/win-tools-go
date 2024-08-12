package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func askForUsername() (string, error) {
	var results string

	validate := func(str string) error {
		if len(str) == 0 {
			return errors.New("please enter your username")
		}
		return nil
	}

	err := huh.NewInput().
		Title("\nPlease enter your username").
		Placeholder("The username you use to logon").
		Validate(validate).
		Value(&results).
		Run()

	return results, err
}

func AutoLogon(username *string, domain *string, autoLogonCount *int, removeLegalPrompt *bool, backupFile *string) {
	// has admin privileges
	isAdmin := Powershell.IsAdmin()

	if !isAdmin {
		Log.Error("\nyou need admin privileges to run this command\n")
		Log.Info("Please run this command from an elevated powershell session\n")
		return
	}

	// check if username is provided, if not ask for it
	if username == nil {

		answer, err := askForUsername()
		if err != nil {
			Log.Error("\nFailed to get user input\n")
			return
		}

		username = &answer
	}

	scriptArgs := fmt.Sprintf(`-Username "%s"`, *username)

	if domain != nil {
		scriptArgs = fmt.Sprintf(`%s -Domain "%s"`, scriptArgs, *domain)
	}

	if autoLogonCount != nil {
		scriptArgs = fmt.Sprintf(`%s -AutoLogonCount "%d"`, scriptArgs, *autoLogonCount)
	}

	if removeLegalPrompt != nil {
		scriptArgs = fmt.Sprintf(`%s -RemoveLegalPrompt`, scriptArgs)
	}

	scriptPath := filepath.Join(AssetsPath, "autologon.ps1")

	err := Powershell.RunPathThroughCmd(
		"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
		fmt.Sprintf(`&"%s"`, scriptPath),
		scriptArgs,
	)

	if err != nil {
		Log.Fatal("\n"+err.Error(), "\n")
	}

	Log.Success("\nDone!\n")
}
