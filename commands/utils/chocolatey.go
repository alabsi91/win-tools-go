package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
)

const ChocolateyInstallPath = "C:\\ProgramData\\chocolatey\\bin\\choco.exe"

type chocolatey struct{}

var Chocolatey = &chocolatey{}

// IsInstalled checks if Chocolatey is installed on the system.
//
// This method determines the presence of Chocolatey by verifying whether the
// command 'choco --version' can be executed successfully. If the command exists
// and returns a version, it indicates that Chocolatey is installed.
//
// Returns:
//   - true if Chocolatey is installed on the system.
//   - false if Chocolatey is not installed or the command cannot be found.
func (chocolatey) IsInstalled() bool {
	cmd := exec.Command("cmd", "/C", "choco", "--version")

	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return true
	}

	return IsPathExists(ChocolateyInstallPath)
}

// GetExecutablePath retrieves the path to the Chocolatey executable.
//  1. First, it checks the default installation path for Chocolatey.
//  2. If the executable is not found in the default path, it attempts to retrieve the path from the environment variables.
//  3. If neither method succeeds, the program will exit with an error, indicating that the Chocolatey executable could not be located.
//
// Returns: The absolute path to the Chocolatey executable as a string.
func (chocolatey) GetExecutablePath() string {
	// first try to the default path
	if IsPathExists(ChocolateyInstallPath) {
		return ChocolateyInstallPath
	}

	// try to get the path from where chocolatey is installed
	cmd := exec.Command("cmd", "/C", "where", "choco")

	output, err := cmd.Output()
	if err != nil {
		Log.Fatal("\nFailed to get chocolatey executable path\n")
	}

	// verify if the path exists
	path := string(output)

	if !IsPathExists(path) {
		Log.Fatal("\nFailed to get chocolatey executable path\n")
	}

	return path
}

// InstallSelf installs Chocolatey on the system using a PowerShell command.
//   - Uses a PowerShell command to perform the installation.
//   - Must be run with elevated privileges (administrator rights).
//   - Does not check for admin privileges before running the command.
//   - Exits the program if the installation command fails.
func (chocolatey) InstallSelf() {
	powershell := Powershell.GetShellName()

	cmd := exec.Command(
		powershell,
		"-Command",
		"Set-ExecutionPolicy Bypass -Scope Process -Force;",
		"[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072;",
		"iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		Log.Fatal("\n"+err.Error(), "\n")
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal("\n"+err.Error(), "\n")
	}
}

// InstallPackage installs a package using Chocolatey.
//   - Uses a PowerShell command to perform the installation and streams the output.
//   - On error, prints an error message without exiting the program or returning an error.
//   - Takes the package name as the first parameter (packageName).
//   - Takes a boolean as the second parameter (openInNewWindow):
//   - If true, starts the process in a new terminal window without waiting for it to finish.
//   - If false, runs the process in the current context.
func (chocolatey *chocolatey) InstallPackage(packageName string, openInNewWindow bool) {
	chocolateyPath := chocolatey.GetExecutablePath()
	chocolateyPath = fmt.Sprintf(`. "%s"`, chocolateyPath)

	powershell := Powershell.GetShellName()

	var err error
	if openInNewWindow {
		err = Powershell.RunPathThroughCmd(
			"Start-Process", powershell,
			"-ArgumentList", fmt.Sprintf(`'-C', '%s install %s -yf --ignore-checksum; pause'`, chocolateyPath, packageName),
			"-Verb", "RunAs",
		)
	} else {
		err = Powershell.RunPathThroughCmd(
			chocolateyPath,
			"install", packageName,
			"-yf", "--ignore-checksum",
		)
	}

	if err != nil {
		Log.Error("\nfailed to install chocolatey package:", packageName, "\n")
		return
	}

}

// AskForInstallConfirmation asks the user if Chocolatey should be installed.
//   - Displays a prompt asking the user if Chocolatey should be installed.
//   - If the user confirms/denies, returns a boolean.
//   - If the user cancels the prompt using Ctrl+C for example, returns an error.
//
// Returns:
//   - A boolean indicating if Chocolatey should be installed.
//   - An error if the user cancels the prompt or denies the installation.
func (chocolatey) AskForInstallConfirmation() (bool, error) {
	var answer bool = false

	err := huh.NewConfirm().
		Title("Chocolatey is not installed. Do you want to install it?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&answer).Run()

	return answer, err
}
