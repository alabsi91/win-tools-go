package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func askForUsername() string {
	var results string

	validate := func(str string) error {
		if len(str) == 0 {
			return errors.New("please enter your username")
		}
		return nil
	}

	huh.NewInput().
		Title("Please enter your username").
		Placeholder("The username you use to logon").
		Validate(validate).
		Value(&results).Run()

	return results
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
		answer := askForUsername()
		username = &answer

		// when the user exit the prompt using CTRL + C
		if len(answer) == 0 {
			Log.Error("\nPlease enter your username\n")
			return
		}
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

	shell := Powershell.GetShellPath()

	cmd := exec.Command(
		shell,
		"-Command",
		"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
		fmt.Sprintf(`&"%s"`, scriptPath),
		scriptArgs,
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

	Log.Success("\nDone!\n")
}
