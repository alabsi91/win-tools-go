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

func (chocolatey) IsInstalled() bool {
	cmd := exec.Command("cmd", "/C", "choco", "--version")

	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return true
	}

	return IsPathExists(ChocolateyInstallPath)
}

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
		os.Exit(1)
	}

	// verify if the path exists
	path := string(output)

	if !IsPathExists(path) {
		Log.Fatal("\nFailed to get chocolatey executable path\n")
		os.Exit(1)
	}

	return path
}

func (chocolatey) InstallSelf() {
	powershell := Powershell.GetShellPath()

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
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal("\n"+err.Error(), "\n")
		os.Exit(1)
	}
}

func (chocolatey *chocolatey) InstallPackage(packageName string) {
	chocolateyPath := chocolatey.GetExecutablePath()
	powershell := Powershell.GetShellPath()

	cmd := exec.Command(
		powershell,
		"-Command",
		fmt.Sprintf(`&"%s"`, chocolateyPath),
		"install",
		packageName,
		"-yf",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		Log.Fatal("\nFailed to install chocolatey package:", packageName, "\n")
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal("\nFailed to install chocolatey package:", packageName, "\n")
		os.Exit(1)
	}
}

func (chocolatey) AskForInstallConfirmation() (bool, error) {
	var answer bool = false

	err := huh.NewConfirm().
		Title("Chocolatey is not installed. Do you want to install it?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&answer).Run()

	return answer, err
}
